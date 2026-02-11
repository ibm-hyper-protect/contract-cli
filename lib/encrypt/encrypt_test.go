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

package encrypt

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testContractPath   = "../../samples/contract.yaml"
	testCertPath       = "../../samples/certificate/active.crt"
	testPrivateKeyPath = "../../samples/sign/private.pem"
	testCaCertPath     = "../../samples/contract-expiry/personal_ca.crt"
	testCaKeyPath      = "../../samples/contract-expiry/personal_ca.pem"
	testCsrPath        = "../../samples/contract-expiry/csr.pem"
	testCsrParamPath   = "../../samples/contract-expiry/csr.pem"
	testOutputPath     = "../../build/test_encrypt_output.txt"
	testInvalidPath    = "../../build/file/file_not_exists.txt"
	testOsVersion      = "hpvs"
)

// TestValidateInput_Success tests ValidateInput with all required flags
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(OsVersionFlagName, testOsVersion, "")
	cmd.Flags().String(CertFlagName, testCertPath, "")
	cmd.Flags().String(PrivateKeyFlagName, testPrivateKeyPath, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")

	inputData, osVersion, certPath, privateKeyPath, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, inputData)
	assert.Equal(t, testOsVersion, osVersion)
	assert.Equal(t, testCertPath, certPath)
	assert.Equal(t, testPrivateKeyPath, privateKeyPath)
	assert.Equal(t, testOutputPath, outputPath)
}

// TestValidateInput_WithoutFlags tests ValidateInput when flags are not set
func TestValidateInput_WithoutFlags(t *testing.T) {
	cmd := &cobra.Command{}
	_, _, _, _, _, err := ValidateInput(cmd)
	assert.Error(t, err)
}

// TestValidateInputEncryptContractExpiry_Success tests contract expiry validation
func TestValidateInputEncryptContractExpiry_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool(ContractExpiryFlag, true, "")
	cmd.Flags().String(CaCertFlag, testCaCertPath, "")
	cmd.Flags().String(CaKeyFlag, testCaKeyPath, "")
	cmd.Flags().String(CsrDataFlag, testCsrParamPath, "")
	cmd.Flags().String(CsrFlag, testCsrPath, "")
	cmd.Flags().Int(ExpiryDaysFlag, 365, "")

	contractExpiryFlag, caCert, caKey, csrParam, csr, expiryDays, err := ValidateInputEncryptContractExpiry(cmd)

	assert.NoError(t, err)
	assert.True(t, contractExpiryFlag)
	assert.Equal(t, testCaCertPath, caCert)
	assert.Equal(t, testCaKeyPath, caKey)
	assert.Equal(t, testCsrParamPath, csrParam)
	assert.Equal(t, testCsrPath, csr)
	assert.Equal(t, 365, expiryDays)
}

// TestValidateInputEncryptContractExpiry_Disabled tests with contract expiry disabled
func TestValidateInputEncryptContractExpiry_Disabled(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool(ContractExpiryFlag, false, "")
	cmd.Flags().String(CaCertFlag, "", "")
	cmd.Flags().String(CaKeyFlag, "", "")
	cmd.Flags().String(CsrDataFlag, "", "")
	cmd.Flags().String(CsrFlag, "", "")
	cmd.Flags().Int(ExpiryDaysFlag, 0, "")

	contractExpiryFlag, caCert, caKey, csrParam, csr, expiryDays, err := ValidateInputEncryptContractExpiry(cmd)

	assert.NoError(t, err)
	assert.False(t, contractExpiryFlag)
	assert.Equal(t, "", caCert)
	assert.Equal(t, "", caKey)
	assert.Equal(t, "", csrParam)
	assert.Equal(t, "", csr)
	assert.Equal(t, 0, expiryDays)
}

// TestValidateInputEncryptContractExpiry_WithoutFlags tests error handling
func TestValidateInputEncryptContractExpiry_WithoutFlags(t *testing.T) {
	cmd := &cobra.Command{}
	_, _, _, _, _, _, err := ValidateInputEncryptContractExpiry(cmd)
	assert.Error(t, err)
}

// TestGenerateSignedEncryptContract_Success tests successful contract generation
func TestGenerateSignedEncryptContract_Success(t *testing.T) {
	result, err := GenerateSignedEncryptContract(testContractPath, testOsVersion, testCertPath, testPrivateKeyPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "hyper-protect-basic")
}

// TestGenerateSignedEncryptContract_InvalidContractPath tests with invalid contract path
func TestGenerateSignedEncryptContract_InvalidContractPath(t *testing.T) {
	result, err := GenerateSignedEncryptContract(testInvalidPath, testOsVersion, testCertPath, testPrivateKeyPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "doesn't exist")
}

// TestGenerateSignedEncryptContract_InvalidCertPath tests with invalid cert path
func TestGenerateSignedEncryptContract_InvalidCertPath(t *testing.T) {
	result, err := GenerateSignedEncryptContract(testContractPath, testOsVersion, testInvalidPath, testPrivateKeyPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestGenerateSignedEncryptContract_InvalidPrivateKeyPath tests with invalid private key path
func TestGenerateSignedEncryptContract_InvalidPrivateKeyPath(t *testing.T) {
	result, err := GenerateSignedEncryptContract(testContractPath, testOsVersion, testCertPath, testInvalidPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestGenerateSignedEncryptContract_InvalidOsVersion tests with invalid OS version
func TestGenerateSignedEncryptContract_InvalidOsVersion(t *testing.T) {
	result, err := GenerateSignedEncryptContract(testContractPath, "invalid_os", testCertPath, testPrivateKeyPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestGenerateSignedEncryptContract_EmptyOsVersion tests with empty OS version
func TestGenerateSignedEncryptContract_EmptyOsVersion(t *testing.T) {
	result, err := GenerateSignedEncryptContract(testContractPath, "", testCertPath, testPrivateKeyPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestGenerateSignedEncryptContractExpiry_InvalidContractPath tests with invalid contract path
func TestGenerateSignedEncryptContractExpiry_InvalidContractPath(t *testing.T) {
	result, err := GenerateSignedEncryptContractExpiry(
		testInvalidPath,
		testOsVersion,
		testCertPath,
		testPrivateKeyPath,
		testCaCertPath,
		testCaKeyPath,
		testCsrParamPath,
		testCsrPath,
		365,
	)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "doesn't exist")
}

// TestGenerateSignedEncryptContractExpiry_InvalidCaCertPath tests with invalid CA cert path
func TestGenerateSignedEncryptContractExpiry_InvalidCaCertPath(t *testing.T) {
	result, err := GenerateSignedEncryptContractExpiry(
		testContractPath,
		testOsVersion,
		testCertPath,
		testPrivateKeyPath,
		testInvalidPath,
		testCaKeyPath,
		testCsrParamPath,
		testCsrPath,
		365,
	)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestGenerateSignedEncryptContractExpiry_InvalidCaKeyPath tests with invalid CA key path
func TestGenerateSignedEncryptContractExpiry_InvalidCaKeyPath(t *testing.T) {
	result, err := GenerateSignedEncryptContractExpiry(
		testContractPath,
		testOsVersion,
		testCertPath,
		testPrivateKeyPath,
		testCaCertPath,
		testInvalidPath,
		testCsrParamPath,
		testCsrPath,
		365,
	)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestGenerateSignedEncryptContractExpiry_InvalidCsrParamPath tests with invalid CSR param path
func TestGenerateSignedEncryptContractExpiry_InvalidCsrParamPath(t *testing.T) {
	result, err := GenerateSignedEncryptContractExpiry(
		testContractPath,
		testOsVersion,
		testCertPath,
		testPrivateKeyPath,
		testCaCertPath,
		testCaKeyPath,
		testInvalidPath,
		testCsrPath,
		365,
	)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestGenerateSignedEncryptContractExpiry_InvalidCsrPath tests with invalid CSR path
func TestGenerateSignedEncryptContractExpiry_InvalidCsrPath(t *testing.T) {
	result, err := GenerateSignedEncryptContractExpiry(
		testContractPath,
		testOsVersion,
		testCertPath,
		testPrivateKeyPath,
		testCaCertPath,
		testCaKeyPath,
		testCsrParamPath,
		testInvalidPath,
		365,
	)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestOutput_WithFilePath tests Output function with valid file path
func TestOutput_WithFilePath(t *testing.T) {
	testData := "hyper-protect-basic.test-data"
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
	testData := "hyper-protect-basic.test-data"
	err := Output(testData, "")
	assert.NoError(t, err)
}

// TestOutput_InvalidPath tests Output function with invalid file path
func TestOutput_InvalidPath(t *testing.T) {
	testData := "hyper-protect-basic.test-data"
	err := Output(testData, testInvalidPath)
	assert.Error(t, err)
}

// TestGenerateSignedEncryptContract_CorruptedContract tests with corrupted contract file
func TestGenerateSignedEncryptContract_CorruptedContract(t *testing.T) {
	corruptedFile := "../../build/corrupted_contract.yaml"
	err := os.WriteFile(corruptedFile, []byte("invalid: yaml: content: ["), 0644)
	assert.NoError(t, err)
	defer os.Remove(corruptedFile)

	result, err := GenerateSignedEncryptContract(corruptedFile, testOsVersion, testCertPath, testPrivateKeyPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestGenerateSignedEncryptContract_EmptyContract tests with empty contract file
func TestGenerateSignedEncryptContract_EmptyContract(t *testing.T) {
	emptyFile := "../../build/empty_contract.yaml"
	err := os.WriteFile(emptyFile, []byte(""), 0644)
	assert.NoError(t, err)
	defer os.Remove(emptyFile)

	result, err := GenerateSignedEncryptContract(emptyFile, testOsVersion, testCertPath, testPrivateKeyPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestGenerateSignedEncryptContract_CorruptedCert tests with corrupted certificate
func TestGenerateSignedEncryptContract_CorruptedCert(t *testing.T) {
	corruptedCert := "../../build/corrupted_cert.crt"
	err := os.WriteFile(corruptedCert, []byte("not a valid certificate"), 0644)
	assert.NoError(t, err)
	defer os.Remove(corruptedCert)

	result, err := GenerateSignedEncryptContract(testContractPath, testOsVersion, corruptedCert, testPrivateKeyPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestGenerateSignedEncryptContract_CorruptedPrivateKey tests with corrupted private key
func TestGenerateSignedEncryptContract_CorruptedPrivateKey(t *testing.T) {
	corruptedKey := "../../build/corrupted_key.pem"
	err := os.WriteFile(corruptedKey, []byte("not a valid private key"), 0644)
	assert.NoError(t, err)
	defer os.Remove(corruptedKey)

	result, err := GenerateSignedEncryptContract(testContractPath, testOsVersion, testCertPath, corruptedKey)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestGenerateSignedEncryptContractExpiry_ZeroExpiryDays tests with zero expiry days
func TestGenerateSignedEncryptContractExpiry_ZeroExpiryDays(t *testing.T) {
	result, err := GenerateSignedEncryptContractExpiry(
		testContractPath,
		testOsVersion,
		testCertPath,
		testPrivateKeyPath,
		testCaCertPath,
		testCaKeyPath,
		testCsrParamPath,
		testCsrPath,
		0,
	)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestGenerateSignedEncryptContractExpiry_NegativeExpiryDays tests with negative expiry days
func TestGenerateSignedEncryptContractExpiry_NegativeExpiryDays(t *testing.T) {
	result, err := GenerateSignedEncryptContractExpiry(
		testContractPath,
		testOsVersion,
		testCertPath,
		testPrivateKeyPath,
		testCaCertPath,
		testCaKeyPath,
		testCsrParamPath,
		testCsrPath,
		-1,
	)
	assert.Error(t, err)
	assert.Equal(t, "", result)
}
