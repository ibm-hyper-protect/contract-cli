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
	"github.com/ibm-hyper-protect/contract-cli/lib/validateEncryptionCertificate"
	"github.com/ibm-hyper-protect/contract-go/v2/certificate"
	"github.com/spf13/cobra"
)

var (
	// validateEncryptionCertificateCmd represents the encryption certificate validity command
	validateEncryptionCertificateCmd = &cobra.Command{
		Use:   validateEncryptionCertificate.ParameterName,
		Short: validateEncryptionCertificate.ParameterShortDescription,
		Long:  validateEncryptionCertificate.ParameterLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			encryptionCertsPath, err := validateEncryptionCertificate.ValidateInput(cmd)
			if err != nil {
				log.Fatal(err)
			}

			encryptionCert, _ := validateEncryptionCertificate.GetEncryptionCertfile(encryptionCertsPath)

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

	requiredFlags := map[string]bool{
		"in": true,
	}
	validateEncryptionCertificateCmd.PersistentFlags().String(validateEncryptionCertificate.InputFlagName, "", validateEncryptionCertificate.CertVersionFlagDescription)
	common.SetCustomHelpTemplate(validateEncryptionCertificateCmd, requiredFlags)
}
