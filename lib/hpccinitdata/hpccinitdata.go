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

package hpccinitdata

import (
	"fmt"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "hpcc-init"
	ParameterShortDescription = "Gzip and Encoded hpcc initdata annotation"
	ParameterLongDescription  = `Gzip and Encoded initdata annotation for HyperProtect Confidential Containers`

	InputFlagName        = "in"
	InputFlagDescription = "Path of Signed and Encrypted contract"

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

	outputPath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", err
	}
	return inputData, outputPath, nil
}

// GenerateGzippedInitdata - function to generate gzipped initdata
func GenerateHpccInitdata(inputDataPath string) (string, error) {
	if !common.CheckFileFolderExists(inputDataPath) {
		return "", fmt.Errorf("the contract path doesn't exist")
	}
	inputData, err := common.ReadDataFromFile(inputDataPath)
	if err != nil {
		return "", err
	}
	gzipInitdata, _, _, err := contract.HpccInitdata(inputData)
	if err != nil {
		return "", err
	}
	return gzipInitdata, nil
}

// PrintHpccInitdata - function to print generated gzipped initdata value
func PrintHpccInitdata(gzippedData, outputPath string) error {
	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, gzippedData)
		if err != nil {
			return err
		}
	} else {
		fmt.Println(gzippedData)
	}
	fmt.Println("Successfully generated gzipped initdata for HPCC")
	return nil
}
