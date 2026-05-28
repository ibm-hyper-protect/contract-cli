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

package sealedSecret

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testSecretData     = "value123"
	testOutputPath     = "../../build/test_sealedsecret_output.txt"
	testInvalidPath    = "../../build/file/file_not_exists.txt"
	testSecretFilePath = "../../build/test_secret.txt"
)

// TestValidateInput_Success tests ValidateInput with all required flags
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testSecretData, "")
	cmd.Flags().String(TypeFlagName, "env", "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")
	cmd.Flags().String(EncryptionKeyFlagName, "", "")
	cmd.Flags().String(SigningKeyFlagName, "", "")

	inputData, secretType, outputPath, encKeyPath, signKeyPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testSecretData, inputData)
	assert.Equal(t, "env", secretType)
	assert.Equal(t, testOutputPath, outputPath)
	assert.Equal(t, "", encKeyPath)
	assert.Equal(t, "", signKeyPath)
}

// TestValidateInput_WorkloadType tests ValidateInput with workload type
func TestValidateInput_WorkloadType(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testSecretData, "")
	cmd.Flags().String(TypeFlagName, "workload", "")
	cmd.Flags().String(OutputFlagName, "", "")
	cmd.Flags().String(EncryptionKeyFlagName, "", "")
	cmd.Flags().String(SigningKeyFlagName, "", "")

	inputData, secretType, outputPath, encKeyPath, signKeyPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testSecretData, inputData)
	assert.Equal(t, "workload", secretType)
	assert.Equal(t, "", outputPath)
	assert.Equal(t, "", encKeyPath)
	assert.Equal(t, "", signKeyPath)
}

// TestValidateInput_WithAllOptionalFlags tests ValidateInput with all optional flags
func TestValidateInput_WithAllOptionalFlags(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testSecretData, "")
	cmd.Flags().String(TypeFlagName, "env", "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")
	cmd.Flags().String(EncryptionKeyFlagName, "/tmp/enc.key", "")
	cmd.Flags().String(SigningKeyFlagName, "/tmp/sign.key", "")

	inputData, secretType, outputPath, encKeyPath, signKeyPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testSecretData, inputData)
	assert.Equal(t, "env", secretType)
	assert.Equal(t, testOutputPath, outputPath)
	assert.Equal(t, "/tmp/enc.key", encKeyPath)
	assert.Equal(t, "/tmp/sign.key", signKeyPath)
}

// TestValidateInput_WithoutFlags tests ValidateInput when flags are not set
func TestValidateInput_WithoutFlags(t *testing.T) {
	cmd := &cobra.Command{}
	_, _, _, _, _, err := ValidateInput(cmd)
	assert.Error(t, err)
}

// TestGenerateSealedSecret_EnvType tests GenerateSealedSecret with env type
func TestGenerateSealedSecret_EnvType(t *testing.T) {
	sealedSecret, decryptionKey, verificationKey, inputSha, encryptedSha, err := GenerateSealedSecret(testSecretData, "env", "", "")
	assert.NoError(t, err)
	assert.NotEmpty(t, sealedSecret)
	assert.NotEmpty(t, decryptionKey)
	assert.NotEmpty(t, verificationKey)
	assert.NotEmpty(t, inputSha)
	assert.NotEmpty(t, encryptedSha)
}

// TestGenerateSealedSecret_WorkloadType tests GenerateSealedSecret with workload type
func TestGenerateSealedSecret_WorkloadType(t *testing.T) {
	sealedSecret, decryptionKey, verificationKey, _, _, err := GenerateSealedSecret("workload-secret-data", "workload", "", "")
	assert.NoError(t, err)
	assert.NotEmpty(t, sealedSecret)
	assert.NotEmpty(t, decryptionKey)
	assert.NotEmpty(t, verificationKey)
}

// TestGenerateSealedSecret_InvalidType tests GenerateSealedSecret with invalid type
func TestGenerateSealedSecret_InvalidType(t *testing.T) {
	sealedSecret, decryptionKey, verificationKey, inputSha, encryptedSha, err := GenerateSealedSecret(testSecretData, "invalid", "", "")
	assert.Error(t, err)
	assert.Empty(t, sealedSecret)
	assert.Empty(t, decryptionKey)
	assert.Empty(t, verificationKey)
	assert.Empty(t, inputSha)
	assert.Empty(t, encryptedSha)
	assert.Contains(t, err.Error(), "invalid secret type")
}

// TestGenerateSealedSecret_WithFile tests GenerateSealedSecret with file input
func TestGenerateSealedSecret_WithFile(t *testing.T) {
	// Create a temporary file with test data
	tmpDir := "../../build"
	os.MkdirAll(tmpDir, 0755)
	testFile := filepath.Join(tmpDir, "test_secret.txt")
	testData := "test123"

	err := os.WriteFile(testFile, []byte(testData), 0644)
	assert.NoError(t, err)
	defer os.Remove(testFile)

	sealedSecret, decryptionKey, verificationKey, _, _, err := GenerateSealedSecret(testFile, "env", "", "")
	assert.NoError(t, err)
	assert.NotEmpty(t, sealedSecret)
	assert.NotEmpty(t, decryptionKey)
	assert.NotEmpty(t, verificationKey)
}

// TestGenerateSealedSecret_InvalidFilePath tests GenerateSealedSecret with invalid file path
// Note: Invalid file paths are treated as direct string input, so this will succeed
func TestGenerateSealedSecret_InvalidFilePath(t *testing.T) {
	sealedSecret, decryptionKey, verificationKey, _, _, err := GenerateSealedSecret(testInvalidPath, "env", "", "")
	assert.NoError(t, err)
	assert.NotEmpty(t, sealedSecret)
	assert.NotEmpty(t, decryptionKey)
	assert.NotEmpty(t, verificationKey)
}

// TestGenerateSealedSecret_InvalidEncryptionKeyPath tests with invalid encryption key path
func TestGenerateSealedSecret_InvalidEncryptionKeyPath(t *testing.T) {
	sealedSecret, decryptionKey, verificationKey, inputSha, encryptedSha, err := GenerateSealedSecret(testSecretData, "env", testInvalidPath, "")
	assert.Error(t, err)
	assert.Empty(t, sealedSecret)
	assert.Empty(t, decryptionKey)
	assert.Empty(t, verificationKey)
	assert.Empty(t, inputSha)
	assert.Empty(t, encryptedSha)
}

// TestGenerateSealedSecret_InvalidSigningKeyPath tests with invalid signing key path
func TestGenerateSealedSecret_InvalidSigningKeyPath(t *testing.T) {
	sealedSecret, decryptionKey, verificationKey, inputSha, encryptedSha, err := GenerateSealedSecret(testSecretData, "env", "", testInvalidPath)
	assert.Error(t, err)
	assert.Empty(t, sealedSecret)
	assert.Empty(t, decryptionKey)
	assert.Empty(t, verificationKey)
	assert.Empty(t, inputSha)
	assert.Empty(t, encryptedSha)
}

// TestOutput_WithFilePath tests Output function with valid file path
func TestOutput_WithFilePath(t *testing.T) {
	sealedSecret := "sealed-secret-data"
	decryptionKey := "-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----"
	verificationKey := "-----BEGIN PUBLIC KEY-----\ntest\n-----END PUBLIC KEY-----"

	tmpDir := "../../build"
	os.MkdirAll(tmpDir, 0755)
	outputFile := filepath.Join(tmpDir, "test_output.txt")
	os.Remove(outputFile)

	err := Output(sealedSecret, decryptionKey, verificationKey, outputFile)
	assert.NoError(t, err)

	_, statErr := os.Stat(outputFile)
	assert.NoError(t, statErr)

	content, readErr := os.ReadFile(outputFile)
	assert.NoError(t, readErr)
	assert.Contains(t, string(content), "Sealed Secret:")
	assert.Contains(t, string(content), "SECRET_DECRYPTION_KEY=")
	assert.Contains(t, string(content), "SECRET_VERIFICATION_KEY=")

	os.Remove(outputFile)
}

// TestOutput_WithoutFilePath tests Output function without file path (prints to stdout)
func TestOutput_WithoutFilePath(t *testing.T) {
	sealedSecret := "sealed-secret-data"
	decryptionKey := "-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----"
	verificationKey := "-----BEGIN PUBLIC KEY-----\ntest\n-----END PUBLIC KEY-----"

	err := Output(sealedSecret, decryptionKey, verificationKey, "")
	assert.NoError(t, err)
}

// TestOutput_InvalidPath tests Output function with invalid file path
func TestOutput_InvalidPath(t *testing.T) {
	sealedSecret := "sealed-secret-data"
	decryptionKey := "test-key"
	verificationKey := "test-key"

	err := Output(sealedSecret, decryptionKey, verificationKey, testInvalidPath)
	assert.Error(t, err)
}

// TestFormatKeyWithEscapedNewlines_WithNewlines tests formatKeyWithEscapedNewlines with newlines
func TestFormatKeyWithEscapedNewlines_WithNewlines(t *testing.T) {
	input := "-----BEGIN KEY-----\nline1\nline2\n-----END KEY-----"
	expected := "-----BEGIN KEY-----\\nline1\\nline2\\n-----END KEY-----"
	result := formatKeyWithEscapedNewlines(input)
	assert.Equal(t, expected, result)
}

// TestFormatKeyWithEscapedNewlines_WithoutNewlines tests formatKeyWithEscapedNewlines without newlines
func TestFormatKeyWithEscapedNewlines_WithoutNewlines(t *testing.T) {
	input := "single-line-key"
	expected := "single-line-key"
	result := formatKeyWithEscapedNewlines(input)
	assert.Equal(t, expected, result)
}

// TestFormatKeyWithEscapedNewlines_EmptyString tests formatKeyWithEscapedNewlines with empty string
func TestFormatKeyWithEscapedNewlines_EmptyString(t *testing.T) {
	input := ""
	expected := ""
	result := formatKeyWithEscapedNewlines(input)
	assert.Equal(t, expected, result)
}

// TestGenerateSealedSecret_EmptySecret tests GenerateSealedSecret with empty secret
func TestGenerateSealedSecret_EmptySecret(t *testing.T) {
	sealedSecret, decryptionKey, verificationKey, inputSha, encryptedSha, err := GenerateSealedSecret("", "env", "", "")
	assert.Error(t, err)
	assert.Empty(t, sealedSecret)
	assert.Empty(t, decryptionKey)
	assert.Empty(t, verificationKey)
	assert.Empty(t, inputSha)
	assert.Empty(t, encryptedSha)
	assert.Contains(t, err.Error(), "secret cannot be empty")
}

// TestGenerateSealedSecret_MultilineSecret tests GenerateSealedSecret with multiline secret data
func TestGenerateSealedSecret_MultilineSecret(t *testing.T) {
	multilineSecret := "SECRET_VAR1=value1\nSECRET_VAR2=value2\nSECRET_VAR3=value3"

	sealedSecret, decryptionKey, verificationKey, _, _, err := GenerateSealedSecret(multilineSecret, "env", "", "")
	assert.NoError(t, err)
	assert.NotEmpty(t, sealedSecret)
	assert.NotEmpty(t, decryptionKey)
	assert.NotEmpty(t, verificationKey)
}
