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

package validateEncryptionCertificate

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testCertPath    = "../../samples/certificate/active.crt"
	testInvalidPath = "../../build/file/file_not_exists.crt"
)

// TestValidateInput_Success tests ValidateInput with valid input
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testCertPath, "")

	encryptionCertsPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testCertPath, encryptionCertsPath)
}

// TestValidateInput_WithRelativePath tests ValidateInput with relative path
func TestValidateInput_WithRelativePath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, "./certificate.crt", "")
	encryptionCertsPath, err := ValidateInput(cmd)
	assert.NoError(t, err)
	assert.Equal(t, "./certificate.crt", encryptionCertsPath)
}

// TestValidateInput_FlagErrors tests error handling for flag retrieval
func TestValidateInput_FlagErrors(t *testing.T) {
	cmd := &cobra.Command{}
	_, err := ValidateInput(cmd)
	assert.Error(t, err)
}

// TestGetEncryptionCertfile_Success tests successful certificate file reading
func TestGetEncryptionCertfile_Success(t *testing.T) {
	result, err := GetEncryptionCertfile(testCertPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "BEGIN CERTIFICATE")
	assert.Contains(t, result, "END CERTIFICATE")
}

// TestGetEncryptionCertfile_InvalidPath tests with non-existent file path
func TestGetEncryptionCertfile_InvalidPath(t *testing.T) {
	result, err := GetEncryptionCertfile(testInvalidPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "does not contain the encryption certificates")
}

// TestGetEncryptionCertfile_EmptyPath tests with empty path
func TestGetEncryptionCertfile_EmptyPath(t *testing.T) {
	result, err := GetEncryptionCertfile("")
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestGetEncryptionCertfile_WithPemExtension tests reading .pem file
func TestGetEncryptionCertfile_WithPemExtension(t *testing.T) {
	testPemPath := "../../samples/contract-expiry/personal_ca.pem"
	result, err := GetEncryptionCertfile(testPemPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestGetEncryptionCertfile_WithCrtExtension tests reading .crt file
func TestGetEncryptionCertfile_WithCrtExtension(t *testing.T) {
	testCrtPath := "../../samples/contract-expiry/personal_ca.crt"
	result, err := GetEncryptionCertfile(testCrtPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}
