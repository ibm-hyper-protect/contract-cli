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

package base64

import (
	"fmt"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "base64"
	ParameterShortDescription = "Encode input as Base64"
	ParameterLongDescription  = `Encode text or JSON data to Base64 format.

Useful for encoding data that needs to be included in contracts or configurations.`

	InputFlagName        = "in"
	InputFlagDescription = "Input data to encode (text or JSON)"

	FormatFlagName        = "format"
	FormatFlagDescription = "Input data format (text or json)"

	OutputFlagName        = "out"
	OutputFlagDescription = "Path to save Base64 encoded output"

	inputMissingMessageBase64 = "input data is missing"
	invalidInputMessageBase64 = "invalid input format"
	successMessageBase64      = "successfully generated Base64"

	TextFormat = "text"
	JsonFormat = "json"
)

// ValidateInput - function to validate base64 command input
func ValidateInput(cmd *cobra.Command) (string, string, string, error) {
	inputData, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", "", err
	}

	if inputData == "" {
		err := fmt.Errorf("required flag '--in' is missing")
		common.SetMandatoryFlagError(cmd, err)
	}

	formatType, err := cmd.Flags().GetString(FormatFlagName)
	if err != nil {
		return "", "", "", err
	}

	outputPath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", "", err
	}

	return inputData, formatType, outputPath, nil
}

// Process - function to generate base64 of input
func Process(inputData, formatType string) (string, error) {
	var base64String string
	var err error

	if formatType == TextFormat {
		base64String, _, _, err = contract.HpcrText(inputData)
		if err != nil {
			return "", err
		}
	} else if formatType == JsonFormat {
		base64String, _, _, err = contract.HpcrJson(inputData)
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf(invalidInputMessageBase64)
	}

	return base64String, nil
}

// Output - function to print base64 or redirect to file
func Output(base64String, outputPath string) error {
	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, base64String)
		if err != nil {
			return err
		}
		fmt.Println(successMessageBase64)
	} else {
		fmt.Println(base64String)
	}

	return nil
}
