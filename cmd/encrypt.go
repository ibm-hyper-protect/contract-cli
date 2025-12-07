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
	"log"

	"github.com/ibm-hyper-protect/contract-cli/lib/encrypt"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   encrypt.ParameterName,
	Short: encrypt.ParameterShortDescription,
	Long:  encrypt.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		inputDataPath, osVersion, certPath, privateKeyPath, outputPath, err := encrypt.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		contractExpiryFlag, caCert, caKey, csrParam, csr, expiryDays, err := encrypt.ValidateInputEncryptContractExpiry(cmd)
		if err != nil {
			log.Fatal(err)
		}

		var signedEncryptContract string
		if !contractExpiryFlag {
			signedEncryptContract, err = encrypt.GenerateSignedEncryptContract(inputDataPath, osVersion, certPath, privateKeyPath)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			signedEncryptContract, err = encrypt.GenerateSignedEncryptContractExpiry(inputDataPath, osVersion, certPath, privateKeyPath, caCert, caKey, csrParam, csr, expiryDays)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = encrypt.Output(signedEncryptContract, outputPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.PersistentFlags().String(encrypt.InputFlagName, "", encrypt.InputFlagDescription)
	encryptCmd.PersistentFlags().String(encrypt.OsVersionFlagName, "", encrypt.OsVersionFlagDescription)
	encryptCmd.PersistentFlags().String(encrypt.CertFlagName, "", encrypt.CertFlagDescription)
	encryptCmd.PersistentFlags().String(encrypt.PrivateKeyFlagName, "", encrypt.PrivateKeyFlagDescription)
	encryptCmd.PersistentFlags().String(encrypt.OutputFlagName, "", encrypt.OutputFlagDescription)
	encryptCmd.PersistentFlags().Bool(encrypt.ContractExpiryFlag, encrypt.DefaultContractExpiryFlag, encrypt.ContractExpiryFlagDescription)
	encryptCmd.PersistentFlags().String(encrypt.CaCertFlag, "", encrypt.CaCertFlagDescription)
	encryptCmd.PersistentFlags().String(encrypt.CaKeyFlag, "", encrypt.CaKeyFlagDescription)
	encryptCmd.PersistentFlags().String(encrypt.CsrDataFlag, "", encrypt.CsrDataFlagDescription)
	encryptCmd.PersistentFlags().String(encrypt.CsrFlag, "", encrypt.CsrFlagDescription)
	encryptCmd.PersistentFlags().Int(encrypt.ExpiryDaysFlag, 0, encrypt.ExpiryDaysFlagDescription)
}
