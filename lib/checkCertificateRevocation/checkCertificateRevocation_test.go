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

package checkCertificateRevocation

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testCertPath    = "../../samples/certificate/active.crt"
	testCRLPath     = "../../samples/certificate/active.crt"
	testInvalidPath = "../../build/file/file_not_exists.crt"
)

// TestValidateInput_Success tests ValidateInput with valid input
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testCertPath, "")
	cmd.Flags().String(CRLFlagName, testCRLPath, "")

	certPath, crlPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testCertPath, certPath)
	assert.Equal(t, testCRLPath, crlPath)
}

// TestValidateInput_WithRelativePath tests ValidateInput with relative path
func TestValidateInput_WithRelativePath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, "./certificate.crt", "")
	cmd.Flags().String(CRLFlagName, "./crl.pem", "")
	certPath, crlPath, err := ValidateInput(cmd)
	assert.NoError(t, err)
	assert.Equal(t, "./certificate.crt", certPath)
	assert.Equal(t, "./crl.pem", crlPath)
}

// TestValidateInput_FlagErrors tests error handling for flag retrieval
func TestValidateInput_FlagErrors(t *testing.T) {
	cmd := &cobra.Command{}
	_, _, err := ValidateInput(cmd)
	assert.Error(t, err)
}

// TestGetCertificateData_Success tests successful certificate file reading
func TestGetCertificateData_Success(t *testing.T) {
	result, err := GetCertificateData(testCertPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "BEGIN CERTIFICATE")
	assert.Contains(t, result, "END CERTIFICATE")
}

// TestGetCertificateData_InvalidPath tests with non-existent file path
func TestGetCertificateData_InvalidPath(t *testing.T) {
	result, err := GetCertificateData(testInvalidPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "doesn't exist")
}

// TestGetCRLData_Success tests successful CRL file reading
func TestGetCRLData_Success(t *testing.T) {
	// Using certificate file as placeholder for CRL
	result, err := GetCRLData(testCRLPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestGetCRLData_InvalidPath tests with non-existent file path
func TestGetCRLData_InvalidPath(t *testing.T) {
	result, err := GetCRLData(testInvalidPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "doesn't exist")
}

// TestGetCRLData_WithPemExtension tests reading .pem file
func TestGetCRLData_WithPemExtension(t *testing.T) {
	testPemPath := "../../samples/contract-expiry/personal_ca.pem"
	result, err := GetCRLData(testPemPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}
