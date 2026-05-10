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
	"github.com/ibm-hyper-protect/contract-cli/lib/listCertVersions"
	"github.com/spf13/cobra"
)

// listCertVersionsCmd represents the list-cert-versions command
var listCertVersionsCmd = &cobra.Command{
	Use:   listCertVersions.ParameterName,
	Short: listCertVersions.ParameterShortDescription,
	Long:  listCertVersions.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		osVersion, outputPath, format, err := listCertVersions.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		result, err := listCertVersions.Process(osVersion, format)
		if err != nil {
			log.Fatal(err)
		}

		formattedResult := listCertVersions.Output(result)

		if outputPath != "" {
			err = common.WriteDataToFile(outputPath, formattedResult)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Successfully stored certificate versions")
		} else {
			fmt.Println(formattedResult)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCertVersionsCmd)

	listCertVersionsCmd.Flags().StringP(
		listCertVersions.OsVersionFlagName,
		"",
		"",
		listCertVersions.OsVersionFlagDescription,
	)

	listCertVersionsCmd.Flags().StringP(
		listCertVersions.OutputFlagName,
		"",
		"",
		listCertVersions.OutputFlagDescription,
	)

	listCertVersionsCmd.Flags().StringP(
		listCertVersions.FormatFlagName,
		"",
		"",
		listCertVersions.FormatFlagDescription,
	)
}

// Made with Bob
