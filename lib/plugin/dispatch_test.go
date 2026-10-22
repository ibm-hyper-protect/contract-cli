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
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// CheckHandshake
// ---------------------------------------------------------------------------

func TestCheckHandshake_WithEnvVar(t *testing.T) {
	t.Setenv(HandshakeEnvVar, HandshakeEnvVal)
	err := CheckHandshake()
	assert.NoError(t, err)
}

func TestCheckHandshake_MissingEnvVar(t *testing.T) {
	os.Unsetenv(HandshakeEnvVar)
	err := CheckHandshake()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contract-cli")
}

func TestCheckHandshake_WrongValue(t *testing.T) {
	t.Setenv(HandshakeEnvVar, "0")
	err := CheckHandshake()
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// FindPlugin
// ---------------------------------------------------------------------------

func TestFindPlugin_Found(t *testing.T) {
	tmpDir := t.TempDir()
	binName := "contract-plugin-ccrv"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	binPath := filepath.Join(tmpDir, binName)
	err := os.WriteFile(binPath, []byte("stub"), 0755) // #nosec G306
	assert.NoError(t, err)

	path, found := FindPlugin("ccrv", tmpDir)
	assert.True(t, found)
	assert.Equal(t, binPath, path)
}

func TestFindPlugin_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	_, found := FindPlugin("doesnotexist", tmpDir)
	assert.False(t, found)
}

func TestFindPlugin_NotExecutable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("executable bit not applicable on Windows")
	}
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "contract-plugin-foo")
	err := os.WriteFile(binPath, []byte("stub"), 0600) // not executable
	assert.NoError(t, err)

	_, found := FindPlugin("foo", tmpDir)
	assert.False(t, found)
}

// ---------------------------------------------------------------------------
// pluginNameFromArg
// ---------------------------------------------------------------------------

func TestPluginNameFromArg_PlainName(t *testing.T) {
	assert.Equal(t, "ccrv", pluginNameFromArg("ccrv"))
}

func TestPluginNameFromArg_WithPrefix(t *testing.T) {
	assert.Equal(t, "ccrv", pluginNameFromArg("contract-plugin-ccrv"))
}
