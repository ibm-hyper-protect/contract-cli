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
	"github.com/ibm-hyper-protect/contract-go/v2/certificate"
	"github.com/spf13/cobra"
)

var (
	// validateEncryptionCertificateCmd represents the encryption certificate validity command
	validateEncryptionCertificateCmd = &cobra.Command{
		Use:   common.ValidateEncryptionCertParamName,
		Short: common.ValidateEncryptionCertParamShortDescription,
		Long:  common.ValidateEncryptionCertParamLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			encryptionCertsPath, err := validateInputEncryptionCertificates(cmd)
			if err != nil {
				log.Fatal(err)
			}

			encryptionCert, _ := validateEncryptionCertfile(encryptionCertsPath)

			msg, err := certificate.HpcrValidateEncryptionCertificate(encryptionCert)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Encryption certificate status ", msg)
		},
	}
)

// init - cobra init function
func init() {
	rootCmd.AddCommand(validateEncryptionCertificateCmd)

	validateEncryptionCertificateCmd.PersistentFlags().String(common.FileInFlagName, "", common.ValidateEncryptionCertVersionFlagDescription)
}

// validateInputEncryptionCertificates - function to validate validate encryption certificate input
func validateInputEncryptionCertificates(cmd *cobra.Command) (string, error) {
	encryptionCertsPath, err := cmd.Flags().GetString(common.FileInFlagName)
	if err != nil {
		return "", err
	}

	return encryptionCertsPath, nil
}

func validateEncryptionCertfile(encryptionCertsPath string) (string, error) {
	if !common.CheckFileFolderExists(encryptionCertsPath) {
		return "", fmt.Errorf("The specified path does not contain the encryption certificates.")
	}

	encryptionCert, err := common.ReadDataFromFile(encryptionCertsPath)
	if err != nil {
		return "", err
	}

	return encryptionCert, nil
}
