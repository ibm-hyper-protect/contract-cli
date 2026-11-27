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
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/lib/plugin"
	"github.com/stretchr/testify/assert"
)

// TestDispatch_KnownCommandNotDelegated verifies that a built-in command (encrypt)
// is handled by Cobra, not dispatched to a plugin binary.
func TestDispatch_KnownCommandNotDelegated(t *testing.T) {
	pluginDir := setPluginDir(t)
	// Even if a binary named contract-plugin-encrypt exists, the built-in wins.
	makeExecutable(t, pluginDir, "encrypt")

	// The encrypt command will fail due to missing --in flag, but it should NOT be
	// dispatched to the plugin binary (which is a shell stub, not a real encrypt binary).
	_, err := runPluginCmd("encrypt", "--help")
	// --help on a known command always succeeds.
	assert.NoError(t, err)
}

// TestDispatch_UnknownCommandNoPlugin verifies the standard "unknown command" error
// is returned when no matching plugin binary exists.
func TestDispatch_UnknownCommandNoPlugin(t *testing.T) {
	setPluginDir(t)
	_, err := runPluginCmd("totallymadeupcommand")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown command")
}

// TestDispatch_FindPlugin verifies FindPlugin detects a plugin binary correctly.
// (The Exec call itself cannot be tested without replacing the process.)
func TestDispatch_FindPlugin_Present(t *testing.T) {
	pluginDir := setPluginDir(t)
	makeExecutable(t, pluginDir, "ccrv")

	path, found := plugin.FindPlugin("ccrv", pluginDir)
	assert.True(t, found)

	expectedBin := "contract-plugin-ccrv"
	if runtime.GOOS == "windows" {
		expectedBin += ".exe"
	}
	assert.Equal(t, filepath.Join(pluginDir, expectedBin), path)
}

// TestDispatch_FindPlugin_Absent verifies FindPlugin returns false for a missing plugin.
func TestDispatch_FindPlugin_Absent(t *testing.T) {
	pluginDir := setPluginDir(t)
	_, found := plugin.FindPlugin("doesnotexist", pluginDir)
	assert.False(t, found)
}

// TestDispatch_HandshakeEnvVarSet verifies the handshake constant values.
func TestDispatch_HandshakeEnvVarSet(t *testing.T) {
	assert.Equal(t, "CONTRACT_CLI_PLUGIN", plugin.HandshakeEnvVar)
	assert.Equal(t, "1", plugin.HandshakeEnvVal)
}

// TestDispatch_CheckHandshake_DirectExecution verifies CheckHandshake fails without env var.
func TestDispatch_CheckHandshake_DirectExecution(t *testing.T) {
	os.Unsetenv(plugin.HandshakeEnvVar)
	err := plugin.CheckHandshake()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contract-cli")
}
