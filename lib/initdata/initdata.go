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

package initdata

import (
	"fmt"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "initdata"
	ParameterShortDescription = "Gzip and Encoded initdata annotation"
	ParameterLongDescription  = `Gzip and Encoded initdata annotation`

	InputFlagName        = "in"
	InputFlagDescription = "Path of Signed and Encrypted contract (use '-' for standard input)"

	OutputFlagName        = "out"
	OutputFlagDescription = "Path to save Gzipped and encoded initdata value"
)

// ValidateInput - function to validate inputs of initdata
func ValidateInput(cmd *cobra.Command) (string, string, error) {
	inputData, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", err
	}

	if inputData == "" {
		err := fmt.Errorf("Error: required flag '--in' is missing")
		common.SetMandatoryFlagError(cmd, err)
	}

	// Validate stdin input
	common.ValidateStdinInput(cmd, inputData)

	outputPath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", err
	}
	return inputData, outputPath, nil
}

// GenerateInitdata - function to generate gzipped initdata
func GenerateInitdata(inputDataPath string) (string, error) {
	var inputData string
	var err error

	// Handle stdin input
	if inputDataPath == "-" {
		inputData, err = common.ReadDataFromStdin()
		if err != nil {
			return "", fmt.Errorf("unable to read input from standard input: %w", err)
		}
	} else {
		if !common.CheckFileFolderExists(inputDataPath) {
			return "", fmt.Errorf("the contract path doesn't exist")
		}
		inputData, err = common.ReadDataFromFile(inputDataPath)
		if err != nil {
			return "", err
		}
	}

	gzipInitdata, _, _, err := contract.HpccInitdata(inputData)
	if err != nil {
		return "", err
	}
	return gzipInitdata, nil
}

// PrintInitdata - function to print generated gzipped initdata value
func PrintInitdata(gzippedData, outputPath string) error {
	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, gzippedData)
		if err != nil {
			return err
		}
		fmt.Println("Successfully generated gzipped initdata annotation")
	} else {
		fmt.Println(gzippedData)
	}
	return nil
}
