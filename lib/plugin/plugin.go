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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ibm-hyper-protect/contract-cli/common"
)

const (
	// ParameterName is the Cobra command name for the plugin group.
	ParameterName = "plugin"
	// ParameterShortDescription is the one-line description shown in help.
	ParameterShortDescription = "Manage contract-cli plugins"
	// ParameterLongDescription is the full description shown in plugin --help.
	ParameterLongDescription = `Manage contract-cli plugins.

Plugins are standalone binaries named contract-plugin-<name> stored in
~/.contract-cli/plugins/. Any unknown sub-command is transparently delegated
to the matching plugin binary.

Plugin repositories follow the naming convention contract-plugin-<name>.`

	// PluginPrefix is the required prefix for all plugin binaries.
	PluginPrefix = "contract-plugin-"
	// PluginDirName is the directory under $HOME where plugins are stored.
	PluginDirName = ".contract-cli/plugins"

	// InsecureFlag is the flag that skips release checksum verification.
	InsecureFlag = "insecure"
	// InsecureFlagDescription is the help text for the --insecure flag.
	InsecureFlagDescription = "skip SHA-256 checksum verification of the downloaded plugin (not recommended)"

	// SourcesFile is the JSON sidecar that maps plugin name → GitHub repo.
	SourcesFile = ".sources.json"

	// versionTimeout is the maximum time allowed for a plugin --version call.
	versionTimeout = 3 * time.Second
)

// PluginInfo holds metadata about an installed plugin.
type PluginInfo struct {
	Name    string
	Path    string
	Version string
}

// PluginDir returns the absolute path to the plugin directory (~/.contract-cli/plugins),
// creating it if it does not exist.
func PluginDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to determine home directory: %w", err)
	}
	dir := filepath.Join(home, PluginDirName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("unable to create plugin directory %s: %w", dir, err)
	}
	return dir, nil
}

// ListPlugins returns the names (without path) of installed plugin binaries in dir.
// Only entries whose names start with PluginPrefix are returned.
func ListPlugins(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to read plugin directory: %w", err)
	}
	var plugins []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), PluginPrefix) {
			plugins = append(plugins, e.Name())
		}
	}
	return plugins, nil
}

// sourcesMap is the in-memory representation of .sources.json.
type sourcesMap map[string]string

// readSources loads the sources sidecar from the plugin directory.
// Returns an empty map (not an error) if the file does not exist yet.
func readSources(pluginDir string) (sourcesMap, error) {
	path := filepath.Join(pluginDir, SourcesFile)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return make(sourcesMap), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin sources: %w", err)
	}
	var m sourcesMap
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse plugin sources: %w", err)
	}
	return m, nil
}

// writeSources persists the sources map to disk.
func writeSources(pluginDir string, m sourcesMap) error {
	path := filepath.Join(pluginDir, SourcesFile)
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal plugin sources: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// ValidateInstallInput validates that the positional argument for install/uninstall/update is non-empty.
func ValidateInstallInput(args []string) (string, error) {
	if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
		return "", fmt.Errorf("plugin reference argument is required")
	}
	return args[0], nil
}

// Install installs a plugin from a GitHub release URL or a local file path.
//
// ref may be:
//   - "github.com/org/contract-plugin-foo"          (latest)
//   - "github.com/org/contract-plugin-foo@v1.2.0"   (pinned)
//   - "/path/to/contract-plugin-foo"                 (local binary)
//
// For GitHub refs the release's SHA-256 checksum file is verified unless insecure is true.
// insecure has no effect on local installs.
func Install(ref, pluginDir string, insecure bool) error {
	if strings.HasPrefix(ref, "github.com/") || strings.Contains(ref, "@") && strings.HasPrefix(ref[:strings.Index(ref, "@")], "github.com/") {
		return installFromGitHub(ref, pluginDir, insecure)
	}
	return installFromLocal(ref, pluginDir)
}

func installFromGitHub(ref, pluginDir string, insecure bool) error {
	repo, version, err := ParsePluginRef(ref)
	if err != nil {
		return err
	}

	destPath, err := DownloadGitHubRelease(repo, version, pluginDir, insecure)
	if err != nil {
		return err
	}

	// Record the source so 'plugin update' can find the repo later.
	sources, err := readSources(pluginDir)
	if err != nil {
		return err
	}
	// Key is the binary name without extension.
	binaryName := filepath.Base(destPath)
	binaryName = strings.TrimSuffix(binaryName, ".exe")
	sources[binaryName] = repo
	return writeSources(pluginDir, sources)
}

func installFromLocal(ref, pluginDir string) error {
	if !common.CheckFileFolderExists(ref) {
		return fmt.Errorf("local plugin path does not exist: %s", ref)
	}
	base := filepath.Base(ref)
	if !strings.HasPrefix(base, PluginPrefix) {
		return fmt.Errorf("local binary name must start with %q, got: %s", PluginPrefix, base)
	}
	destPath := filepath.Join(pluginDir, base)
	src, err := os.Open(ref)
	if err != nil {
		return fmt.Errorf("failed to open source binary: %w", err)
	}
	defer src.Close()

	dst, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755) // #nosec G306
	if err != nil {
		return fmt.Errorf("failed to create destination binary: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy plugin binary: %w", err)
	}
	return nil
}

// PluginVersion calls the plugin binary with --version (with CONTRACT_CLI_PLUGIN=1 set)
// and returns the trimmed first line of its output.
// Returns "unknown" on any error or timeout.
func PluginVersion(binaryPath string) string {
	ctx, cancel := context.WithTimeout(context.Background(), versionTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, binaryPath, "--version") // #nosec G204
	cmd.Env = append(os.Environ(), HandshakeEnvVar+"="+HandshakeEnvVal)
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	line := strings.SplitN(strings.TrimSpace(string(out)), "\n", 2)[0]
	if line == "" {
		return "unknown"
	}
	return line
}

// List returns metadata for all installed plugins, including their versions.
func List(pluginDir string) ([]PluginInfo, error) {
	names, err := ListPlugins(pluginDir)
	if err != nil {
		return nil, err
	}
	var infos []PluginInfo
	for _, name := range names {
		path := filepath.Join(pluginDir, name)
		infos = append(infos, PluginInfo{
			Name:    strings.TrimPrefix(name, PluginPrefix),
			Path:    path,
			Version: PluginVersion(path),
		})
	}
	return infos, nil
}

// Uninstall removes the named plugin from the plugin directory.
// name is the short plugin name (e.g. "ccrv", NOT "contract-plugin-ccrv").
func Uninstall(name, pluginDir string) error {
	binaryName := PluginPrefix + name
	destPath := filepath.Join(pluginDir, binaryName)
	if !common.CheckFileFolderExists(destPath) {
		return fmt.Errorf("plugin %q is not installed", name)
	}
	if err := os.Remove(destPath); err != nil {
		return fmt.Errorf("failed to remove plugin %q: %w", name, err)
	}

	// Remove from sources tracking (ignore missing entry).
	sources, err := readSources(pluginDir)
	if err == nil {
		delete(sources, binaryName)
		_ = writeSources(pluginDir, sources)
	}
	return nil
}

// Update reinstalls one or all plugins.
//
//   - ref == ""              → update all GitHub-sourced installed plugins
//   - ref == "foo"           → update plugin "foo" to latest
//   - ref == "foo@v2.0.0"   → update plugin "foo" to v2.0.0
func Update(ref, pluginDir string, insecure bool) error {
	if ref == "" {
		return updateAll(pluginDir, insecure)
	}
	return updateOne(ref, pluginDir, insecure)
}

func updateOne(ref, pluginDir string, insecure bool) error {
	// Strip any @version suffix to get the name, then reconstruct the full ref.
	name := ref
	version := "latest"
	if idx := strings.LastIndex(ref, "@"); idx != -1 {
		name = ref[:idx]
		version = ref[idx+1:]
	}

	binaryName := PluginPrefix + name
	sources, err := readSources(pluginDir)
	if err != nil {
		return err
	}
	repo, ok := sources[binaryName]
	if !ok {
		return fmt.Errorf("plugin %q was not installed from GitHub and cannot be auto-updated; re-install manually", name)
	}

	fullRef := repo
	if version != "latest" {
		fullRef = repo + "@" + version
	}
	return installFromGitHub(fullRef, pluginDir, insecure)
}

func updateAll(pluginDir string, insecure bool) error {
	names, err := ListPlugins(pluginDir)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return fmt.Errorf("no plugins are installed")
	}
	sources, err := readSources(pluginDir)
	if err != nil {
		return err
	}

	var errs []string
	for _, binaryName := range names {
		repo, ok := sources[binaryName]
		if !ok {
			fmt.Fprintf(os.Stderr, "warning: plugin %q was installed from a local path and cannot be auto-updated; skipping\n", binaryName)
			continue
		}
		if err := installFromGitHub(repo, pluginDir, insecure); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", binaryName, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("some plugins failed to update:\n  %s", strings.Join(errs, "\n  "))
	}
	return nil
}
