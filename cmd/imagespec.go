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
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-cli/lib/imagespec"
	"github.com/spf13/cobra"
)

var imageSpecCmd = &cobra.Command{
	Use:   imagespec.ParameterName,
	Short: imagespec.ParameterShortDescription,
	Long:  imagespec.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := imagespec.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		result, err := imagespec.Process(input.ImageRef, input.ContainerName, input.Username, input.Password)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Image user: %s\n", result.ImageUser)

		if input.OutputPath != "" {
			err = common.WriteDataToFile(input.OutputPath, result.YAML)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Successfully wrote pod YAML template to %s\n", input.OutputPath)
		} else {
			fmt.Print(result.YAML)
		}
	},
}

func init() {
	rootCmd.AddCommand(imageSpecCmd)

	requiredFlags := map[string]bool{
		imagespec.InputFlagName: true,
	}

	imageSpecCmd.PersistentFlags().String(imagespec.InputFlagName, "", imagespec.ImageRefDescription)
	imageSpecCmd.PersistentFlags().String(imagespec.OutputFlagName, "", imagespec.OutputDescription)
	imageSpecCmd.PersistentFlags().String(imagespec.UserFlagName, "", imagespec.UsernameDescription)
	imageSpecCmd.PersistentFlags().String(imagespec.PassFlagName, "", imagespec.PasswordDescription)
	imageSpecCmd.PersistentFlags().String(imagespec.ContainerNameFlagName, "", imagespec.ContainerNameDescription)

	common.SetCustomHelpTemplate(imageSpecCmd, requiredFlags)
	common.SetCustomErrorTemplate(imageSpecCmd)
}
