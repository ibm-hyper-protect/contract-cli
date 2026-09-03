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
	"bytes"
	"encoding/json"
	"io"
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

// TestOutput_WithSingleFilePath tests Output when --out has 1 value;
// key files should default to sealed_decryption.pem / sealed_encryption.pem
// in the same directory.
func TestOutput_WithSingleFilePath(t *testing.T) {
	sealedSecretVal := "sealed-secret-data"
	decryptionKey := "-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----"
	verificationKey := "-----BEGIN PUBLIC KEY-----\ntest\n-----END PUBLIC KEY-----"

	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "test_output.txt")

	err := Output(sealedSecretVal, decryptionKey, verificationKey, outputFile)
	assert.NoError(t, err)

	sealedContent, readErr := os.ReadFile(outputFile)
	assert.NoError(t, readErr)
	assert.Contains(t, string(sealedContent), "sealed-secret-data")

	decryptionContent, readErr := os.ReadFile(filepath.Join(tmpDir, DefaultDecryptionKeyFileName))
	assert.NoError(t, readErr)
	assert.Contains(t, string(decryptionContent), "BEGIN PRIVATE KEY")

	verificationContent, readErr := os.ReadFile(filepath.Join(tmpDir, DefaultVerificationKeyFileName))
	assert.NoError(t, readErr)
	assert.Contains(t, string(verificationContent), "BEGIN PUBLIC KEY")
}

// TestOutput_WithThreeFilePaths tests Output when --out has 3 comma-separated values;
// each file should receive its designated content.
func TestOutput_WithThreeFilePaths(t *testing.T) {
	sealedSecretVal := "sealed-secret-data"
	decryptionKey := "-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----"
	verificationKey := "-----BEGIN PUBLIC KEY-----\ntest\n-----END PUBLIC KEY-----"

	tmpDir := t.TempDir()
	sealedPath := filepath.Join(tmpDir, "my_secret.txt")
	decryptPath := filepath.Join(tmpDir, "my_decrypt.pem")
	verifyPath := filepath.Join(tmpDir, "my_verify.pem")

	outFlag := sealedPath + "," + decryptPath + "," + verifyPath

	err := Output(sealedSecretVal, decryptionKey, verificationKey, outFlag)
	assert.NoError(t, err)

	sealedContent, readErr := os.ReadFile(sealedPath)
	assert.NoError(t, readErr)
	assert.Contains(t, string(sealedContent), "sealed-secret-data")

	decryptionContent, readErr := os.ReadFile(decryptPath)
	assert.NoError(t, readErr)
	assert.Contains(t, string(decryptionContent), "BEGIN PRIVATE KEY")

	verificationContent, readErr := os.ReadFile(verifyPath)
	assert.NoError(t, readErr)
	assert.Contains(t, string(verificationContent), "BEGIN PUBLIC KEY")
}

// TestOutput_WithoutFilePath tests Output function without file path (prints JSON to stdout)
func TestOutput_WithoutFilePath(t *testing.T) {
	sealedSecretVal := "sealed-secret-data"
	decryptionKey := "-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----"
	verificationKey := "-----BEGIN PUBLIC KEY-----\ntest\n-----END PUBLIC KEY-----"

	// Capture stdout
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := Output(sealedSecretVal, decryptionKey, verificationKey, "")

	w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)

	assert.NoError(t, err)

	// Verify the output is valid JSON with the expected fields
	var out SealedSecretOutput
	jsonErr := json.Unmarshal(buf.Bytes(), &out)
	assert.NoError(t, jsonErr)
	assert.Equal(t, sealedSecretVal, out.SealedSecret)
	assert.Contains(t, out.DecryptionKey, "BEGIN PRIVATE KEY")
	assert.Contains(t, out.VerificationKey, "BEGIN PUBLIC KEY")
}

// TestOutput_InvalidPath tests Output function with invalid file path
func TestOutput_InvalidPath(t *testing.T) {
	sealedSecret := "sealed-secret-data"
	decryptionKey := "test-key"
	verificationKey := "test-key"

	err := Output(sealedSecret, decryptionKey, verificationKey, testInvalidPath)
	assert.Error(t, err)
}

// TestOutput_InvalidOutCount tests Output function with wrong number of comma-separated values
func TestOutput_InvalidOutCount(t *testing.T) {
	sealedSecret := "sealed-secret-data"
	decryptionKey := "test-key"
	verificationKey := "test-key"

	// 2 values — should be rejected
	err := Output(sealedSecret, decryptionKey, verificationKey, "file1.txt,file2.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--out accepts 1 or 3 comma-separated values")

	// 4 values — should be rejected
	err = Output(sealedSecret, decryptionKey, verificationKey, "f1.txt,f2.txt,f3.txt,f4.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--out accepts 1 or 3 comma-separated values")
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
