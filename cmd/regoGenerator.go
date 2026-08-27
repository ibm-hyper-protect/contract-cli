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

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-cli/lib/regoGenerator"
	"github.com/spf13/cobra"
)

// regoGeneratorCmd represents the rego-generator command
var regoGeneratorCmd = &cobra.Command{
	Use:   regoGenerator.ParameterName,
	Short: regoGenerator.ParameterShortDescription,
	Long:  regoGenerator.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		inputPath, outputPath, format, err := regoGenerator.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		policy, policyBase64, err := regoGenerator.Process(inputPath)
		if err != nil {
			log.Fatal(err)
		}

		err = regoGenerator.Output(outputPath, format, policy, policyBase64)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(regoGeneratorCmd)

	requiredFlags := map[string]bool{
		"in": true,
	}
	regoGeneratorCmd.PersistentFlags().String(regoGenerator.InputFlagName, "", regoGenerator.InputFlagDescription)
	regoGeneratorCmd.PersistentFlags().String(regoGenerator.OutputFlagName, "", regoGenerator.OutputFlagDescription)
	regoGeneratorCmd.PersistentFlags().String(regoGenerator.FormatFlagName, regoGenerator.FormatBase64, regoGenerator.FormatFlagDescription)
	common.SetCustomHelpTemplate(regoGeneratorCmd, requiredFlags)
	common.SetCustomErrorTemplate(regoGeneratorCmd)
}
