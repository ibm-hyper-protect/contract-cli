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

// getCertificateCmd represents the get-certificate command
var getCertificateCmd = &cobra.Command{
	Use:   common.GetCertParamName,
	Short: common.GetCertParamShortDescription,
	Long:  common.GetCertParamLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		encryptionCertsPath, version, encryptionCertOutputPath, err := validateInputGetCertificate(cmd)
		if err != nil {
			log.Fatal(err)
		}

		encryptionCertificate, err := getEncryptionCertificate(encryptionCertsPath, version)
		if err != nil {
			log.Fatal(err)
		}

		err = printCertificate(encryptionCertificate, encryptionCertOutputPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(getCertificateCmd)

	getCertificateCmd.PersistentFlags().String(common.FileInFlagName, "", common.GetCertFileInFlagDescription)
	getCertificateCmd.PersistentFlags().String(common.VersionFlagName, "", common.GetCertVersionFlagDescription)
	getCertificateCmd.PersistentFlags().String(common.FileOutFlagName, "", common.GetCertFileOutFlagDescription)
}

// validateInputGetCertificate - function to validate get-certificate input
func validateInputGetCertificate(cmd *cobra.Command) (string, string, string, error) {
	encryptionCertsPath, err := cmd.Flags().GetString(common.FileInFlagName)
	if err != nil {
		return "", "", "", err
	}

	version, err := cmd.Flags().GetString(common.VersionFlagName)
	if err != nil {
		return "", "", "", err
	}

	encryptionCertificatePath, err := cmd.Flags().GetString(common.FileOutFlagName)
	if err != nil {
		return "", "", "", err
	}

	return encryptionCertsPath, version, encryptionCertificatePath, nil
}

// getEncryptionCertificate - function to get encryption certificate
func getEncryptionCertificate(encryptionCertsPath, version string) (string, error) {
	if !common.CheckFileFolderExists(encryptionCertsPath) {
		return "", fmt.Errorf("the path to encryption certificates doesn't exist")
	}

	encryptionCertsJson, err := common.ReadDataFromFile(encryptionCertsPath)
	if err != nil {
		return "", err
	}

	_, outputCertificate, err := certificate.HpcrGetEncryptionCertificateFromJson(encryptionCertsJson, version)
	if err != nil {
		return "", err
	}

	return outputCertificate, nil
}

// printCertificate - function to print encryption certificate or redirect it to a file
func printCertificate(cert, certPath string) error {
	if certPath != "" {
		err := common.WriteDataToFile(certPath, cert)
		if err != nil {
			return err
		}
		fmt.Println("Successfully added encryption certificate ")
	} else {
		fmt.Println(cert)
	}

	return nil
}
