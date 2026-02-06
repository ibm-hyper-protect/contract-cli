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
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testCertsJsonPath = "../../samples/certificate/certs.json"
	testOutputPath    = "../../build/test_get_certificate_output.crt"
	testInvalidPath   = "../../build/file/file_not_exists.json"
	testVersion1      = "1.0.21"
	testVersion2      = "1.0.22"
	testVersion3      = "1.0.23"
	testInvalidVer    = "9.9.99"
)

// TestValidateInput_Success tests ValidateInput with all required flags
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testCertsJsonPath, "")
	cmd.Flags().String(VersionFlagName, testVersion2, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")

	encryptionCertsPath, version, encryptionCertificatePath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testCertsJsonPath, encryptionCertsPath)
	assert.Equal(t, testVersion2, version)
	assert.Equal(t, testOutputPath, encryptionCertificatePath)
}

// TestValidateInput_WithoutOutputPath tests ValidateInput without output path (optional)
func TestValidateInput_WithoutOutputPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testCertsJsonPath, "")
	cmd.Flags().String(VersionFlagName, testVersion2, "")
	cmd.Flags().String(OutputFlagName, "", "")

	encryptionCertsPath, version, encryptionCertificatePath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testCertsJsonPath, encryptionCertsPath)
	assert.Equal(t, testVersion2, version)
	assert.Equal(t, "", encryptionCertificatePath)
}

// TestValidateInput_DifferentVersions tests ValidateInput with different version numbers
func TestValidateInput_DifferentVersions(t *testing.T) {
	versions := []string{testVersion1, testVersion2, testVersion3}

	for _, ver := range versions {
		cmd := &cobra.Command{}
		cmd.Flags().String(InputFlagName, testCertsJsonPath, "")
		cmd.Flags().String(VersionFlagName, ver, "")
		cmd.Flags().String(OutputFlagName, testOutputPath, "")

		encryptionCertsPath, version, encryptionCertificatePath, err := ValidateInput(cmd)

		assert.NoError(t, err)
		assert.Equal(t, testCertsJsonPath, encryptionCertsPath)
		assert.Equal(t, ver, version)
		assert.Equal(t, testOutputPath, encryptionCertificatePath)
	}
}

// TestValidateInput_FlagErrors tests error handling for flag retrieval
func TestValidateInput_FlagErrors(t *testing.T) {
	cmd := &cobra.Command{}
	_, _, _, err := ValidateInput(cmd)
	assert.Error(t, err)
}

// TestProcess_Success tests successful certificate extraction
func TestProcess_Success(t *testing.T) {
	result, err := Process(testCertsJsonPath, testVersion2)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "BEGIN CERTIFICATE")
	assert.Contains(t, result, "END CERTIFICATE")
}

// TestProcess_Version1 tests extraction of version 1.0.21
func TestProcess_Version1(t *testing.T) {
	result, err := Process(testCertsJsonPath, testVersion1)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "BEGIN CERTIFICATE")
}

// TestProcess_Version3 tests extraction of version 1.0.23
func TestProcess_Version3(t *testing.T) {
	result, err := Process(testCertsJsonPath, testVersion3)
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
	result, err := Process(testInvalidPath, testVersion2)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "doesn't exist")
}
