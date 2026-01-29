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
	"github.com/ibm-hyper-protect/contract-cli/lib/hpccinitdata"
	"github.com/spf13/cobra"
)

// initdataCmd represents the hpccinidata command
var initdataCmd = &cobra.Command{
	Use:   hpccinitdata.ParameterName,
	Short: hpccinitdata.ParameterShortDescription,
	Long:  hpccinitdata.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		inputDataPath, outputPath, err := hpccinitdata.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		gzipInitdata, err := hpccinitdata.GenerateHpccInitdata(inputDataPath)
		if err != nil {
			log.Fatal(err)
		}

		err = hpccinitdata.PrintHpccInitdata(gzipInitdata, outputPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(initdataCmd)
	requiredFlags := map[string]bool{
		"in": true,
	}

	initdataCmd.PersistentFlags().String(hpccinitdata.InputFlagName, "", hpccinitdata.InputFlagDescription)
	initdataCmd.PersistentFlags().String(hpccinitdata.OutputFlagName, "", hpccinitdata.OutputFlagDescription)
	common.SetCustomHelpTemplate(initdataCmd, requiredFlags)
	common.SetCustomErrorTemplate(initdataCmd)
}
