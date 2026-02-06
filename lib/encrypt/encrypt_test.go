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

// TestValidateInput_WithoutOutputPath tests ValidateInput without output path (optional)
func TestValidateInput_WithoutOutputPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(OsVersionFlagName, testOsVersion, "")
	cmd.Flags().String(CertFlagName, testCertPath, "")
	cmd.Flags().String(PrivateKeyFlagName, testPrivateKeyPath, "")
	cmd.Flags().String(OutputFlagName, "", "")

	inputData, osVersion, certPath, privateKeyPath, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, inputData)
	assert.Equal(t, testOsVersion, osVersion)
	assert.Equal(t, testCertPath, certPath)
	assert.Equal(t, testPrivateKeyPath, privateKeyPath)
	assert.Equal(t, "", outputPath)
}

// TestValidateInput_WithEmptyCertPath tests ValidateInput with empty cert path (optional)
func TestValidateInput_WithEmptyCertPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(OsVersionFlagName, testOsVersion, "")
	cmd.Flags().String(CertFlagName, "", "")
	cmd.Flags().String(PrivateKeyFlagName, testPrivateKeyPath, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")

	inputData, osVersion, certPath, privateKeyPath, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, inputData)
	assert.Equal(t, testOsVersion, osVersion)
	assert.Equal(t, "", certPath)
	assert.Equal(t, testPrivateKeyPath, privateKeyPath)
	assert.Equal(t, testOutputPath, outputPath)
}

// TestValidateInput_WithEmptyPrivateKeyPath tests ValidateInput with empty private key path (optional)
func TestValidateInput_WithEmptyPrivateKeyPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(OsVersionFlagName, testOsVersion, "")
	cmd.Flags().String(CertFlagName, testCertPath, "")
	cmd.Flags().String(PrivateKeyFlagName, "", "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")

	inputData, osVersion, certPath, privateKeyPath, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, inputData)
	assert.Equal(t, testOsVersion, osVersion)
	assert.Equal(t, testCertPath, certPath)
	assert.Equal(t, "", privateKeyPath)
	assert.Equal(t, testOutputPath, outputPath)
}

// TestValidateInput_FlagErrors tests error handling for flag retrieval
func TestValidateInput_FlagErrors(t *testing.T) {
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

// TestValidateInputEncryptContractExpiry_FlagErrors tests error handling
func TestValidateInputEncryptContractExpiry_FlagErrors(t *testing.T) {
	cmd := &cobra.Command{}
	_, _, _, _, _, _, err := ValidateInputEncryptContractExpiry(cmd)
	assert.Error(t, err)
}

// TestGenerateSignedEncryptContract_InvalidContractPath tests with invalid contract path
func TestGenerateSignedEncryptContract_InvalidContractPath(t *testing.T) {
	result, err := GenerateSignedEncryptContract(testInvalidPath, testOsVersion, testCertPath, testPrivateKeyPath)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "doesn't exist")
}
