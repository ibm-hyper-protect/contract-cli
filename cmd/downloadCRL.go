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
	"github.com/ibm-hyper-protect/contract-cli/lib/downloadCRL"
	"github.com/ibm-hyper-protect/contract-go/v2/certificate"
	"github.com/spf13/cobra"
)

// downloadCRLCmd represents the download-crl command
var downloadCRLCmd = &cobra.Command{
	Use:   downloadCRL.ParameterName,
	Short: downloadCRL.ParameterShortDescription,
	Long:  downloadCRL.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		crlURL, outputPath, err := downloadCRL.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		crlData, err := certificate.HpcrDownloadCRL(crlURL)
		if err != nil {
			log.Fatal(err)
		}

		err = downloadCRL.Output(crlData, outputPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(downloadCRLCmd)

	requiredFlags := map[string]bool{
		"url": true,
	}
	downloadCRLCmd.PersistentFlags().String(downloadCRL.URLFlagName, "", downloadCRL.URLFlagDescription)
	downloadCRLCmd.PersistentFlags().String(downloadCRL.OutputFlagName, "", downloadCRL.OutputPathDescription)
	common.SetCustomHelpTemplate(downloadCRLCmd, requiredFlags)
	common.SetCustomErrorTemplate(downloadCRLCmd)
}
