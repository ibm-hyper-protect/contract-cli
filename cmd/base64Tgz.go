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

	"github.com/ibm-hyper-protect/contract-cli/lib/base64Tgz"
	"github.com/spf13/cobra"
)

// base64TgzCmd represents the base64-tgz command
var base64TgzCmd = &cobra.Command{
	Use:   base64Tgz.ParameterName,
	Short: base64Tgz.ParameterShortDescription,
	Long:  base64Tgz.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		inputData, outputFormat, hyperProtectVersion, encCert, outputPath, err := base64Tgz.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		base64TgzData, err := base64Tgz.Process(inputData, outputFormat, hyperProtectVersion, encCert)
		if err != nil {
			log.Fatal(err)
		}

		err = base64Tgz.Output(base64TgzData, outputPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(base64TgzCmd)

	base64TgzCmd.PersistentFlags().String(base64Tgz.InputFlagName, "", base64Tgz.InputFlagDescription)
	base64TgzCmd.PersistentFlags().String(base64Tgz.OutputFormatFlag, base64Tgz.DefaultOutput, base64Tgz.OutputFlagDescription)
	base64TgzCmd.PersistentFlags().String(base64Tgz.OsVersionFlagName, "", base64Tgz.OsVersionFlagDescription)
	base64TgzCmd.PersistentFlags().String(base64Tgz.CertFlagName, "", base64Tgz.CertPathDescription)
	base64TgzCmd.PersistentFlags().String(base64Tgz.OutputFlagName, "", base64Tgz.OutputPathDescription)
}
