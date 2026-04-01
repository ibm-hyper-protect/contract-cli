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
	"os"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-cli/lib/checkCertificateRevocation"
	"github.com/ibm-hyper-protect/contract-go/v2/certificate"
	"github.com/spf13/cobra"
)

// checkCertificateRevocationCmd represents the check-crl command
var checkCertificateRevocationCmd = &cobra.Command{
	Use:   checkCertificateRevocation.ParameterName,
	Short: checkCertificateRevocation.ParameterShortDescription,
	Long:  checkCertificateRevocation.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		certPath, crlPath, err := checkCertificateRevocation.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		cert, err := checkCertificateRevocation.GetCertificateData(certPath)
		if err != nil {
			log.Fatal(err)
		}

		crl, err := checkCertificateRevocation.GetCRLData(crlPath)
		if err != nil {
			log.Fatal(err)
		}

		revoked, msg, err := certificate.HpcrCheckCertificateRevocation(cert, crl)
		if err != nil {
			log.Fatal(err)
		}

		if revoked {
			fmt.Printf("WARNING: Certificate has been REVOKED: %s\n", msg)
			os.Exit(1)
		}

		fmt.Printf("Certificate is valid (not revoked): %s\n", msg)
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(checkCertificateRevocationCmd)

	requiredFlags := map[string]bool{
		"cert": true,
		"crl":  true,
	}

	checkCertificateRevocationCmd.PersistentFlags().String(checkCertificateRevocation.InputFlagName, "", checkCertificateRevocation.CertFlagDescription)
	checkCertificateRevocationCmd.PersistentFlags().String(checkCertificateRevocation.CRLFlagName, "", checkCertificateRevocation.CRLFlagDescription)
	common.SetCustomHelpTemplate(checkCertificateRevocationCmd, requiredFlags)
	common.SetCustomErrorTemplate(checkCertificateRevocationCmd)
}
