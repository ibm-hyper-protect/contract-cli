// Copyright (c) 2026 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plugin

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	// githubAPIBaseURL is the GitHub REST API root. It is a variable so tests can point it at a local server.
	githubAPIBaseURL = "https://api.github.com"
	// httpClient is used for all plugin downloads and has a bounded timeout.
	httpClient = &http.Client{Timeout: 5 * time.Minute}
)

// maxChecksumFileSize caps how much of a checksum file is read.
const maxChecksumFileSize = 1 << 20

// githubAsset is a single downloadable file attached to a GitHub release.
type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// githubRelease is a minimal representation of the GitHub Releases API response.
type githubRelease struct {
	Assets []githubAsset `json:"assets"`
}

// ParsePluginRef splits a plugin reference such as
// "github.com/org/contract-plugin-foo@v1.2.0" into repo and version.
// Version defaults to "latest" when the @version suffix is omitted.
// Returns an error if the ref is empty or does not start with "github.com/".
func ParsePluginRef(ref string) (repo, version string, err error) {
	if ref == "" {
		return "", "", fmt.Errorf("plugin reference must not be empty")
	}
	version = "latest"
	if idx := strings.LastIndex(ref, "@"); idx != -1 {
		version = ref[idx+1:]
		ref = ref[:idx]
	}
	if !strings.HasPrefix(ref, "github.com/") {
		return "", "", fmt.Errorf("plugin reference must start with github.com/, got: %s", ref)
	}
	return ref, version, nil
}

// findChecksumAsset locates the checksum file for assetName in a release.
// A per-asset "<assetName>.sha256" file is preferred; otherwise an aggregate file
// named "checksums.txt", "SHA256SUMS", "SHA256SUMS.txt" or "*_checksums.txt" (GoReleaser default) is used.
func findChecksumAsset(assets []githubAsset, assetName string) (githubAsset, bool) {
	for _, a := range assets {
		if a.Name == assetName+".sha256" {
			return a, true
		}
	}
	for _, a := range assets {
		n := strings.ToLower(a.Name)
		if n == "checksums.txt" || n == "sha256sums" || n == "sha256sums.txt" || strings.HasSuffix(n, "_checksums.txt") {
			return a, true
		}
	}
	return githubAsset{}, false
}

// parseChecksum extracts the SHA-256 hex digest for assetName from the content of a
// checksum file in "sha256sum" format ("<hex>  <name>"). A file containing a single
// bare digest is also accepted.
func parseChecksum(content, assetName string) (string, error) {
	var lines []string
	for _, l := range strings.Split(content, "\n") {
		if strings.TrimSpace(l) != "" {
			lines = append(lines, l)
		}
	}

	var digest string
	if len(lines) == 1 && len(strings.Fields(lines[0])) == 1 {
		digest = strings.Fields(lines[0])[0]
	} else {
		for _, l := range lines {
			fields := strings.Fields(l)
			if len(fields) == 2 && strings.TrimPrefix(fields[1], "*") == assetName {
				digest = fields[0]
				break
			}
		}
	}
	if digest == "" {
		return "", fmt.Errorf("no checksum entry for %q found in checksum file", assetName)
	}
	digest = strings.ToLower(digest)
	if b, err := hex.DecodeString(digest); err != nil || len(b) != sha256.Size {
		return "", fmt.Errorf("invalid SHA-256 checksum for %q in checksum file", assetName)
	}
	return digest, nil
}

// fetchChecksum downloads a checksum file and returns the expected digest for assetName.
func fetchChecksum(url, assetName string) (string, error) {
	resp, err := httpClient.Get(url) // #nosec G107 — URL comes from the GitHub API
	if err != nil {
		return "", fmt.Errorf("failed to download checksum file: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("checksum download returned HTTP %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxChecksumFileSize))
	if err != nil {
		return "", fmt.Errorf("failed to read checksum file: %w", err)
	}
	return parseChecksum(string(data), assetName)
}

// DownloadGitHubRelease downloads the platform-matching binary from a GitHub release,
// writes it to destDir as "contract-plugin-<name>", sets 0755 permissions, and returns
// the full path to the installed binary.
//
// repo must be in the form "github.com/org/contract-plugin-<name>".
// version is a tag name (e.g. "v1.2.0") or "latest".
// Asset naming convention: contract-plugin-<name>-<os>-<arch>[.exe]
// where os ∈ {linux, darwin, windows} and arch ∈ {amd64, arm64, s390x, ppc64le}.
//
// Unless insecure is true, the release must contain a SHA-256 checksum file (see
// findChecksumAsset) and the downloaded binary must match it; on any failure the
// binary is not installed.
func DownloadGitHubRelease(repo, version, destDir string, insecure bool) (string, error) {
	// Extract owner/repo path from "github.com/owner/repo"
	parts := strings.SplitN(strings.TrimPrefix(repo, "github.com/"), "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid GitHub repo path: %s", repo)
	}
	ownerRepo := parts[0] + "/" + parts[1]
	pluginName := parts[1] // e.g. "contract-plugin-foo"

	// Build GitHub API URL
	var apiURL string
	if version == "latest" {
		apiURL = fmt.Sprintf("%s/repos/%s/releases/latest", githubAPIBaseURL, ownerRepo)
	} else {
		apiURL = fmt.Sprintf("%s/repos/%s/releases/tags/%s", githubAPIBaseURL, ownerRepo, version)
	}

	// Fetch release metadata
	resp, err := httpClient.Get(apiURL) // #nosec G107 — URL is constructed from validated user input
	if err != nil {
		return "", fmt.Errorf("failed to fetch release info from %s: %w", apiURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned %d for %s", resp.StatusCode, apiURL)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse release JSON: %w", err)
	}

	// Determine expected asset name for the current platform
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	expectedName := fmt.Sprintf("%s-%s-%s", pluginName, goos, goarch)
	if goos == "windows" {
		expectedName += ".exe"
	}

	// Find matching asset
	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == expectedName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		return "", fmt.Errorf("no asset named %q found in release; available assets may not include your platform (%s/%s)", expectedName, goos, goarch)
	}

	// Resolve the expected checksum before downloading the binary.
	var expectedSum string
	if insecure {
		fmt.Fprintf(os.Stderr, "warning: --insecure set; skipping checksum verification for %s\n", expectedName)
	} else {
		checksumAsset, ok := findChecksumAsset(release.Assets, expectedName)
		if !ok {
			return "", fmt.Errorf("no checksum file found in release for %q (expected %q.sha256, checksums.txt or SHA256SUMS); use --insecure to skip verification", expectedName, expectedName)
		}
		expectedSum, err = fetchChecksum(checksumAsset.BrowserDownloadURL, expectedName)
		if err != nil {
			return "", err
		}
	}

	// Download the binary
	binResp, err := httpClient.Get(downloadURL) // #nosec G107 — URL comes from authenticated GitHub API
	if err != nil {
		return "", fmt.Errorf("failed to download plugin binary: %w", err)
	}
	defer binResp.Body.Close()
	if binResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download returned HTTP %d", binResp.StatusCode)
	}

	// Write to a temp file in destDir so a failed or mismatching download never replaces an installed plugin.
	tmp, err := os.CreateTemp(destDir, ".download-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file in %s: %w", destDir, err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath) // no-op after a successful rename

	hasher := sha256.New()
	if _, err := io.Copy(io.MultiWriter(tmp, hasher), binResp.Body); err != nil {
		tmp.Close()
		return "", fmt.Errorf("failed to write plugin binary: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return "", fmt.Errorf("failed to write plugin binary: %w", err)
	}

	if !insecure {
		if actual := hex.EncodeToString(hasher.Sum(nil)); actual != expectedSum {
			return "", fmt.Errorf("checksum mismatch for %q: expected %s, got %s", expectedName, expectedSum, actual)
		}
	}

	if err := os.Chmod(tmpPath, 0755); err != nil { // #nosec G302 — plugins must be executable
		return "", fmt.Errorf("failed to set permissions on plugin binary: %w", err)
	}

	destName := pluginName
	if goos == "windows" {
		destName += ".exe"
	}
	destPath := filepath.Join(destDir, destName)
	if err := os.Rename(tmpPath, destPath); err != nil {
		return "", fmt.Errorf("failed to install plugin binary at %s: %w", destPath, err)
	}

	return destPath, nil
}
