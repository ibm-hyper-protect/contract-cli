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

package cmd

import (
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/contract"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   common.EncryptParamName,
	Short: common.EncryptParamShortDescription,
	Long:  common.EncryptParamLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		inputDataPath, osVersion, certPath, privateKeyPath, outputPath, err := validateInputEncrypt(cmd)
		if err != nil {
			log.Fatal(err)
		}

		contractExpiryFlag, caCert, caKey, csrParam, csr, expiryDays, err := validateInputEncryptContractExpiry(cmd)
		if err != nil {
			log.Fatal(err)
		}

		var signedEncryptContract string
		if !contractExpiryFlag {
			signedEncryptContract, err = generateSignedEncryptContract(inputDataPath, osVersion, certPath, privateKeyPath)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			signedEncryptContract, err = generateSignedEncryptContractExpiry(inputDataPath, osVersion, certPath, privateKeyPath, caCert, caKey, csrParam, csr, expiryDays)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = printSignedEncryptContract(signedEncryptContract, outputPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.PersistentFlags().String(common.FileInFlagName, "", common.EncryptInputFlagDescription)
	encryptCmd.PersistentFlags().String(common.OsVersionFlagName, "", common.OsVersionFlagDescription)
	encryptCmd.PersistentFlags().String(common.CertFlagName, "", common.CertFlagDescription)
	encryptCmd.PersistentFlags().String(common.PrivateKeyFlagName, "", common.PrivateKeyFlagDescription)
	encryptCmd.PersistentFlags().String(common.FileOutFlagName, "", common.EncryptOutputFlagDescription)

	encryptCmd.PersistentFlags().Bool(common.EncryptContractExpiryFlagName, common.EncryptContractExpiryFlagDefault, common.EncryptContractExpiryFlagDescription)
	encryptCmd.PersistentFlags().String(common.EncryptCaCertFlagName, "", common.EncryptCaCertFlagDescription)
	encryptCmd.PersistentFlags().String(common.EncryptCaKeyFlagName, "", common.EncryptCaKeyFlagDescription)
	encryptCmd.PersistentFlags().String(common.EncryptCsrDataFlagName, "", common.EncryptCsrDataFlagDescription)
	encryptCmd.PersistentFlags().String(common.EncryptCsrFlagName, "", common.EncryptCsrFlagDescription)
	encryptCmd.PersistentFlags().Int(common.EncryptExpiryDaysFlagName, 0, common.EncryptExpiryDaysFlagDescription)
}

// validateInputEncrypt - function to validate inputs of encrypt
func validateInputEncrypt(cmd *cobra.Command) (string, string, string, string, string, error) {
	inputData, err := cmd.Flags().GetString(common.FileInFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	osVersion, err := cmd.Flags().GetString(common.OsVersionFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	certPath, err := cmd.Flags().GetString(common.CertFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	privateKeyPath, err := cmd.Flags().GetString(common.PrivateKeyFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	outputPath, err := cmd.Flags().GetString(common.FileOutFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	return inputData, osVersion, certPath, privateKeyPath, outputPath, nil
}

// validateInputEncryptContractExpiry - function to validate input of contract expiry input
func validateInputEncryptContractExpiry(cmd *cobra.Command) (bool, string, string, string, string, int, error) {
	contractExpiryFlag, err := cmd.Flags().GetBool(common.EncryptContractExpiryFlagName)
	if err != nil {
		return false, "", "", "", "", 0, err
	}

	caCert, err := cmd.Flags().GetString(common.EncryptCaCertFlagName)
	if err != nil {
		return false, "", "", "", "", 0, err
	}

	caKey, err := cmd.Flags().GetString(common.EncryptCaKeyFlagName)
	if err != nil {
		return false, "", "", "", "", 0, err
	}

	csrParam, err := cmd.Flags().GetString(common.EncryptCsrDataFlagName)
	if err != nil {
		return false, "", "", "", "", 0, err
	}

	csr, err := cmd.Flags().GetString(common.EncryptCsrFlagName)
	if err != nil {
		return false, "", "", "", "", 0, err
	}

	expiryDays, err := cmd.Flags().GetInt(common.EncryptExpiryDaysFlagName)
	if err != nil {
		return false, "", "", "", "", 0, err
	}

	return contractExpiryFlag, caCert, caKey, csrParam, csr, expiryDays, nil
}

// generateSignedEncryptContract - function to generate signed and encrypted contract
func generateSignedEncryptContract(inputDataPath, osVersion, certPath, privateKeyPath string) (string, error) {
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

// generateSignedEncryptContractExpiry - function to generated signed and encrypted contract with contract expiry
func generateSignedEncryptContractExpiry(inputDataPath, osVersion, certPath, privateKeyPath, caCertPath, caKeyPath, csrParamPath, csrPath string, expiryDays int) (string, error) {

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
	if !common.CheckFileFolderExists(inputDataPath) {
		return "", "", "", fmt.Errorf("the contract path doesn't exist")
	}

	inputData, err := common.ReadDataFromFile(inputDataPath)
	if err != nil {
		return "", "", "", err
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

// printSignedEncryptContract - function to print signed and encrypted contract or redirect it to a file
func printSignedEncryptContract(signedEncryptContract, outputPath string) error {
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
