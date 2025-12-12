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
	"github.com/ibm-hyper-protect/contract-cli/lib/encryptString"
	"github.com/spf13/cobra"
)

// encryptStringCmd represents the encrypt-string command
var encryptStringCmd = &cobra.Command{
	Use:   encryptString.ParameterName,
	Short: encryptString.ParameterShortDescription,
	Long:  encryptString.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		inputData, inputFormat, hyperProtectVersion, encCertPath, outputPath, err := encryptString.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		encryptedString, err := encryptString.Process(inputData, inputFormat, hyperProtectVersion, encCertPath)
		if err != nil {
			log.Fatal(err)
		}

		err = encryptString.Output(outputPath, encryptedString)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(encryptStringCmd)

	requiredFlags := map[string]bool{
		"in": true,
	}
	encryptStringCmd.PersistentFlags().String(encryptString.InputFlagName, "", encryptString.InputFlagDescription)
	encryptStringCmd.PersistentFlags().String(encryptString.FormatFlag, encryptString.TextFormat, encryptString.FormatFlagDescription)
	encryptStringCmd.PersistentFlags().String(encryptString.OsVersionFlagName, "", encryptString.OsVersionFlagDescription)
	encryptStringCmd.PersistentFlags().String(encryptString.CertFlagName, "", encryptString.CertFlagDescription)
	encryptStringCmd.PersistentFlags().String(encryptString.OutputFlagName, "", encryptString.OutputFlagDescription)
	common.SetCustomHelpTemplate(encryptStringCmd, requiredFlags)
}
