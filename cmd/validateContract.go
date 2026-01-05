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
	"github.com/ibm-hyper-protect/contract-cli/lib/validateContract"
	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/spf13/cobra"
)

// validateContractCmd represents the validate-contract command
var validateContractCmd = &cobra.Command{
	Use:   validateContract.ParameterName,
	Short: validateContract.ParameterShortDescription,
	Long:  validateContract.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		contractPath, version, err := validateContract.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		contractData, err := common.ReadDataFromFile(contractPath)
		if err != nil {
			log.Fatal(err)
		}

		if !common.CheckFileFolderExists(contractPath) {
			log.Fatal("The path to contract doesn't exist")
		}

		err = contract.HpcrVerifyContract(contractData, version)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("True")
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(validateContractCmd)

	requiredFlags := map[string]bool{
		"in": true,
	}
	validateContractCmd.PersistentFlags().String(validateContract.InputFlagName, "", validateContract.InputFlagDescription)
	validateContractCmd.PersistentFlags().String(validateContract.OsVersionFlagName, "", validateContract.OsVersionFlagDescription)
	common.SetCustomHelpTemplate(validateContractCmd, requiredFlags)
	common.SetCustomErrorTemplate(validateContractCmd)
}
