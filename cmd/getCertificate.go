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
	"github.com/ibm-hyper-protect/contract-cli/lib/getCertificate"
	"github.com/spf13/cobra"
)

// getCertificateCmd represents the get-certificate command
var getCertificateCmd = &cobra.Command{
	Use:   getCertificate.ParameterName,
	Short: getCertificate.ParameterShortDescription,
	Long:  getCertificate.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		encryptionCertsPath, version, encryptionCertOutputPath, err := getCertificate.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		encryptionCertificate, err := getCertificate.Process(encryptionCertsPath, version)
		if err != nil {
			log.Fatal(err)
		}

		err = getCertificate.Output(encryptionCertificate, encryptionCertOutputPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(getCertificateCmd)

	requiredFlags := map[string]bool{
		"in":      true,
		"version": true,
	}
	getCertificateCmd.PersistentFlags().String(getCertificate.InputFlagName, "", getCertificate.FileInFlagDescription)
	getCertificateCmd.PersistentFlags().String(getCertificate.VersionFlagName, "", getCertificate.VersionFlagDescription)
	getCertificateCmd.PersistentFlags().String(getCertificate.OutputFlagName, "", getCertificate.FileOutFlagDescription)
	common.SetCustomHelpTemplate(getCertificateCmd, requiredFlags)
}
