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

package base64Tgz

import (
	"fmt"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/spf13/cobra"
)

const (
	// Base64 tar
	ParameterName             = "base64-tgz"
	ParameterShortDescription = "Create Base64 tar archive of container configurations"
	ParameterLongDescription  = `Generate Base64-encoded tar.gz archive of docker-compose.yaml or pods.yaml.

Creates a compressed archive of your container configuration files, encoded as Base64
for inclusion in Hyper Protect contracts. Supports both plain and encrypted output.`
	InputFlagName            = "in"
	InputFlagDescription     = "Path to folder containing docker-compose.yaml or pods.yaml"
	OutputFlagName           = "out"
	OutputFormatFlag         = "output"
	OutputFlagDescription    = "Output format (plain or encrypted)"
	OutputFormatUnencrypted  = "plain"
	OutputFormatEncrypted    = "encrypt"
	DefaultOutput            = OutputFormatUnencrypted
	OutputPathDescription    = "Path to save Base64 tar.gz output"
	OsVersionFlagName        = "os"
	OsVersionFlagDescription = "Target Hyper Protect platform (hpvs, hpcr-rhvs, or hpcc-peerpod)"
	CertFlagName             = "cert"
	CertPathDescription      = "Path to encryption certificate file"
)

// ValidateInput - function to validate base64-tgz inputs
func ValidateInput(cmd *cobra.Command) (string, string, string, string, string, error) {
	inputData, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	if inputData == "" {
		_ = cmd.Help()
		return "", "", "", "", "", fmt.Errorf("Error: required flag '--in' is missing.")
	}

	outputFormat, err := cmd.Flags().GetString(OutputFormatFlag)
	if err != nil {
		return "", "", "", "", "", err
	}

	hyperProtectVersion, err := cmd.Flags().GetString(OsVersionFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	encCertPath, err := cmd.Flags().GetString(CertFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	outputPath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	return inputData, outputFormat, hyperProtectVersion, encCertPath, outputPath, nil
}

// Process - function to process base64-tgz inputs
func Process(inputData, outputFormat, hyperProtectVersion, encCertPath string) (string, error) {
	if outputFormat == OutputFormatUnencrypted {
		if !common.CheckFileFolderExists(inputData) {
			return "", fmt.Errorf("the path to docker-compose.yaml or pods.yaml is not accessible")
		}

		base64Data, _, _, err := contract.HpcrTgz(inputData)
		if err != nil {
			return "", fmt.Errorf("failed to generate base64 tar - %v", err)
		}

		return base64Data, nil
	} else if outputFormat == OutputFormatEncrypted {
		encCert, err := common.GetDataFromFile(encCertPath)
		if err != nil {
			return "", err
		}

		encryptedBase64Data, _, _, err := contract.HpcrTgzEncrypted(inputData, hyperProtectVersion, encCert)
		if err != nil {
			return "", fmt.Errorf("failed to generate encrypted base64 tar - %v", err)
		}

		return encryptedBase64Data, nil
	} else {
		return "", fmt.Errorf("invalid output format (supported: plain / encrypt)")
	}
}

// Output - function to print base64 tgz or redirect it to a file
func Output(tarTgzData, outputPath string) error {
	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, tarTgzData)
		if err != nil {
			return err
		}
		fmt.Println("Successfully stored tar tgz data")
	} else {
		fmt.Println(tarTgzData)
	}

	return nil
}
