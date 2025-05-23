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
	"github.com/ibm-hyper-protect/contract-go/contract"
	"github.com/spf13/cobra"
)

// validateContractCmd represents the validate-contract command
var validateContractCmd = &cobra.Command{
	Use:   common.ValidateContractParamName,
	Short: common.ValidateContractParamShortDescription,
	Long:  common.ValidateContractParamLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		contractPath, version, err := validateInputContract(cmd)
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

	validateContractCmd.PersistentFlags().String(common.FileInFlagName, "", common.ValidateContractInputFlagDescription)
	validateContractCmd.PersistentFlags().String(common.OsVersionFlagName, "", common.OsVersionFlagDescription)
}

// validateInputContract - function to validate plain contract
func validateInputContract(cmd *cobra.Command) (string, string, error) {
	contract, err := cmd.Flags().GetString(common.FileInFlagName)
	if err != nil {
		return "", "", err
	}

	version, err := cmd.Flags().GetString(common.OsVersionFlagName)
	if err != nil {
		return "", "", err
	}

	return contract, version, nil
}
