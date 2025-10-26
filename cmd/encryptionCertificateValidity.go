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

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/certificate"
	"github.com/spf13/cobra"
)

var (
	// validateEncryptionCertificagteCmd represents the encryption certificate validity command
	validateEncryptionCertificagteCmd = &cobra.Command{
		Use:   common.ValidateCertParamName,
		Short: common.ValidateCertParamShortDescription,
		Long:  common.ValidateCertParamLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			formatType, err := validateInputEncryptionCertificates(cmd)
			if err != nil {
				log.Fatal(err)
			}

			certificates, err := certificate.HpcrDownloadEncryptionCertificates(encrption_cert_versions, formatType, "")
			if err != nil {
				log.Fatal(err)
			}

			_, err = certificate.HpcrEncryptionCertificatesValidation(certificates)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	encrption_cert_versions []string
)

// init - cobra init function
func init() {
	rootCmd.AddCommand(validateEncryptionCertificagteCmd)

	validateEncryptionCertificagteCmd.PersistentFlags().StringSliceVarP(&encrption_cert_versions, common.VersionFlagName, "", []string{}, common.ValidateCertVersionFlagDescription)
	validateEncryptionCertificagteCmd.PersistentFlags().String(common.DataFormatFlagName, common.DataFormatDefault, common.DataFormatFlagDescription)
}

// validateInputEncryptionCertificates - function to validate encryption certificate validity inputs
func validateInputEncryptionCertificates(cmd *cobra.Command) (string, error) {
	formatType, err := cmd.Flags().GetString(common.DataFormatFlagName)
	if err != nil {
		return "", err
	}

	return formatType, nil
}
