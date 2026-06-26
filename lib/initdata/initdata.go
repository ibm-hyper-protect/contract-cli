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

	SehdrBinFlagName        = "sehdr"
	SehdrBinFlagDescription = "Path to SE header binary file (.bin) for baremetal solution"

	OutputFlagName        = "out"
	OutputFlagDescription = "Path to save Gzipped and encoded initdata value"
)

// ValidateInput - function to validate inputs of initdata
func ValidateInput(cmd *cobra.Command) (string, string, string, error) {
	inputData, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", "", err
	}

	if inputData == "" {
		err := fmt.Errorf("Error: required flag '--in' is missing")
		common.SetMandatoryFlagError(cmd, err)
	}

	// Validate stdin input
	common.ValidateStdinInput(cmd, inputData)

	sehdrBinPath, err := cmd.Flags().GetString(SehdrBinFlagName)
	if err != nil {
		return "", "", "", err
	}

	outputPath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", "", err
	}
	return inputData, sehdrBinPath, outputPath, nil
}

// GenerateInitdata - function to generate gzipped initdata
func GenerateInitdata(inputDataPath, sehdrBinPath string) (string, bool, error) {
	var inputData string
	var err error

	// Handle stdin input
	if inputDataPath == "-" {
		inputData, err = common.ReadDataFromStdin()
		if err != nil {
			return "", false, fmt.Errorf("unable to read input from standard input: %w", err)
		}
	} else {
		if !common.CheckFileFolderExists(inputDataPath) {
			return "", false, fmt.Errorf("the contract path doesn't exist")
		}
		inputData, err = common.ReadDataFromFile(inputDataPath)
		if err != nil {
			return "", false, err
		}
	}

	// Handle SE header binary file if provided
	var sehdrBase64 string
	isBaremetal := false
	if sehdrBinPath != "" {
		if !common.CheckFileFolderExists(sehdrBinPath) {
			return "", false, fmt.Errorf("the SE header binary file path doesn't exist")
		}
		binData, err := common.ReadDataFromFile(sehdrBinPath)
		if err != nil {
			return "", false, fmt.Errorf("unable to read SE header binary file: %w", err)
		}
		// Generate base64 from binary data using HpcrText
		sehdrBase64, _, _, err = contract.HpcrText(binData)
		if err != nil {
			return "", false, fmt.Errorf("unable to encode SE header binary to base64: %w", err)
		}
		isBaremetal = true
	}

	gzipInitdata, _, _, err := contract.HpccInitdata(inputData, sehdrBase64)
	if err != nil {
		return "", false, err
	}
	return gzipInitdata, isBaremetal, nil
}

// PrintInitdata - function to print generated gzipped initdata value
func PrintInitdata(gzippedData, outputPath string, isBaremetal bool) error {
	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, gzippedData)
		if err != nil {
			return err
		}
		if isBaremetal {
			fmt.Println("Successfully generated gzipped initdata annotation for baremetal solution")
		} else {
			fmt.Println("Successfully generated gzipped initdata annotation for peerpod solution")
		}
	} else {
		fmt.Println(gzippedData)
	}
	return nil
}
