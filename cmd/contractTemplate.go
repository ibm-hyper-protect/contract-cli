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
	"github.com/ibm-hyper-protect/contract-cli/lib/contractTemplate"
	"github.com/spf13/cobra"
)

// contractTemplateCmd represents the contract-template command
var contractTemplateCmd = &cobra.Command{
	Use:   contractTemplate.ParameterName,
	Short: contractTemplate.ParameterShortDescription,
	Long:  contractTemplate.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		templateType, osVersion, outputPath, err := contractTemplate.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		template, err := contractTemplate.GenerateContractTemplate(templateType, osVersion)
		if err != nil {
			log.Fatal(err)
		}

		err = contractTemplate.Output(template, outputPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(contractTemplateCmd)

	requiredFlags := map[string]bool{}

	contractTemplateCmd.PersistentFlags().String(contractTemplate.TypeFlagName, contractTemplate.TypeContract, contractTemplate.TypeFlagDescription)
	contractTemplateCmd.PersistentFlags().String(contractTemplate.OsVersionFlagName, contractTemplate.OsHpvs, contractTemplate.OsVersionFlagDescription)
	contractTemplateCmd.PersistentFlags().String(contractTemplate.OutputFlagName, "", contractTemplate.OutputFlagDescription)

	common.SetCustomHelpTemplate(contractTemplateCmd, requiredFlags)
	common.SetCustomErrorTemplate(contractTemplateCmd)
}
