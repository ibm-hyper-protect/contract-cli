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
	"github.com/ibm-hyper-protect/contract-go/certificate"
	"github.com/spf13/cobra"
)

const (
	successMessageDownloadCertificate = "Successfully stored certificates"
)

var (
	// downloadCertificatesCmd represents the download-certificate command
	downloadCertificatesCmd = &cobra.Command{
		Use:   common.DownloadCertParamName,
		Short: common.DownloadCertParamShortDescription,
		Long:  common.DownloadCertParamLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			formatType, certificatePath, err := validateInputDownloadCertificates(cmd)
			if err != nil {
				log.Fatal(err)
			}

			certificates, err := certificate.HpcrDownloadEncryptionCertificates(versions, formatType, "")
			if err != nil {
				log.Fatal(err)
			}

			err = printDownloadCertificates(certificates, certificatePath)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	versions []string
)

// init - cobra init function
func init() {
	rootCmd.AddCommand(downloadCertificatesCmd)

	downloadCertificatesCmd.PersistentFlags().StringSliceVarP(&versions, common.VersionFlagName, "", []string{}, common.EncryptionCertVersionFlagDescription)
	downloadCertificatesCmd.PersistentFlags().String(common.DataFormatFlagName, common.DataFormatDefault, common.DataFormatFlagDescription)
	downloadCertificatesCmd.PersistentFlags().String(common.FileOutFlagName, "", common.EncryptionCertsFlagDescription)
}

// validateInputDownloadCertificates - function to validate download-certificate inputs
func validateInputDownloadCertificates(cmd *cobra.Command) (string, string, error) {
	formatType, err := cmd.Flags().GetString(common.DataFormatFlagName)
	if err != nil {
		return "", "", err
	}

	certificatePath, err := cmd.Flags().GetString(common.FileOutFlagName)
	if err != nil {
		return "", "", err
	}

	return formatType, certificatePath, nil
}

// printDownloadCertificates - function to print encryption certificates or redirect output to a file
func printDownloadCertificates(certificates, certificatePath string) error {
	if certificatePath != "" {
		err := common.WriteDataToFile(certificatePath, certificates)
		if err != nil {
			return err
		}
		fmt.Println(successMessageDownloadCertificate)
	} else {
		fmt.Println(certificates)
	}

	return nil
}
