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
	"testing"

	cryptoLib "github.com/ibm-hyper-protect/contract-cli/lib/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// silenceStdout redirects os.Stdout to a temp file for the duration of fn
// so that fmt.Printf output from the command does not pollute test logs.
func silenceStdout(t *testing.T, fn func()) {
	t.Helper()
	tmpFile, err := os.CreateTemp(t.TempDir(), "stdout-*")
	require.NoError(t, err)

	orig := os.Stdout
	os.Stdout = tmpFile
	defer func() {
		os.Stdout = orig
		tmpFile.Close()
	}()

	fn()
}

// TestCryptoCmd_CommandProperties verifies the command metadata.
func TestCryptoCmd_CommandProperties(t *testing.T) {
	assert.NotNil(t, cryptoCmd)
	assert.Equal(t, cryptoLib.ParameterName, cryptoCmd.Use)
	assert.Equal(t, cryptoLib.ParameterShortDescription, cryptoCmd.Short)
	assert.Equal(t, cryptoLib.ParameterLongDescription, cryptoCmd.Long)
}

// TestCryptoCmd_Flags verifies all flags are registered.
// Uses cmd.Flag() which searches both local and persistent flag sets.
func TestCryptoCmd_Flags(t *testing.T) {
	flags := []string{
		cryptoLib.TypeFlagName,
		cryptoLib.OutFlagName,
		cryptoLib.SizeFlagName,
		cryptoLib.PasswordFlagName,
		cryptoLib.SANFlagName,
		cryptoLib.DaysFlagName,
		cryptoLib.CNFlagName,
	}
	for _, name := range flags {
		// cmd.Flag() searches persistent + local; Flags().Lookup() misses PersistentFlags.
		f := cryptoCmd.Flag(name)
		assert.NotNil(t, f, "expected flag '%s' to be registered", name)
	}
}

// TestCryptoCmd_FlagDefaults verifies the cobra-level default values.
// Uses cmd.Flag() which searches both local and persistent flag sets.
func TestCryptoCmd_FlagDefaults(t *testing.T) {
	stringFlags := []string{
		cryptoLib.TypeFlagName,
		cryptoLib.OutFlagName,
		cryptoLib.PasswordFlagName,
		cryptoLib.SANFlagName,
		cryptoLib.CNFlagName,
	}
	for _, name := range stringFlags {
		f := cryptoCmd.Flag(name)
		require.NotNil(t, f)
		assert.Equal(t, "", f.DefValue, "flag --%s should default to empty string at cobra level", name)
	}

	// Int flags default to 0 at cobra level; ValidateInput applies real defaults.
	for _, name := range []string{cryptoLib.SizeFlagName, cryptoLib.DaysFlagName} {
		f := cryptoCmd.Flag(name)
		require.NotNil(t, f)
		assert.Equal(t, "0", f.DefValue, "flag --%s should default to 0 at cobra level", name)
	}
}

// TestCryptoCmd_KeyTypeToFile verifies end-to-end execution for --type key
// writing to two explicitly named output files.
func TestCryptoCmd_KeyTypeToFile(t *testing.T) {
	tmpDir := t.TempDir()
	privPath := filepath.Join(tmpDir, "testkey_private.pem")
	pubPath := filepath.Join(tmpDir, "testkey_public.pem")
	outVal := privPath + "," + pubPath

	rootCmd.SetArgs([]string{
		cryptoLib.ParameterName,
		"--" + cryptoLib.TypeFlagName, "key",
		"--" + cryptoLib.SizeFlagName, "2048",
		"--" + cryptoLib.OutFlagName, outVal,
	})
	silenceStdout(t, func() {
		require.NoError(t, cryptoCmd.Execute())
	})

	assert.FileExists(t, privPath)
	assert.FileExists(t, pubPath)

	privData, err := os.ReadFile(privPath)
	require.NoError(t, err)
	assert.Contains(t, string(privData), "PRIVATE KEY")

	info, err := os.Stat(privPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

// TestCryptoCmd_CertTypeToFile verifies end-to-end execution for --type cert
// writing to three explicitly named output files.
func TestCryptoCmd_CertTypeToFile(t *testing.T) {
	tmpDir := t.TempDir()
	caPath := filepath.Join(tmpDir, "testcert-ca.crt")
	clientCertPath := filepath.Join(tmpDir, "testcert-client.crt")
	clientKeyPath := filepath.Join(tmpDir, "testcert-client.pem")
	outVal := caPath + "," + clientCertPath + "," + clientKeyPath

	rootCmd.SetArgs([]string{
		cryptoLib.ParameterName,
		"--" + cryptoLib.TypeFlagName, "cert",
		"--" + cryptoLib.SizeFlagName, "2048",
		"--" + cryptoLib.SANFlagName, "example.com",
		"--" + cryptoLib.OutFlagName, outVal,
	})
	silenceStdout(t, func() {
		require.NoError(t, cryptoCmd.Execute())
	})

	assert.FileExists(t, caPath)
	assert.FileExists(t, clientCertPath)
	assert.FileExists(t, clientKeyPath)

	caData, err := os.ReadFile(caPath)
	require.NoError(t, err)
	assert.Contains(t, string(caData), "BEGIN CERTIFICATE")

	keyInfo, err := os.Stat(clientKeyPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), keyInfo.Mode().Perm())
}

// TestCryptoCmd_KeyTypeWithPassword verifies that a password-protected
// private key file is encrypted on disk when two output names are given.
func TestCryptoCmd_KeyTypeWithPassword(t *testing.T) {
	tmpDir := t.TempDir()
	privPath := filepath.Join(tmpDir, "enckey_private.pem")
	pubPath := filepath.Join(tmpDir, "enckey_public.pem")
	outVal := privPath + "," + pubPath

	rootCmd.SetArgs([]string{
		cryptoLib.ParameterName,
		"--" + cryptoLib.TypeFlagName, "key",
		"--" + cryptoLib.SizeFlagName, "2048",
		"--" + cryptoLib.PasswordFlagName, "test-passphrase-cmd",
		"--" + cryptoLib.OutFlagName, outVal,
	})
	silenceStdout(t, func() {
		require.NoError(t, cryptoCmd.Execute())
	})

	privData, err := os.ReadFile(privPath)
	require.NoError(t, err)
	assert.Contains(t, string(privData), "ENCRYPTED",
		"private key on disk must be encrypted when --password is provided")
}
