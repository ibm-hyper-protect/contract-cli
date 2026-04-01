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
	"github.com/ibm-hyper-protect/contract-cli/lib/validateCertChain"
	"github.com/ibm-hyper-protect/contract-go/v2/certificate"
	"github.com/spf13/cobra"
)

// validateCertChainCmd represents the validate-cert-chain command
var validateCertChainCmd = &cobra.Command{
	Use:   validateCertChain.ParameterName,
	Short: validateCertChain.ParameterShortDescription,
	Long:  validateCertChain.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		certPath, intermediatePath, rootPath, err := validateCertChain.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		cert, err := validateCertChain.GetCertificateData(certPath)
		if err != nil {
			log.Fatal(err)
		}

		intermediate, err := validateCertChain.GetIntermediateCertData(intermediatePath)
		if err != nil {
			log.Fatal(err)
		}

		root, err := validateCertChain.GetRootCertData(rootPath)
		if err != nil {
			log.Fatal(err)
		}

		valid, msg, err := certificate.HpcrValidateCertChain(cert, intermediate, root)
		if err != nil {
			log.Fatal(err)
		}

		if !valid {
			fmt.Printf("Certificate chain validation failed: %s\n", msg)
			os.Exit(1)
		}

		fmt.Printf("Certificate chain validated successfully: %s\n", msg)
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(validateCertChainCmd)

	requiredFlags := map[string]bool{
		"cert":         true,
		"intermediate": true,
		"root":         true,
	}

	validateCertChainCmd.PersistentFlags().String(validateCertChain.InputFlagName, "", validateCertChain.CertFlagDescription)
	validateCertChainCmd.PersistentFlags().String(validateCertChain.IntermediateFlagName, "", validateCertChain.IntermediateFlagDescription)
	validateCertChainCmd.PersistentFlags().String(validateCertChain.RootFlagName, "", validateCertChain.RootFlagDescription)
	common.SetCustomHelpTemplate(validateCertChainCmd, requiredFlags)
	common.SetCustomErrorTemplate(validateCertChainCmd)
}
