// Copyright (c) 2025 IBM Corp.
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
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/lib/signContract"
	"github.com/stretchr/testify/assert"
)

const (
	testContractPath   = "../samples/contract.yaml"
	testPrivateKeyPath = "../samples/sign/private.pem"
	testOutputPath     = "../build/test_cmd_sign_contract_output.txt"
	testInvalidPath    = "../build/file/file_not_exists.txt"
)

// TestSignContractCmd_Success tests successful contract signing via command
func TestSignContractCmd_Success(t *testing.T) {
	// Clean up any existing output file
	os.Remove(testOutputPath)

	// Execute the sign contract command
	signContractCmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testContractPath,
		"--" + signContract.PrivateKeyFlagName, testPrivateKeyPath,
		"--" + signContract.OutputFlagName, testOutputPath,
	})

	err := signContractCmd.Execute()
	assert.NoError(t, err)

	// Verify output file was created
	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	// Verify output contains signed contract
	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Contains(t, string(content), "hyper-protect-basic")

	// Clean up
	os.Remove(testOutputPath)
}

// TestSignContractCmd_WithPassword tests signing with password parameter
func TestSignContractCmd_WithPassword(t *testing.T) {
	os.Remove(testOutputPath)

	signContractCmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testContractPath,
		"--" + signContract.PrivateKeyFlagName, testPrivateKeyPath,
		"--" + signContract.PasswordFlagName, "testPassword123",
		"--" + signContract.OutputFlagName, testOutputPath,
	})

	err := signContractCmd.Execute()
	assert.NoError(t, err)

	// Verify output file was created
	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	// Verify output contains signed contract
	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Contains(t, string(content), "hyper-protect-basic")

	os.Remove(testOutputPath)
}

// TestSignContractCmd_WithEmptyPassword tests signing with empty password
func TestSignContractCmd_WithEmptyPassword(t *testing.T) {
	os.Remove(testOutputPath)

	signContractCmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testContractPath,
		"--" + signContract.PrivateKeyFlagName, testPrivateKeyPath,
		"--" + signContract.PasswordFlagName, "",
		"--" + signContract.OutputFlagName, testOutputPath,
	})

	err := signContractCmd.Execute()
	assert.NoError(t, err)

	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	os.Remove(testOutputPath)
}

// TestSignContractCmd_WithoutOutputPath tests signing without output path (stdout)
func TestSignContractCmd_WithoutOutputPath(t *testing.T) {
	signContractCmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testContractPath,
		"--" + signContract.PrivateKeyFlagName, testPrivateKeyPath,
	})

	err := signContractCmd.Execute()
	assert.NoError(t, err)
}

// TestSignContractCmd_MissingInputFlag tests error when input flag is missing
func TestSignContractCmd_MissingInputFlag(t *testing.T) {
	signContractCmd.SetArgs([]string{
		"--" + signContract.PrivateKeyFlagName, testPrivateKeyPath,
	})

	err := signContractCmd.Execute()
	assert.Error(t, err)
}

// TestSignContractCmd_MissingPrivateKeyFlag tests error when private key flag is missing
func TestSignContractCmd_MissingPrivateKeyFlag(t *testing.T) {
	signContractCmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testContractPath,
	})

	err := signContractCmd.Execute()
	assert.Error(t, err)
}

// TestSignContractCmd_InvalidContractPath tests error with invalid contract path
func TestSignContractCmd_InvalidContractPath(t *testing.T) {
	signContractCmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testInvalidPath,
		"--" + signContract.PrivateKeyFlagName, testPrivateKeyPath,
	})

	err := signContractCmd.Execute()
	assert.Error(t, err)
}

// TestSignContractCmd_InvalidPrivateKeyPath tests error with invalid private key path
func TestSignContractCmd_InvalidPrivateKeyPath(t *testing.T) {
	signContractCmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testContractPath,
		"--" + signContract.PrivateKeyFlagName, testInvalidPath,
	})

	err := signContractCmd.Execute()
	assert.Error(t, err)
}

// TestSignContractCmd_CorruptedContract tests error with corrupted contract
func TestSignContractCmd_CorruptedContract(t *testing.T) {
	corruptedFile := "../build/corrupted_contract_cmd.yaml"
	err := os.WriteFile(corruptedFile, []byte("invalid: yaml: content: ["), 0644)
	assert.NoError(t, err)
	defer os.Remove(corruptedFile)

	signContractCmd.SetArgs([]string{
		"--" + signContract.InputFlagName, corruptedFile,
		"--" + signContract.PrivateKeyFlagName, testPrivateKeyPath,
	})

	err = signContractCmd.Execute()
	assert.Error(t, err)
}

// TestSignContractCmd_CorruptedPrivateKey tests error with corrupted private key
func TestSignContractCmd_CorruptedPrivateKey(t *testing.T) {
	corruptedKey := "../build/corrupted_key_cmd.pem"
	err := os.WriteFile(corruptedKey, []byte("not a valid private key"), 0644)
	assert.NoError(t, err)
	defer os.Remove(corruptedKey)

	signContractCmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testContractPath,
		"--" + signContract.PrivateKeyFlagName, corruptedKey,
	})

	err = signContractCmd.Execute()
	assert.Error(t, err)
}

// TestSignContractCmd_WithPasswordAndOutput tests complete workflow with password and output
func TestSignContractCmd_WithPasswordAndOutput(t *testing.T) {
	os.Remove(testOutputPath)

	signContractCmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testContractPath,
		"--" + signContract.PrivateKeyFlagName, testPrivateKeyPath,
		"--" + signContract.PasswordFlagName, "securePass456",
		"--" + signContract.OutputFlagName, testOutputPath,
	})

	err := signContractCmd.Execute()
	assert.NoError(t, err)

	// Verify output
	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.NotEmpty(t, content)
	assert.Contains(t, string(content), "hyper-protect-basic")

	os.Remove(testOutputPath)
}

// TestSignContractCmd_EmptyContract tests error with empty contract file
func TestSignContractCmd_EmptyContract(t *testing.T) {
	emptyFile := "../build/empty_contract_cmd.yaml"
	err := os.WriteFile(emptyFile, []byte(""), 0644)
	assert.NoError(t, err)
	defer os.Remove(emptyFile)

	signContractCmd.SetArgs([]string{
		"--" + signContract.InputFlagName, emptyFile,
		"--" + signContract.PrivateKeyFlagName, testPrivateKeyPath,
	})

	err = signContractCmd.Execute()
	assert.Error(t, err)
}
