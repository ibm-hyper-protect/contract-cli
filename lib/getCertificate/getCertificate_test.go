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

package getCertificate

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testCertsJsonPath = "../../samples/certificate/certs.json"
	testOutputPath    = "../../build/test_get_certificate_output.crt"
	testInvalidPath   = "../../build/file/file_not_exists.json"
	testVersion       = "1.0.23"
	testInvalidVer    = "9.9.99"
)

// TestValidateInput_Success tests ValidateInput with all required flags
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testCertsJsonPath, "")
	cmd.Flags().String(VersionFlagName, testVersion, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")

	encryptionCertsPath, version, encryptionCertificatePath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testCertsJsonPath, encryptionCertsPath)
	assert.Equal(t, testVersion, version)
	assert.Equal(t, testOutputPath, encryptionCertificatePath)
}

// TestValidateInput_WithoutOutputPath tests ValidateInput without output path (optional)
func TestValidateInput_WithoutOutputPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testCertsJsonPath, "")
	cmd.Flags().String(VersionFlagName, testVersion, "")
	cmd.Flags().String(OutputFlagName, "", "")

	encryptionCertsPath, version, encryptionCertificatePath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testCertsJsonPath, encryptionCertsPath)
	assert.Equal(t, testVersion, version)
	assert.Equal(t, "", encryptionCertificatePath)
}

// TestValidateInput_WithoutFlags tests ValidateInput when flags are not set
func TestValidateInput_WithoutFlags(t *testing.T) {
	cmd := &cobra.Command{}
	_, _, _, err := ValidateInput(cmd)
	assert.Error(t, err)
}

// TestProcess_Version3 tests extraction of version 1.0.23
func TestProcess_Version3(t *testing.T) {
	result, err := Process(testCertsJsonPath, testVersion)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "BEGIN CERTIFICATE")
}

// TestProcess_InvalidVersion tests with non-existent version
func TestProcess_InvalidVersion(t *testing.T) {
	result, err := Process(testCertsJsonPath, testInvalidVer)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestProcess_InvalidPath tests with non-existent file path
func TestProcess_InvalidPath(t *testing.T) {
	result, err := Process(testInvalidPath, testVersion)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "doesn't exist")
}

// TestProcess_EmptyPath tests with empty file path
func TestProcess_EmptyPath(t *testing.T) {
	result, err := Process("", testVersion)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "doesn't exist")
}

// TestProcess_EmptyVersion tests with empty version
func TestProcess_EmptyVersion(t *testing.T) {
	result, err := Process(testCertsJsonPath, "")
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestProcess_EmptyJson tests with empty JSON file
func TestProcess_EmptyJson(t *testing.T) {
	emptyFile := "../../build/empty_certs.json"
	err := os.WriteFile(emptyFile, []byte(""), 0644)
	assert.NoError(t, err)
	defer os.Remove(emptyFile)

	result, err := Process(emptyFile, testVersion)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestProcess_InvalidJsonStructure tests with valid JSON but wrong structure
func TestProcess_InvalidJsonStructure(t *testing.T) {
	invalidStructure := "../../build/invalid_structure.json"
	err := os.WriteFile(invalidStructure, []byte(`{"wrong": "structure"}`), 0644)
	assert.NoError(t, err)
	defer os.Remove(invalidStructure)

	result, err := Process(invalidStructure, testVersion)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestOutput_ToFile tests writing certificate to file
func TestOutput_ToFile(t *testing.T) {
	testCert := "-----BEGIN CERTIFICATE-----\ntest certificate data\n-----END CERTIFICATE-----"
	os.Remove(testOutputPath)
	err := Output(testCert, testOutputPath)
	assert.NoError(t, err)

	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Equal(t, testCert, string(content))
	os.Remove(testOutputPath)
}

// TestOutput_ToStdout tests printing to stdout (empty path)
func TestOutput_ToStdout(t *testing.T) {
	testCert := "-----BEGIN CERTIFICATE-----\ntest certificate data\n-----END CERTIFICATE-----"
	err := Output(testCert, "")
	assert.NoError(t, err)
}

// TestOutput_InvalidPath tests with invalid output path
func TestOutput_InvalidPath(t *testing.T) {
	testCert := "-----BEGIN CERTIFICATE-----\ntest certificate data\n-----END CERTIFICATE-----"
	err := Output(testCert, testInvalidPath)
	assert.Error(t, err)
}
