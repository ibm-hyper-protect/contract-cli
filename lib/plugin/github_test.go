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
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPluginDir_CreatesDirectory verifies PluginDir creates the directory and returns a valid path.
func TestPluginDir_CreatesDirectory(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	dir, err := PluginDir()
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(tmpHome, PluginDirName), dir)
	_, statErr := os.Stat(dir)
	assert.NoError(t, statErr)
}

// TestListPlugins_ReturnsOnlyPluginPrefixedFiles verifies ListPlugins filters correctly.
func TestListPlugins_ReturnsOnlyPluginPrefixedFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create plugin binaries and a non-plugin file
	for _, name := range []string{"contract-plugin-foo", "contract-plugin-bar", "notaplugin"} {
		f, err := os.Create(filepath.Join(tmpDir, name))
		assert.NoError(t, err)
		f.Close()
	}

	plugins, err := ListPlugins(tmpDir)
	assert.NoError(t, err)
	assert.Len(t, plugins, 2)
	assert.Contains(t, plugins, "contract-plugin-foo")
	assert.Contains(t, plugins, "contract-plugin-bar")
	assert.NotContains(t, plugins, "notaplugin")
}

// TestListPlugins_EmptyDir verifies ListPlugins returns empty slice for a directory with no plugins.
func TestListPlugins_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	plugins, err := ListPlugins(tmpDir)
	assert.NoError(t, err)
	assert.Empty(t, plugins)
}

// TestParsePluginRef_WithVersion verifies version extraction.
func TestParsePluginRef_WithVersion(t *testing.T) {
	repo, version, err := ParsePluginRef("github.com/org/contract-plugin-foo@v1.2.0")
	assert.NoError(t, err)
	assert.Equal(t, "github.com/org/contract-plugin-foo", repo)
	assert.Equal(t, "v1.2.0", version)
}

// TestParsePluginRef_WithoutVersion verifies default of "latest".
func TestParsePluginRef_WithoutVersion(t *testing.T) {
	repo, version, err := ParsePluginRef("github.com/org/contract-plugin-foo")
	assert.NoError(t, err)
	assert.Equal(t, "github.com/org/contract-plugin-foo", repo)
	assert.Equal(t, "latest", version)
}

// TestParsePluginRef_EmptyRef verifies error on empty input.
func TestParsePluginRef_EmptyRef(t *testing.T) {
	_, _, err := ParsePluginRef("")
	assert.Error(t, err)
}

// TestParsePluginRef_NonGitHub verifies error for non-github.com refs.
func TestParsePluginRef_NonGitHub(t *testing.T) {
	_, _, err := ParsePluginRef("gitlab.com/org/contract-plugin-foo")
	assert.Error(t, err)
}

const (
	testPluginRepo   = "github.com/org/contract-plugin-myplugin"
	testPluginBinary = "contract-plugin-myplugin"
)

// testAssetName returns the release asset name expected for the current platform.
func testAssetName() string {
	name := testPluginBinary + "-" + runtime.GOOS + "-" + runtime.GOARCH
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return name
}

// newReleaseServer starts a fake GitHub API + download server. extraAssets maps an extra
// asset name to its content; the platform binary is always served with binaryContent.
// It points githubAPIBaseURL at the server for the duration of the test.
func newReleaseServer(t *testing.T, binaryContent []byte, extraAssets map[string]string) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/repos/") {
			assets := []map[string]string{
				{"name": testAssetName(), "browser_download_url": "http://" + r.Host + "/download/" + testAssetName()},
			}
			for name := range extraAssets {
				assets = append(assets, map[string]string{"name": name, "browser_download_url": "http://" + r.Host + "/download/" + name})
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"assets": assets})
			return
		}
		if name := strings.TrimPrefix(r.URL.Path, "/download/"); name != r.URL.Path {
			if name == testAssetName() {
				_, _ = w.Write(binaryContent)
				return
			}
			if content, ok := extraAssets[name]; ok {
				_, _ = w.Write([]byte(content))
				return
			}
		}
		http.NotFound(w, r)
	}))
	t.Cleanup(srv.Close)

	orig := githubAPIBaseURL
	githubAPIBaseURL = srv.URL
	t.Cleanup(func() { githubAPIBaseURL = orig })
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// TestDownloadGitHubRelease_PerAssetChecksum_Success verifies install with a matching "<asset>.sha256" file.
func TestDownloadGitHubRelease_PerAssetChecksum_Success(t *testing.T) {
	content := []byte("fake-binary-content")
	newReleaseServer(t, content, map[string]string{
		testAssetName() + ".sha256": sha256Hex(content) + "  " + testAssetName() + "\n",
	})
	dir := t.TempDir()

	path, err := DownloadGitHubRelease(testPluginRepo, "latest", dir, false)
	assert.NoError(t, err)
	got, readErr := os.ReadFile(path)
	assert.NoError(t, readErr)
	assert.Equal(t, content, got)
}

// TestDownloadGitHubRelease_AggregateChecksum_Success verifies install with a checksums.txt file.
func TestDownloadGitHubRelease_AggregateChecksum_Success(t *testing.T) {
	content := []byte("fake-binary-content")
	newReleaseServer(t, content, map[string]string{
		"checksums.txt": strings.Repeat("0", 64) + "  other-asset\n" + sha256Hex(content) + "  " + testAssetName() + "\n",
	})

	_, err := DownloadGitHubRelease(testPluginRepo, "v1.0.0", t.TempDir(), false)
	assert.NoError(t, err)
}

// TestDownloadGitHubRelease_ChecksumMismatch_Error verifies a bad checksum fails and installs nothing.
func TestDownloadGitHubRelease_ChecksumMismatch_Error(t *testing.T) {
	newReleaseServer(t, []byte("tampered"), map[string]string{
		"checksums.txt": strings.Repeat("a", 64) + "  " + testAssetName() + "\n",
	})
	dir := t.TempDir()

	_, err := DownloadGitHubRelease(testPluginRepo, "latest", dir, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "checksum mismatch")
	entries, _ := os.ReadDir(dir)
	assert.Empty(t, entries, "no file (including temp files) should remain after a failed install")
}

// TestDownloadGitHubRelease_MissingChecksumFile_Error verifies install fails when the release has no checksum file.
func TestDownloadGitHubRelease_MissingChecksumFile_Error(t *testing.T) {
	newReleaseServer(t, []byte("content"), nil)

	_, err := DownloadGitHubRelease(testPluginRepo, "latest", t.TempDir(), false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--insecure")
}

// TestDownloadGitHubRelease_NoEntryInChecksumFile_Error verifies install fails when the asset is not listed.
func TestDownloadGitHubRelease_NoEntryInChecksumFile_Error(t *testing.T) {
	newReleaseServer(t, []byte("content"), map[string]string{
		"checksums.txt": strings.Repeat("a", 64) + "  other-asset\n",
	})

	_, err := DownloadGitHubRelease(testPluginRepo, "latest", t.TempDir(), false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no checksum entry")
}

// TestDownloadGitHubRelease_Insecure_SkipsChecksum verifies --insecure installs without any checksum file.
func TestDownloadGitHubRelease_Insecure_SkipsChecksum(t *testing.T) {
	content := []byte("content")
	newReleaseServer(t, content, nil)

	path, err := DownloadGitHubRelease(testPluginRepo, "latest", t.TempDir(), true)
	assert.NoError(t, err)
	got, readErr := os.ReadFile(path)
	assert.NoError(t, readErr)
	assert.Equal(t, content, got)
}

// TestDownloadGitHubRelease_Insecure_IgnoresMismatch verifies --insecure does not verify a present checksum file.
func TestDownloadGitHubRelease_Insecure_IgnoresMismatch(t *testing.T) {
	newReleaseServer(t, []byte("content"), map[string]string{
		"checksums.txt": strings.Repeat("a", 64) + "  " + testAssetName() + "\n",
	})

	_, err := DownloadGitHubRelease(testPluginRepo, "latest", t.TempDir(), true)
	assert.NoError(t, err)
}

// TestDownloadGitHubRelease_ChecksumMismatch_KeepsExistingPlugin verifies a failed update leaves the old binary intact.
func TestDownloadGitHubRelease_ChecksumMismatch_KeepsExistingPlugin(t *testing.T) {
	newReleaseServer(t, []byte("tampered"), map[string]string{
		"checksums.txt": strings.Repeat("a", 64) + "  " + testAssetName() + "\n",
	})
	dir := t.TempDir()
	dest := filepath.Join(dir, testPluginBinary)
	if runtime.GOOS == "windows" {
		dest += ".exe"
	}
	assert.NoError(t, os.WriteFile(dest, []byte("old"), 0755))

	_, err := DownloadGitHubRelease(testPluginRepo, "latest", dir, false)
	assert.Error(t, err)
	got, _ := os.ReadFile(dest)
	assert.Equal(t, "old", string(got))
}

// TestParseChecksum_BareDigest verifies a single bare digest is accepted.
func TestParseChecksum_BareDigest(t *testing.T) {
	d := strings.Repeat("ab", 32)
	got, err := parseChecksum(d+"\n", "any")
	assert.NoError(t, err)
	assert.Equal(t, d, got)
}

// TestParseChecksum_BinaryModeMarker verifies the "*name" sha256sum binary marker is handled.
func TestParseChecksum_BinaryModeMarker(t *testing.T) {
	d := strings.Repeat("cd", 32)
	got, err := parseChecksum(d+" *asset\n", "asset")
	assert.NoError(t, err)
	assert.Equal(t, d, got)
}

// TestParseChecksum_InvalidDigest verifies a malformed digest is rejected.
func TestParseChecksum_InvalidDigest(t *testing.T) {
	_, err := parseChecksum("nothex  asset\n", "asset")
	assert.Error(t, err)
	_, err = parseChecksum("abcd  asset\n", "asset")
	assert.Error(t, err)
}
