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
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// makePluginBinary creates a minimal executable stub in tmpDir named contract-plugin-<name>.
func makePluginBinary(t *testing.T, dir, name string) string {
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

// ---------------------------------------------------------------------------
// ValidateInstallInput
// ---------------------------------------------------------------------------

func TestValidateInstallInput_Success(t *testing.T) {
	ref, err := ValidateInstallInput([]string{"github.com/org/contract-plugin-foo"})
	assert.NoError(t, err)
	assert.Equal(t, "github.com/org/contract-plugin-foo", ref)
}

func TestValidateInstallInput_Empty(t *testing.T) {
	_, err := ValidateInstallInput([]string{})
	assert.Error(t, err)
}

func TestValidateInstallInput_Whitespace(t *testing.T) {
	_, err := ValidateInstallInput([]string{"   "})
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// Uninstall
// ---------------------------------------------------------------------------

func TestUninstall_Success(t *testing.T) {
	tmpDir := t.TempDir()
	makePluginBinary(t, tmpDir, "foo")

	err := Uninstall("foo", tmpDir)
	assert.NoError(t, err)

	binName := "contract-plugin-foo"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	assert.NoFileExists(t, filepath.Join(tmpDir, binName))
}

func TestUninstall_NotInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	err := Uninstall("nonexistent", tmpDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not installed")
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

func TestList_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	infos, err := List(tmpDir)
	assert.NoError(t, err)
	assert.Empty(t, infos)
}

func TestList_WithPlugins(t *testing.T) {
	tmpDir := t.TempDir()
	makePluginBinary(t, tmpDir, "ccrv")
	makePluginBinary(t, tmpDir, "bar")

	infos, err := List(tmpDir)
	assert.NoError(t, err)
	assert.Len(t, infos, 2)

	names := make([]string, 0, len(infos))
	for _, i := range infos {
		names = append(names, i.Name)
	}
	assert.Contains(t, names, "ccrv")
	assert.Contains(t, names, "bar")
}

func TestList_InstallMethod_GitHubAndLocal(t *testing.T) {
	tmpDir := t.TempDir()
	makePluginBinary(t, tmpDir, "ccrv")
	makePluginBinary(t, tmpDir, "bar")
	// Only ccrv is tracked in .sources.json, i.e. installed from GitHub.
	err := writeSources(tmpDir, sourcesMap{"contract-plugin-ccrv": "github.com/ibm-hyper-protect/contract-plugin-ccrv"})
	assert.NoError(t, err)

	infos, err := List(tmpDir)
	assert.NoError(t, err)

	methods := map[string]string{}
	for _, i := range infos {
		methods[i.Name] = i.InstallMethod
	}
	assert.Equal(t, InstallMethodGitHub, methods["ccrv"])
	assert.Equal(t, InstallMethodLocal, methods["bar"])
}

func TestInstall_LocalOverGitHub_ClearsSource(t *testing.T) {
	pluginDir := t.TempDir()
	makePluginBinary(t, pluginDir, "foo")
	err := writeSources(pluginDir, sourcesMap{"contract-plugin-foo": "github.com/org/contract-plugin-foo"})
	assert.NoError(t, err)

	srcDir := t.TempDir()
	src := makePluginBinary(t, srcDir, "foo")

	assert.NoError(t, Install(src, pluginDir, false))

	infos, err := List(pluginDir)
	assert.NoError(t, err)
	assert.Len(t, infos, 1)
	assert.Equal(t, InstallMethodLocal, infos[0].InstallMethod)
}

// ---------------------------------------------------------------------------
// PluginVersion
// ---------------------------------------------------------------------------

func TestPluginVersion_NonExistentBinary(t *testing.T) {
	version := PluginVersion("/nonexistent/path/contract-plugin-fake")
	assert.Equal(t, "unknown", version)
}

func TestParseVersion_Variants_ReturnsVersionToken(t *testing.T) {
	assert.Equal(t, "v1.2.0", parseVersion("contract-plugin-ccrv version v1.2.0"))
	assert.Equal(t, "dev", parseVersion("contract-plugin-ccrv version dev"))
	assert.Equal(t, "v1.2.0", parseVersion("ccrv Version v1.2.0 (abc123)"))
	// No "version <token>" pair: returned unchanged.
	assert.Equal(t, "v1.2.0", parseVersion("v1.2.0"))
	assert.Equal(t, "unknown", parseVersion("unknown"))
}

func TestFormatList_Rows_HaveNoTrailingWhitespace(t *testing.T) {
	infos := []PluginInfo{
		{Name: "a", Version: "v1", InstallMethod: InstallMethodGitHub, Path: "/x"},
		{Name: "longer-name", Version: "v10.20.30", InstallMethod: InstallMethodLocal, Path: "/a/much/longer/path"},
	}
	for _, line := range strings.Split(strings.TrimRight(FormatList(infos), "\n"), "\n") {
		assert.Equal(t, strings.TrimRight(line, " "), line)
	}
}

// ---------------------------------------------------------------------------
// installFromLocal
// ---------------------------------------------------------------------------

func TestInstallFromLocal_Success(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	srcPath := filepath.Join(srcDir, "contract-plugin-foo")
	err := os.WriteFile(srcPath, []byte("fake"), 0755) // #nosec G306
	assert.NoError(t, err)

	err = installFromLocal(srcPath, destDir)
	assert.NoError(t, err)
	assert.FileExists(t, filepath.Join(destDir, "contract-plugin-foo"))
}

func TestInstallFromLocal_BadName(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	srcPath := filepath.Join(srcDir, "myplugin")       // missing prefix
	err := os.WriteFile(srcPath, []byte("fake"), 0755) // #nosec G306
	assert.NoError(t, err)

	err = installFromLocal(srcPath, destDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contract-plugin-")
}

func TestInstallFromLocal_PathNotExist(t *testing.T) {
	destDir := t.TempDir()
	err := installFromLocal("/nonexistent/contract-plugin-foo", destDir)
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// sourcesMap read/write
// ---------------------------------------------------------------------------

func TestReadWriteSources_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	m := sourcesMap{
		"contract-plugin-ccrv": "github.com/ibm-hyper-protect/contract-plugin-ccrv",
	}
	err := writeSources(tmpDir, m)
	assert.NoError(t, err)

	loaded, err := readSources(tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, m, loaded)
}

func TestReadSources_Missing(t *testing.T) {
	tmpDir := t.TempDir()
	m, err := readSources(tmpDir)
	assert.NoError(t, err)
	assert.Empty(t, m)
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

func TestUpdate_NoPluginsInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	err := Update("", tmpDir, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no plugins are installed")
}

func TestUpdate_LocalPluginSkipped(t *testing.T) {
	tmpDir := t.TempDir()
	makePluginBinary(t, tmpDir, "foo")
	// No entry in .sources.json — simulates a locally-installed plugin.
	err := Update("", tmpDir, false)
	// No error; local plugins are skipped with a warning.
	// But since all are local, nothing updated — no error returned for the loop itself.
	assert.NoError(t, err)
}

func TestUpdateOne_NotInSources(t *testing.T) {
	tmpDir := t.TempDir()
	err := updateOne("foo", tmpDir, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not installed from GitHub")
}
