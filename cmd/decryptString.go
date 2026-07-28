// Copyright (c) 2026 IBM Corp.
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

	"github.com/spf13/cobra"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-cli/lib/decryptString"
)

// decryptStringCmd represents the decrypt command
var decryptStringCmd = &cobra.Command{
	Use:   decryptString.ParameterName,
	Short: decryptString.ParameterShortDescription,
	Long:  decryptString.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		inputData, privateKeyPath, password, outputPath, err := decryptString.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		decryptedText, err := decryptString.Process(inputData, privateKeyPath, password)
		if err != nil {
			log.Fatal(err)
		}

		err = decryptString.Output(outputPath, decryptedText)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(decryptStringCmd)

	requiredFlags := map[string]bool{
		"in":   true,
		"priv": true,
	}
	decryptStringCmd.PersistentFlags().String(decryptString.InputFlagName, "", decryptString.InputFlagDescription)
	decryptStringCmd.PersistentFlags().String(decryptString.PrivateKeyFlagName, "", decryptString.PrivateKeyFlagDescription)
	decryptStringCmd.PersistentFlags().String(decryptString.PasswordFlagName, "", decryptString.PasswordFlagDescription)
	decryptStringCmd.PersistentFlags().String(decryptString.OutputFlagName, "", decryptString.OutputFlagDescription)
	common.SetCustomHelpTemplate(decryptStringCmd, requiredFlags)
	common.SetCustomErrorTemplate(decryptStringCmd)
}
