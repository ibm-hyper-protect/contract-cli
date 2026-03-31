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
	"github.com/ibm-hyper-protect/contract-cli/lib/signContract"
	"github.com/spf13/cobra"
)

// signContractCmd represents the signContract command
var signContractCmd = &cobra.Command{
	Use:   signContract.ParameterName,
	Short: signContract.ParameterShortDescription,
	Long:  signContract.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		contract, privateKey, output, password, err := signContract.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		contractSign, err := signContract.GenerateSignContract(contract, privateKey, password)
		if err != nil {
			log.Fatal(err)
		}

		err = signContract.Output(contractSign, output)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(signContractCmd)

	requiredFlags := map[string]bool{
		"in": true,
	}

	signContractCmd.PersistentFlags().String(signContract.InputFlagName, "", signContract.InputFlagDescription)
	signContractCmd.PersistentFlags().String(signContract.PrivateKeyFlagName, "", signContract.PrivateKeyFlagDescription)
	signContractCmd.PersistentFlags().String(signContract.PasswordFlagName, "", signContract.PasswordFlagDescription)
	signContractCmd.PersistentFlags().String(signContract.OutputFlagName, "", signContract.OutputFlagDescription)

	common.SetCustomHelpTemplate(signContractCmd, requiredFlags)
	common.SetCustomErrorTemplate(signContractCmd)
}
