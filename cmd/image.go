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
	"github.com/ibm-hyper-protect/contract-cli/lib/image"
	"github.com/spf13/cobra"
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   image.ParameterName,
	Short: image.ParameterShortDescription,
	Long:  image.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		imageListJsonPath, versionName, formatType, hpcrImagePath, err := image.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		imageDetail, err := image.Process(imageListJsonPath, versionName)
		if err != nil {
			log.Fatal(err)
		}

		result, err := image.Output(imageDetail, formatType)
		if err != nil {
			log.Fatal(err)
		}

		if hpcrImagePath != "" {
			err = common.WriteDataToFile(hpcrImagePath, result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Successfully stored HPCR image details")
		} else {
			fmt.Println(result)
		}
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(imageCmd)

	requiredFlags := map[string]bool{
		"in": true,
	}
	imageCmd.PersistentFlags().String(image.InputFlagName, "", image.IbmCloudJsonInputDescription)
	imageCmd.PersistentFlags().String(image.VersionFlagName, "", image.HpcrVersionFlagDescription)
	imageCmd.PersistentFlags().String(image.FormatFlag, image.JsonFormat, image.DataFormatFlagDescription)
	imageCmd.PersistentFlags().String(image.OutputFlagName, "", image.OutputFlagDescription)
	common.SetCustomHelpTemplate(imageCmd, requiredFlags)
	common.SetCustomErrorTemplate(imageCmd)
}
