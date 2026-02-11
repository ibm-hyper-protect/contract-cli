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

package encryptString

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testInputText   = "hello world"
	testInputJson   = `{"key":"value"}`
	testInvalidJson = `{"key":"value"`
	testCertPath    = "../../samples/certificate/active.crt"
	testOutputPath  = "../../build/test_encrypt_string_output.txt"
	testInvalidPath = "../../build/file/file_not_exist.txt"
	testOsVersion   = "hpvs"
)

// TestValidateInput_Success tests ValidateInput with all required flags
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testInputText, "")
	cmd.Flags().String(FormatFlag, TextFormat, "")
	cmd.Flags().String(OsVersionFlagName, testOsVersion, "")
	cmd.Flags().String(CertFlagName, testCertPath, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")

	inputData, inputFormat, hyperProtectVersion, encCertPath, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testInputText, inputData)
	assert.Equal(t, TextFormat, inputFormat)
	assert.Equal(t, testOsVersion, hyperProtectVersion)
	assert.Equal(t, testCertPath, encCertPath)
	assert.Equal(t, testOutputPath, outputPath)
}

// TestValidateInput_WithJsonFormat tests ValidateInput with JSON format
func TestValidateInput_WithJsonFormat(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testInputJson, "")
	cmd.Flags().String(FormatFlag, JsonFormat, "")
	cmd.Flags().String(OsVersionFlagName, testOsVersion, "")
	cmd.Flags().String(CertFlagName, testCertPath, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")

	inputData, inputFormat, hyperProtectVersion, encCertPath, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testInputJson, inputData)
	assert.Equal(t, JsonFormat, inputFormat)
	assert.Equal(t, testOsVersion, hyperProtectVersion)
	assert.Equal(t, testCertPath, encCertPath)
	assert.Equal(t, testOutputPath, outputPath)
}

// TestValidateInput_WithoutOutputPath tests ValidateInput without output path (optional)
func TestValidateInput_WithoutOutputPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testInputText, "")
	cmd.Flags().String(FormatFlag, TextFormat, "")
	cmd.Flags().String(OsVersionFlagName, testOsVersion, "")
	cmd.Flags().String(CertFlagName, testCertPath, "")
	cmd.Flags().String(OutputFlagName, "", "")

	inputData, inputFormat, hyperProtectVersion, encCertPath, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testInputText, inputData)
	assert.Equal(t, TextFormat, inputFormat)
	assert.Equal(t, testOsVersion, hyperProtectVersion)
	assert.Equal(t, testCertPath, encCertPath)
	assert.Equal(t, "", outputPath)
}

// TestValidateInput_WithoutFlags tests ValidateInput when flags are not set
func TestValidateInput_WithoutFlags(t *testing.T) {
	cmd := &cobra.Command{}
	_, _, _, _, _, err := ValidateInput(cmd)
	assert.Error(t, err)
}

// TestProcess_TextFormat tests Process function with text format
func TestProcess_TextFormat(t *testing.T) {
	result, err := Process(testInputText, TextFormat, testOsVersion, testCertPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "hyper-protect-basic")
}

// TestProcess_JsonFormat tests Process function with JSON format
func TestProcess_JsonFormat(t *testing.T) {
	result, err := Process(testInputJson, JsonFormat, testOsVersion, testCertPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "hyper-protect-basic")
}

// TestProcess_InvalidFormat tests Process function with invalid format
func TestProcess_InvalidFormat(t *testing.T) {
	result, err := Process(testInputText, "invalid", testOsVersion, testCertPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "invalid input format")
}

// TestProcess_EmptyFormat tests Process function with empty format
func TestProcess_EmptyFormat(t *testing.T) {
	result, err := Process(testInputText, "", testOsVersion, testCertPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "invalid input format")
}

// TestProcess_InvalidCertPath tests Process function with invalid cert path
func TestProcess_InvalidCertPath(t *testing.T) {
	result, err := Process(testInputText, TextFormat, testOsVersion, testInvalidPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestProcess_EmptyCertPath tests Process function with empty cert path
func TestProcess_EmptyCertPath(t *testing.T) {
	result, err := Process(testInputText, TextFormat, testOsVersion, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestProcess_InvalidJson tests Process function with invalid JSON
func TestProcess_InvalidJson(t *testing.T) {
	result, err := Process(testInvalidJson, JsonFormat, testOsVersion, testCertPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestProcess_EmptyInput tests Process function with empty input
func TestProcess_EmptyInput(t *testing.T) {
	result, err := Process("", TextFormat, testOsVersion, testCertPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestProcess_InvalidOsVersion tests Process function with invalid OS version
func TestProcess_InvalidOsVersion(t *testing.T) {
	result, err := Process(testInputText, TextFormat, "invalid_os", testCertPath)
	// Invalid OS version doesn't fail, it uses default behavior
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestProcess_EmptyOsVersion tests Process function with empty OS version
func TestProcess_EmptyOsVersion(t *testing.T) {
	result, err := Process(testInputText, TextFormat, "", testCertPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestProcess_CorruptedCert tests Process function with corrupted certificate
func TestProcess_CorruptedCert(t *testing.T) {
	// Create a temporary file with corrupted cert
	corruptedCert := "../../build/corrupted_cert.crt"
	err := os.WriteFile(corruptedCert, []byte("not a valid certificate"), 0644)
	assert.NoError(t, err)
	defer os.Remove(corruptedCert)

	result, err := Process(testInputText, TextFormat, testOsVersion, corruptedCert)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestProcess_SpecialCharacters tests Process function with special characters
func TestProcess_SpecialCharacters(t *testing.T) {
	specialInput := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
	result, err := Process(specialInput, TextFormat, testOsVersion, testCertPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestProcess_DifferentOsVersions tests Process function with different OS versions
func TestProcess_DifferentOsVersions(t *testing.T) {
	// Test with hpcr-rhvs
	result, err := Process(testInputText, TextFormat, "hpcr-rhvs", testCertPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	// Test with hpcc-peerpod
	result, err = Process(testInputText, TextFormat, "hpcc-peerpod", testCertPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestOutput_ToFile tests writing encrypted string to file
func TestOutput_ToFile(t *testing.T) {
	testData := "hyper-protect-basic.encrypted.data"
	os.Remove(testOutputPath)
	err := Output(testOutputPath, testData)
	assert.NoError(t, err)

	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Equal(t, testData, string(content))
	os.Remove(testOutputPath)
}

// TestOutput_ToStdout tests printing to stdout (empty path)
func TestOutput_ToStdout(t *testing.T) {
	testData := "hyper-protect-basic.encrypted.data"
	err := Output("", testData)
	assert.NoError(t, err)
}

// TestOutput_InvalidPath tests with invalid output path
func TestOutput_InvalidPath(t *testing.T) {
	testData := "hyper-protect-basic.encrypted.data"
	err := Output(testInvalidPath, testData)
	assert.Error(t, err)
}

// TestValidateInput_EmptyOsVersion tests ValidateInput with empty OS version
func TestValidateInput_EmptyOsVersion(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testInputText, "")
	cmd.Flags().String(FormatFlag, TextFormat, "")
	cmd.Flags().String(OsVersionFlagName, "", "")
	cmd.Flags().String(CertFlagName, testCertPath, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")

	inputData, inputFormat, hyperProtectVersion, encCertPath, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testInputText, inputData)
	assert.Equal(t, TextFormat, inputFormat)
	assert.Equal(t, "", hyperProtectVersion)
	assert.Equal(t, testCertPath, encCertPath)
	assert.Equal(t, testOutputPath, outputPath)
}

// TestValidateInput_EmptyFormat tests ValidateInput with empty format
func TestValidateInput_EmptyFormat(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testInputText, "")
	cmd.Flags().String(FormatFlag, "", "")
	cmd.Flags().String(OsVersionFlagName, testOsVersion, "")
	cmd.Flags().String(CertFlagName, testCertPath, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")

	inputData, inputFormat, hyperProtectVersion, encCertPath, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testInputText, inputData)
	assert.Equal(t, "", inputFormat)
	assert.Equal(t, testOsVersion, hyperProtectVersion)
	assert.Equal(t, testCertPath, encCertPath)
	assert.Equal(t, testOutputPath, outputPath)
}
