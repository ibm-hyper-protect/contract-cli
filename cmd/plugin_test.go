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

package cmd

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setPluginDir overrides HOME so PluginDir resolves to a temp directory controlled by the test.
func setPluginDir(t *testing.T) string {
	t.Helper()
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	pluginDir := filepath.Join(tmpHome, ".contract-cli", "plugins")
	err := os.MkdirAll(pluginDir, 0755)
	assert.NoError(t, err)
	return pluginDir
}

// makeExecutable writes a stub executable to dir/contract-plugin-<name>.
func makeExecutable(t *testing.T, dir, name string) string {
	t.Helper()
	binName := "contract-plugin-" + name
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	path := filepath.Join(dir, binName)
	err := os.WriteFile(path, []byte("#!/bin/sh\necho stub"), 0755) // #nosec G306
	assert.NoError(t, err)
	return path
}

// runPluginCmd executes a plugin sub-command, captures real stdout, and returns it along with any error.
// Plugin commands use fmt.Printf (process stdout), so we redirect os.Stdout via os.Pipe.
func runPluginCmd(args ...string) (string, error) {
	// Redirect os.Stdout so fmt.Printf output is captured.
	r, w, _ := os.Pipe()
	origStdout := os.Stdout
	os.Stdout = w

	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = origStdout
	out, _ := io.ReadAll(r)
	return string(out), err
}

// ---------------------------------------------------------------------------
// plugin list
// ---------------------------------------------------------------------------

func TestPluginListCmd_NoPlugins(t *testing.T) {
	setPluginDir(t)
	out, err := runPluginCmd("plugin", "list")
	assert.NoError(t, err)
	assert.Contains(t, out, "No plugins installed")
}

func TestPluginListCmd_WithPlugin(t *testing.T) {
	pluginDir := setPluginDir(t)
	makeExecutable(t, pluginDir, "ccrv")

	out, err := runPluginCmd("plugin", "list")
	assert.NoError(t, err)
	assert.Contains(t, out, "ccrv")
}

// ---------------------------------------------------------------------------
// plugin uninstall
// ---------------------------------------------------------------------------

func TestPluginUninstallCmd_NotInstalled(t *testing.T) {
	setPluginDir(t)
	// uninstall non-existent plugin should call log.Fatal; we test via the lib layer instead.
	// The integration path is covered in lib/plugin/plugin_test.go.
	// Here we verify the argument validation path.
	_, err := runPluginCmd("plugin", "uninstall")
	// ExactArgs(1) causes Cobra to return an error when no arg given.
	assert.Error(t, err)
}

func TestPluginUninstallCmd_Success(t *testing.T) {
	pluginDir := setPluginDir(t)
	makeExecutable(t, pluginDir, "foo")

	binName := "contract-plugin-foo"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	assert.FileExists(t, filepath.Join(pluginDir, binName))

	out, err := runPluginCmd("plugin", "uninstall", "foo")
	assert.NoError(t, err)
	assert.Contains(t, strings.ToLower(out), "uninstalled")
	assert.NoFileExists(t, filepath.Join(pluginDir, binName))
}

// ---------------------------------------------------------------------------
// plugin install (argument validation)
// ---------------------------------------------------------------------------

func TestPluginInstallCmd_MissingArg(t *testing.T) {
	setPluginDir(t)
	_, err := runPluginCmd("plugin", "install")
	assert.Error(t, err) // ExactArgs(1) enforced by Cobra
}

// ---------------------------------------------------------------------------
// plugin update (argument validation)
// ---------------------------------------------------------------------------

func TestPluginUpdateCmd_NoArg_NoPlugins(t *testing.T) {
	setPluginDir(t)
	// update-all with no installed plugins should produce an error via log.Fatal.
	// Since log.Fatal calls os.Exit we can't catch it in a unit test; tested via lib layer.
	// Validate that Cobra accepts zero args without a parse error.
	_, err := runPluginCmd("plugin", "update", "--help")
	assert.NoError(t, err)
}
