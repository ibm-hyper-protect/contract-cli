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
	"fmt"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "encrypt"
	ParameterShortDescription = "Generate signed and encrypted contract"
	ParameterLongDescription  = `Generate a signed and encrypted contract for Hyper Protect deployment.

Creates a cryptographically signed contract using your private key and encrypts it
with the platform's encryption certificate. Supports optional contract expiry for
enhanced security.`
	InputFlagDescription          = "Path to unencrypted contract YAML file (use '-' for standard input)"
	OutputFlagDescription         = "Path to save signed and encrypted contract"
	ContractExpiryFlag            = "contract-expiry"
	DefaultContractExpiryFlag     = false
	ContractExpiryFlagDescription = "Enable contract expiry feature (requires CA cert and key)"
	CaCertFlag                    = "cacert"
	CaCertFlagDescription         = "Path to CA certificate file (required with --contract-expiry)"
	CaKeyFlag                     = "cakey"
	CaKeyFlagDescription          = "Path to CA private key file (required with --contract-expiry)"
	CsrDataFlag                   = "csrParam"
	CsrDataFlagDescription        = "Path to CSR parameters JSON file"
	CsrFlag                       = "csr"
	CsrFlagDescription            = "Path to Certificate Signing Request (CSR) file"
	ExpiryDaysFlag                = "expiry"
	ExpiryDaysFlagDescription     = "Contract validity period in days (required with --contract-expiry)"
	InputFlagName                 = "in"
	OutputFlagName                = "out"
	OsVersionFlagName             = "os"
	OsVersionFlagDescription      = "Target Hyper Protect platform (hpvs, hpcr-rhvs, or hpcc-peerpod)"
	CertFlagName                  = "cert"
	CertFlagDescription           = "Path to encryption certificate file"
	PrivateKeyFlagName            = "priv"
	PrivateKeyFlagDescription     = "Path to private key file for signing"
)

// ValidateInput - function to validate inputs of encrypt
func ValidateInput(cmd *cobra.Command) (string, string, string, string, string, error) {
	inputData, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	if inputData == "" {
		err := fmt.Errorf("Error: required flag '--in' is missing")
		common.SetMandatoryFlagError(cmd, err)
	}

	// Validate stdin input conflicts
	common.ValidateStdinInput(cmd, inputData)

	osVersion, err := cmd.Flags().GetString(OsVersionFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	certPath, err := cmd.Flags().GetString(CertFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	privateKeyPath, err := cmd.Flags().GetString(PrivateKeyFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	outputPath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	return inputData, osVersion, certPath, privateKeyPath, outputPath, nil
}

// ValidateInputEncryptContractExpiry - function to validate input of contract expiry input
func ValidateInputEncryptContractExpiry(cmd *cobra.Command) (bool, string, string, string, string, int, error) {
	contractExpiryFlag, err := cmd.Flags().GetBool(ContractExpiryFlag)
	if err != nil {
		return false, "", "", "", "", 0, err
	}

	caCert, err := cmd.Flags().GetString(CaCertFlag)
	if err != nil {
		return false, "", "", "", "", 0, err
	}

	caKey, err := cmd.Flags().GetString(CaKeyFlag)
	if err != nil {
		return false, "", "", "", "", 0, err
	}

	csrParam, err := cmd.Flags().GetString(CsrDataFlag)
	if err != nil {
		return false, "", "", "", "", 0, err
	}

	csr, err := cmd.Flags().GetString(CsrFlag)
	if err != nil {
		return false, "", "", "", "", 0, err
	}

	expiryDays, err := cmd.Flags().GetInt(ExpiryDaysFlag)
	if err != nil {
		return false, "", "", "", "", 0, err
	}

	return contractExpiryFlag, caCert, caKey, csrParam, csr, expiryDays, nil
}

// GenerateSignedEncryptContract - function to generate signed and encrypted contract
func GenerateSignedEncryptContract(inputDataPath, osVersion, certPath, privateKeyPath string) (string, error) {
	inputData, cert, privateKey, err := commonParameters(inputDataPath, certPath, privateKeyPath)
	if err != nil {
		return "", err
	}

	signedEncryptContract, _, _, err := contract.HpcrContractSignedEncrypted(inputData, osVersion, cert, privateKey)
	if err != nil {
		return "", err
	}

	return signedEncryptContract, nil
}

// GenerateSignedEncryptContractExpiry - function to generated signed and encrypted contract with contract expiry
func GenerateSignedEncryptContractExpiry(inputDataPath, osVersion, certPath, privateKeyPath, caCertPath, caKeyPath, csrParamPath, csrPath string, expiryDays int) (string, error) {

	inputData, cert, privateKey, err := commonParameters(inputDataPath, certPath, privateKeyPath)
	if err != nil {
		return "", err
	}

	caCert, err := common.GetDataFromFile(caCertPath)
	if err != nil {
		return "", err
	}

	caKey, err := common.GetDataFromFile(caKeyPath)
	if err != nil {
		return "", err
	}

	csrParam, err := common.GetDataFromFile(csrParamPath)
	if err != nil {
		return "", err
	}

	csr, err := common.GetDataFromFile(csrPath)
	if err != nil {
		return "", err
	}

	signedEncryptContract, _, _, err := contract.HpcrContractSignedEncryptedContractExpiry(inputData, osVersion, cert, privateKey, caCert, caKey, csrParam, csr, expiryDays)
	if err != nil {
		return "", err
	}

	return signedEncryptContract, nil
}

// commonParameters - function to fetch common details
func commonParameters(inputDataPath, certPath, privateKeyPath string) (string, string, string, error) {
	var inputData string
	var err error

	if inputDataPath == "-" {
		inputData, err = common.ReadDataFromStdin()
		if err != nil {
			return "", "", "", fmt.Errorf("unable to read input from standard input: %w", err)
		}
	} else {
		if !common.CheckFileFolderExists(inputDataPath) {
			return "", "", "", fmt.Errorf("the contract path doesn't exist")
		}

		inputData, err = common.ReadDataFromFile(inputDataPath)
		if err != nil {
			return "", "", "", err
		}
	}

	cert, err := common.GetDataFromFile(certPath)
	if err != nil {
		return "", "", "", err
	}

	privateKey, err := common.GetPrivateKey(privateKeyPath)
	if err != nil {
		return "", "", "", err
	}

	return inputData, cert, privateKey, nil
}

// Output - function to print signed and encrypted contract or redirect it to a file
func Output(signedEncryptContract, outputPath string) error {
	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, signedEncryptContract)
		if err != nil {
			return err
		}
		fmt.Println("Successfully generated signed and encrypted contract")
	} else {
		fmt.Println(signedEncryptContract)
	}

	return nil
}
