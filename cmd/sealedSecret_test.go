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
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/lib/sealedSecret"
	"github.com/stretchr/testify/assert"
)

const (
	sampleSealedSecretInput = "value123"
)

// TestSealedSecretCmd_EnvType tests if sealed secret command can generate env type sealed secret
func TestSealedSecretCmd_EnvType(t *testing.T) {
	// Create unique output file for this test
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "sealed-secret-env-output.txt")

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{sealedSecret.ParameterName, "--in", sampleSealedSecretInput, "--type", "env", "--out", outputPath})
	err := sealedSecretCmd.Execute()

	assert.NoError(t, err)

	// Verify output file was created and contains expected content
	content, readErr := os.ReadFile(outputPath)
	assert.NoError(t, readErr)
	assert.Contains(t, string(content), "Sealed Secret:")
	assert.Contains(t, string(content), "SECRET_DECRYPTION_KEY=")
	assert.Contains(t, string(content), "SECRET_VERIFICATION_KEY=")
	// t.TempDir() automatically cleans up after test
}

// TestSealedSecretCmd_WorkloadType tests if sealed secret command can generate workload type sealed secret
func TestSealedSecretCmd_WorkloadType(t *testing.T) {
	// Create unique output file for this test
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "sealed-secret-workload-output.txt")

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{sealedSecret.ParameterName, "--in", sampleSealedSecretInput, "--type", "workload", "--out", outputPath})
	err := sealedSecretCmd.Execute()

	assert.NoError(t, err)

	// Verify output file was created and contains expected content
	content, readErr := os.ReadFile(outputPath)
	assert.NoError(t, readErr)
	assert.Contains(t, string(content), "Sealed Secret:")
	assert.Contains(t, string(content), "SECRET_DECRYPTION_KEY=")
	assert.Contains(t, string(content), "SECRET_VERIFICATION_KEY=")
	// t.TempDir() automatically cleans up after test
}

// TestSealedSecretCmd_CommandProperties tests command properties
func TestSealedSecretCmd_CommandProperties(t *testing.T) {
	assert.NotNil(t, sealedSecretCmd)
	assert.Equal(t, sealedSecret.ParameterName, sealedSecretCmd.Use)
	assert.Equal(t, sealedSecret.ParameterShortDescription, sealedSecretCmd.Short)
	assert.Equal(t, sealedSecret.ParameterLongDescription, sealedSecretCmd.Long)
}

// TestSealedSecretCmd_Flags tests if all required flags are present
func TestSealedSecretCmd_Flags(t *testing.T) {
	flags := []string{
		sealedSecret.InputFlagName,
		sealedSecret.TypeFlagName,
		sealedSecret.OutputFlagName,
		sealedSecret.EncryptionKeyFlagName,
		sealedSecret.SigningKeyFlagName,
	}

	for _, flagName := range flags {
		flag := sealedSecretCmd.Flags().Lookup(flagName)
		assert.NotNil(t, flag, "Expected flag '%s' to be present", flagName)
	}
}

// TestSealedSecretCmd_FlagShorthands tests that flags don't have shorthands (using PersistentFlags)
func TestSealedSecretCmd_FlagShorthands(t *testing.T) {
	inputFlag := sealedSecretCmd.Flags().Lookup(sealedSecret.InputFlagName)
	assert.NotNil(t, inputFlag)
	assert.Equal(t, "", inputFlag.Shorthand)

	typeFlag := sealedSecretCmd.Flags().Lookup(sealedSecret.TypeFlagName)
	assert.NotNil(t, typeFlag)
	assert.Equal(t, "", typeFlag.Shorthand)

	outputFlag := sealedSecretCmd.Flags().Lookup(sealedSecret.OutputFlagName)
	assert.NotNil(t, outputFlag)
	assert.Equal(t, "", outputFlag.Shorthand)

	encKeyFlag := sealedSecretCmd.Flags().Lookup(sealedSecret.EncryptionKeyFlagName)
	assert.NotNil(t, encKeyFlag)
	assert.Equal(t, "", encKeyFlag.Shorthand)

	signKeyFlag := sealedSecretCmd.Flags().Lookup(sealedSecret.SigningKeyFlagName)
	assert.NotNil(t, signKeyFlag)
	assert.Equal(t, "", signKeyFlag.Shorthand)
}

// TestSealedSecretCmd_FlagDefaults tests flag default values
func TestSealedSecretCmd_FlagDefaults(t *testing.T) {
	inputFlag := sealedSecretCmd.Flags().Lookup(sealedSecret.InputFlagName)
	assert.NotNil(t, inputFlag)
	assert.Equal(t, "", inputFlag.DefValue)

	typeFlag := sealedSecretCmd.Flags().Lookup(sealedSecret.TypeFlagName)
	assert.NotNil(t, typeFlag)
	assert.Equal(t, "", typeFlag.DefValue)

	outputFlag := sealedSecretCmd.Flags().Lookup(sealedSecret.OutputFlagName)
	assert.NotNil(t, outputFlag)
	assert.Equal(t, "", outputFlag.DefValue)

	encKeyFlag := sealedSecretCmd.Flags().Lookup(sealedSecret.EncryptionKeyFlagName)
	assert.NotNil(t, encKeyFlag)
	assert.Equal(t, "", encKeyFlag.DefValue)

	signKeyFlag := sealedSecretCmd.Flags().Lookup(sealedSecret.SigningKeyFlagName)
	assert.NotNil(t, signKeyFlag)
	assert.Equal(t, "", signKeyFlag.DefValue)
}
