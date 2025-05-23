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
	"encoding/json"
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/image"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type ImageDetails struct {
	Id       string `json:"id"       yaml:"id"`
	Name     string `json:"name"     yaml:"name"`
	Checksum string `json:"checksum" yaml:"checksum"`
	Version  string `json:"version"  yaml:"version"`
}

const (
	invalidImagePathMessage = "The Image details path doesn't exists"
	invalidFormatMessage    = "invalid output format"
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   common.ImageParamName,
	Short: common.ImageParamShortDescription,
	Long:  common.ImageParamLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		imageListJsonPath, versionName, formatType, hpcrImagePath, err := validateInputImage(cmd)
		if err != nil {
			log.Fatal(err)
		}

		imageDetail, err := getImageDetails(imageListJsonPath, versionName)
		if err != nil {
			log.Fatal(err)
		}

		result, err := printDataImage(imageDetail, formatType)
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

	imageCmd.PersistentFlags().String(common.FileInFlagName, "", common.IbmCloudJsonInputDescription)
	imageCmd.PersistentFlags().String(common.VersionFlagName, "", common.HpcrVersionFlagDescription)
	imageCmd.PersistentFlags().String(common.DataFormatFlagName, common.DataFormatDefault, common.DataFormatFlagDescription)
	imageCmd.PersistentFlags().String(common.FileOutFlagName, "", common.HpcrImageFlagDescription)
}

// validateInputImage - function to validate image input
func validateInputImage(cmd *cobra.Command) (string, string, string, string, error) {
	imageListJsonPath, err := cmd.Flags().GetString(common.FileInFlagName)
	if err != nil {
		return "", "", "", "", err
	}
	versionName, err := cmd.Flags().GetString(common.VersionFlagName)
	if err != nil {
		return "", "", "", "", err
	}
	formatType, err := cmd.Flags().GetString(common.DataFormatFlagName)
	if err != nil {
		return "", "", "", "", err
	}
	hpcrImagePath, err := cmd.Flags().GetString(common.FileOutFlagName)
	if err != nil {
		return "", "", "", "", err
	}

	return imageListJsonPath, versionName, formatType, hpcrImagePath, nil
}

// getImageDetails - function to get HPCR image details from JSON input
func getImageDetails(imageDetailsJsonPath, versionName string) (ImageDetails, error) {
	if !common.CheckFileFolderExists(imageDetailsJsonPath) {
		log.Fatal(invalidImagePathMessage)
	}

	imageDataJson, err := common.ReadDataFromFile(imageDetailsJsonPath)
	if err != nil {
		log.Fatal(err)
	}

	imageId, imageName, imageChecksum, ImageVersion, err := image.HpcrSelectImage(imageDataJson, versionName)
	if err != nil {
		log.Fatal(err)
	}

	return ImageDetails{imageId, imageName, imageChecksum, ImageVersion}, nil
}

// printDataImage - function to print image details or redirect output to a file
func printDataImage(imageDetail ImageDetails, format string) (string, error) {
	if format == common.DataFormatJson {
		imageJson, err := json.MarshalIndent(imageDetail, "", "  ")
		if err != nil {
			return "", err
		}
		return string(imageJson), nil
	} else if format == common.DataFormatYaml {
		imageYaml, err := yaml.Marshal(imageDetail)
		if err != nil {
			return "", err
		}
		return string(imageYaml), nil
	} else {
		return "", fmt.Errorf(invalidFormatMessage)
	}
}
