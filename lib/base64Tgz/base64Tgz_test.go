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

package base64Tgz

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testInputPath       = "../../samples/tgz"
	testCertPath        = "../../samples/certificate/active.crt"
	testOutputPath      = "../../build/test_base64tgz_output.txt"
	testInvalidPath     = "../../build/file/file_not_exists.txt"
	testInvalidCertPath = "../../invalid_cert_path/path.crt"
)

// TestValidateInput_Success tests ValidateInput with valid input
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testInputPath, "")
	cmd.Flags().String(OutputFormatFlag, OutputFormatUnencrypted, "")
	cmd.Flags().String(OsVersionFlagName, "hpvs", "")
	cmd.Flags().String(CertFlagName, testCertPath, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")

	inputData, outputFormat, hyperProtectVersion, encCertPath, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testInputPath, inputData)
	assert.Equal(t, OutputFormatUnencrypted, outputFormat)
	assert.Equal(t, "hpvs", hyperProtectVersion)
	assert.Equal(t, testCertPath, encCertPath)
	assert.Equal(t, testOutputPath, outputPath)
}

// TestValidateInput_DefaultValues tests ValidateInput with default values
func TestValidateInput_DefaultValues(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testInputPath, "")
	cmd.Flags().String(OutputFormatFlag, DefaultOutput, "")
	cmd.Flags().String(OsVersionFlagName, "", "")
	cmd.Flags().String(CertFlagName, "", "")
	cmd.Flags().String(OutputFlagName, "", "")

	inputData, outputFormat, hyperProtectVersion, encCertPath, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testInputPath, inputData)
	assert.Equal(t, DefaultOutput, outputFormat)
	assert.Equal(t, "", hyperProtectVersion)
	assert.Equal(t, "", encCertPath)
	assert.Equal(t, "", outputPath)
}

// TestProcess_PlainSuccess tests Process function with plain output format
func TestProcess_PlainSuccess(t *testing.T) {
	result, err := Process(testInputPath, OutputFormatUnencrypted, "", "")

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestProcess_InvalidPath tests Process function with invalid input path
func TestProcess_InvalidPath(t *testing.T) {
	result, err := Process(testInvalidPath, OutputFormatUnencrypted, "", "")

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "not accessible")
}

// TestProcess_EncryptedSuccess tests Process function with encrypted output format
func TestProcess_EncryptedSuccess(t *testing.T) {
	result, err := Process(testInputPath, OutputFormatEncrypted, "hpvs", testCertPath)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestProcess_EncryptedInvalidCert tests Process function with invalid certificate path
func TestProcess_EncryptedInvalidCert(t *testing.T) {
	result, err := Process(testInputPath, OutputFormatEncrypted, "hpvs", testInvalidCertPath)

	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestProcess_InvalidOutputFormat tests Process function with invalid output format
func TestProcess_InvalidOutputFormat(t *testing.T) {
	result, err := Process(testInputPath, "invalid", "", "")

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "invalid output format")
}

// TestOutput_Success tests Output function with valid file path
func TestOutput_Success(t *testing.T) {
	testData := "dGVzdCBkYXRh"
	os.Remove(testOutputPath)
	err := Output(testData, testOutputPath)
	assert.NoError(t, err)

	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	os.Remove(testOutputPath)
}

// TestOutput_InvalidPath tests Output function with invalid file path
func TestOutput_InvalidPath(t *testing.T) {
	testData := "dGVzdCBkYXRh"
	err := Output(testData, testInvalidPath)
	assert.Error(t, err)
}

// TestOutput_EmptyPath tests Output function with empty output path
func TestOutput_EmptyPath(t *testing.T) {
	testData := "dGVzdCBkYXRh"
	err := Output(testData, "")
	assert.NoError(t, err)
}

// TestProcess_EncryptedEmptyCert tests Process function with encrypted format but empty cert path
func TestProcess_EncryptedEmptyCert(t *testing.T) {
	result, err := Process(testInputPath, OutputFormatEncrypted, "hpvs", "")
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}
