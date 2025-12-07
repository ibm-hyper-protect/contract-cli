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

	"github.com/ibm-hyper-protect/contract-cli/lib/base64"
	"github.com/spf13/cobra"
)

// base64Cmd represents the base64 command
var base64Cmd = &cobra.Command{
	Use:   base64.ParameterName,
	Short: base64.ParameterShortDescription,
	Long:  base64.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		inputData, formatType, outputPath, err := base64.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		base64String, err := base64.Process(inputData, formatType)
		if err != nil {
			log.Fatal(err)
		}

		err = base64.Output(base64String, outputPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(base64Cmd)

	base64Cmd.PersistentFlags().String(base64.InputFlagName, "", base64.InputFlagDescription)
	base64Cmd.PersistentFlags().String(base64.FormatFlagName, base64.TextFormat, base64.FormatFlagDescription)
	base64Cmd.PersistentFlags().String(base64.OutputFlagName, "", base64.OutputFlagDescription)
}
