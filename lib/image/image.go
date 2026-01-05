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

package image

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/v2/image"
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
	ParameterName             = "image"
	ParameterShortDescription = "Get HPCR image details from IBM Cloud"
	ParameterLongDescription  = `Retrieve Hyper Protect Container Runtime (HPCR) image details from IBM Cloud.

Parses image information from IBM Cloud API, CLI, or Terraform output to extract
image ID, name, checksum, and version. Supports filtering by specific HPCR version.`
	IbmCloudJsonInputDescription = "Path to IBM Cloud images JSON (from API, CLI, or Terraform)"
	HpcrVersionFlagDescription   = "Specific HPCR version to retrieve (returns latest if not specified)"
	OutputFlagDescription        = "Path to save HPCR image details"
	invalidImagePathMessage      = "The Image details path doesn't exists"
	invalidFormatMessage         = "invalid output format"
	InputFlagName                = "in"
	OutputFlagName               = "out"
	VersionFlagName              = "version"
	FormatFlag                   = "format"
	YamlFormat                   = "yaml"
	JsonFormat                   = "json"
	DataFormatFlagDescription    = "Output format for data (json, yaml, or text)"
)

// ValidateInput - function to validate image input
func ValidateInput(cmd *cobra.Command) (string, string, string, string, error) {
	imageListJsonPath, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", "", "", err
	}
	if imageListJsonPath == "" {
		err := fmt.Errorf("Error: required flag '--in' is missing")
		common.SetMandatoryFlagError(cmd, err)
	}

	versionName, err := cmd.Flags().GetString(VersionFlagName)
	if err != nil {
		return "", "", "", "", err
	}
	formatType, err := cmd.Flags().GetString(FormatFlag)
	if err != nil {
		return "", "", "", "", err
	}
	hpcrImagePath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", "", "", err
	}

	return imageListJsonPath, versionName, formatType, hpcrImagePath, nil
}

// Process - function to get HPCR image details from JSON input
func Process(imageDetailsJsonPath, versionName string) (ImageDetails, error) {
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

// Output - function to print image details or redirect output to a file
func Output(imageDetail ImageDetails, format string) (string, error) {
	if format == JsonFormat {
		imageJson, err := json.MarshalIndent(imageDetail, "", "  ")
		if err != nil {
			return "", err
		}
		return string(imageJson), nil
	} else if format == YamlFormat {
		imageYaml, err := yaml.Marshal(imageDetail)
		if err != nil {
			return "", err
		}
		return string(imageYaml), nil
	} else {
		return "", fmt.Errorf(invalidFormatMessage)
	}
}
