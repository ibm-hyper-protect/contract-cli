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

package signContract

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testContractPath      = "../../samples/contract.yaml"
	testPrivateKeyPath    = "../../samples/sign/private.pem"
	testOutputPath        = "../../build/test_sign_contract_output.txt"
	testInvalidPath       = "../../build/file/file_not_exists.txt"
	testCorruptedContract = "../../build/corrupted_contract_sign.yaml"
	testCorruptedKey      = "../../build/corrupted_key_sign.pem"
	testEncryptedEnvPath  = "../../samples/sign/encrypted-env.txt"
	testEncryptedWlPath   = "../../samples/sign/encrypted-workload.txt"
)

// ─── ValidateInput tests ──────────────────────────────────────────────────────

// TestValidateInput_Success tests ValidateInput with --in and --priv
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(PrivateKeyFlagName, testPrivateKeyPath, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")
	cmd.Flags().String(PasswordFlagName, "", "")
	cmd.Flags().String(EnvFlagName, "", "")
	cmd.Flags().String(WorkloadFlagName, "", "")

	inputData, privateKeyPath, outputPath, password, envPath, workloadPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, inputData)
	assert.Equal(t, testPrivateKeyPath, privateKeyPath)
	assert.Equal(t, testOutputPath, outputPath)
	assert.Equal(t, "", password)
	assert.Equal(t, "", envPath)
	assert.Equal(t, "", workloadPath)
}

// TestValidateInput_WithoutOutputPath tests ValidateInput without output path
func TestValidateInput_WithoutOutputPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(PrivateKeyFlagName, testPrivateKeyPath, "")
	cmd.Flags().String(OutputFlagName, "", "")
	cmd.Flags().String(PasswordFlagName, "", "")
	cmd.Flags().String(EnvFlagName, "", "")
	cmd.Flags().String(WorkloadFlagName, "", "")

	inputData, privateKeyPath, outputPath, password, envPath, workloadPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, inputData)
	assert.Equal(t, testPrivateKeyPath, privateKeyPath)
	assert.Equal(t, "", outputPath)
	assert.Equal(t, "", password)
	assert.Equal(t, "", envPath)
	assert.Equal(t, "", workloadPath)
}

// TestValidateInput_WithPassword tests ValidateInput with password flag
func TestValidateInput_WithPassword(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(PrivateKeyFlagName, testPrivateKeyPath, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")
	cmd.Flags().String(PasswordFlagName, "mySecurePassword", "")
	cmd.Flags().String(EnvFlagName, "", "")
	cmd.Flags().String(WorkloadFlagName, "", "")

	inputData, privateKeyPath, outputPath, password, envPath, workloadPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, inputData)
	assert.Equal(t, testPrivateKeyPath, privateKeyPath)
	assert.Equal(t, testOutputPath, outputPath)
	assert.Equal(t, "mySecurePassword", password)
	assert.Equal(t, "", envPath)
	assert.Equal(t, "", workloadPath)
}

// TestValidateInput_EnvAndWorkload tests ValidateInput with --env and --workload (no --in)
func TestValidateInput_EnvAndWorkload(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, "", "")
	cmd.Flags().String(PrivateKeyFlagName, testPrivateKeyPath, "")
	cmd.Flags().String(OutputFlagName, "", "")
	cmd.Flags().String(PasswordFlagName, "", "")
	cmd.Flags().String(EnvFlagName, testEncryptedEnvPath, "")
	cmd.Flags().String(WorkloadFlagName, testEncryptedWlPath, "")

	inputData, privateKeyPath, outputPath, password, envPath, workloadPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, "", inputData)
	assert.Equal(t, testPrivateKeyPath, privateKeyPath)
	assert.Equal(t, "", outputPath)
	assert.Equal(t, "", password)
	assert.Equal(t, testEncryptedEnvPath, envPath)
	assert.Equal(t, testEncryptedWlPath, workloadPath)
}

// ─── BuildContractFromParts tests ────────────────────────────────────────────

// TestBuildContractFromParts_Success tests successful assembly
func TestBuildContractFromParts_Success(t *testing.T) {
	result, err := BuildContractFromParts(testEncryptedEnvPath, testEncryptedWlPath)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.True(t, strings.HasPrefix(result, "env: contract-basic."), "should start with env: contract-basic.")
	assert.Contains(t, result, "\nworkload: contract-basic.")
}

// TestBuildContractFromParts_InvalidEnvPath tests with missing env file
func TestBuildContractFromParts_InvalidEnvPath(t *testing.T) {
	result, err := BuildContractFromParts(testInvalidPath, testEncryptedWlPath)

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "env file path doesn't exist")
}

// TestBuildContractFromParts_InvalidWorkloadPath tests with missing workload file
func TestBuildContractFromParts_InvalidWorkloadPath(t *testing.T) {
	result, err := BuildContractFromParts(testEncryptedEnvPath, testInvalidPath)

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "workload file path doesn't exist")
}

// TestBuildContractFromParts_TrimsWhitespace tests that trailing newlines in files are trimmed
func TestBuildContractFromParts_TrimsWhitespace(t *testing.T) {
	envFile := "../../build/test_env_with_newline.txt"
	wlFile := "../../build/test_wl_with_newline.txt"
	defer os.Remove(envFile)
	defer os.Remove(wlFile)

	envBlob := "contract-basic.ENVBLOB=="
	wlBlob := "contract-basic.WLBLOB=="
	os.WriteFile(envFile, []byte(envBlob+"\n"), 0644)
	os.WriteFile(wlFile, []byte(wlBlob+"\n"), 0644)

	result, err := BuildContractFromParts(envFile, wlFile)

	assert.NoError(t, err)
	assert.Equal(t, "env: "+envBlob+"\nworkload: "+wlBlob+"\n", result)
}

// ─── GenerateSignContract tests ───────────────────────────────────────────────

// TestGenerateSignContract_Success tests successful contract signing via --in
func TestGenerateSignContract_Success(t *testing.T) {
	result, err := GenerateSignContract(testContractPath, testPrivateKeyPath, "", "", "")

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "envWorkloadSignature")
}

// TestGenerateSignContract_WithPassword tests signing with password parameter
func TestGenerateSignContract_WithPassword(t *testing.T) {
	result, err := GenerateSignContract(testContractPath, testPrivateKeyPath, "anyPassword", "", "")

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "envWorkloadSignature")
}

// TestGenerateSignContract_FromEnvWorkload tests signing via --env + --workload path
func TestGenerateSignContract_FromEnvWorkload(t *testing.T) {
	result, err := GenerateSignContract("", testPrivateKeyPath, "", testEncryptedEnvPath, testEncryptedWlPath)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "envWorkloadSignature")
	assert.Contains(t, result, "env:")
	assert.Contains(t, result, "workload:")
}

// TestGenerateSignContract_FromEnvWorkload_WithPassword tests the new path with a password param
func TestGenerateSignContract_FromEnvWorkload_WithPassword(t *testing.T) {
	result, err := GenerateSignContract("", testPrivateKeyPath, "anyPassword", testEncryptedEnvPath, testEncryptedWlPath)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "envWorkloadSignature")
}

// TestGenerateSignContract_InvalidContractPath tests with invalid --in path
func TestGenerateSignContract_InvalidContractPath(t *testing.T) {
	result, err := GenerateSignContract(testInvalidPath, testPrivateKeyPath, "", "", "")

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "doesn't exist")
}

// TestGenerateSignContract_InvalidPrivateKeyPath tests with invalid private key path
func TestGenerateSignContract_InvalidPrivateKeyPath(t *testing.T) {
	result, err := GenerateSignContract(testContractPath, testInvalidPath, "", "", "")

	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestGenerateSignContract_InvalidEnvPath tests new path with missing env file
func TestGenerateSignContract_InvalidEnvPath(t *testing.T) {
	result, err := GenerateSignContract("", testPrivateKeyPath, "", testInvalidPath, testEncryptedWlPath)

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "env file path doesn't exist")
}

// TestGenerateSignContract_InvalidWorkloadPath tests new path with missing workload file
func TestGenerateSignContract_InvalidWorkloadPath(t *testing.T) {
	result, err := GenerateSignContract("", testPrivateKeyPath, "", testEncryptedEnvPath, testInvalidPath)

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "workload file path doesn't exist")
}

// TestGenerateSignContract_CorruptedContract tests with corrupted --in contract file
func TestGenerateSignContract_CorruptedContract(t *testing.T) {
	err := os.WriteFile(testCorruptedContract, []byte("invalid: yaml: content: ["), 0644)
	assert.NoError(t, err)
	defer os.Remove(testCorruptedContract)

	result, err := GenerateSignContract(testCorruptedContract, testPrivateKeyPath, "", "", "")

	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestGenerateSignContract_CorruptedPrivateKey tests with corrupted private key
func TestGenerateSignContract_CorruptedPrivateKey(t *testing.T) {
	err := os.WriteFile(testCorruptedKey, []byte("not a valid private key"), 0644)
	assert.NoError(t, err)
	defer os.Remove(testCorruptedKey)

	result, err := GenerateSignContract(testContractPath, testCorruptedKey, "", "", "")

	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// ─── Output tests ─────────────────────────────────────────────────────────────

// TestOutput_WithFilePath tests Output function with valid file path
func TestOutput_WithFilePath(t *testing.T) {
	testData := "hyper-protect-basic.test-signed-data"
	os.Remove(testOutputPath)
	err := Output(testData, testOutputPath)
	assert.NoError(t, err)

	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Equal(t, testData, string(content))

	os.Remove(testOutputPath)
}

// TestOutput_WithoutFilePath tests Output function without file path (prints to stdout)
func TestOutput_WithoutFilePath(t *testing.T) {
	testData := "hyper-protect-basic.test-signed-data"
	err := Output(testData, "")
	assert.NoError(t, err)
}

// TestOutput_InvalidPath tests Output function with invalid file path
func TestOutput_InvalidPath(t *testing.T) {
	testData := "hyper-protect-basic.test-signed-data"
	err := Output(testData, testInvalidPath)
	assert.Error(t, err)
}
