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

package encryptString

import (
	"fmt"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "encrypt-string"
	ParameterShortDescription = "Encrypt string in Hyper Protect format"
	ParameterLongDescription  = `Encrypt strings using the Hyper Protect encryption format.

Output format: hyper-protect-basic.<encrypted-password>.<encrypted-string>
Use this to encrypt sensitive data like passwords or API keys for contracts.`
	InputFlagDescription        = "String data to encrypt (text or JSON)"
	FormatFlagDescription       = "Input data format (text or json)"
	OutputFlagDescription       = "Path to save encrypted output"
	successMessageEncryptString = "Successfully stored encrypted text"
	InputFlagName               = "in"
	OutputFlagName              = "out"
	FormatFlag                  = "format"
	TextFormat                  = "text"
	JsonFormat                  = "json"
	OsVersionFlagName           = "os"
	OsVersionFlagDescription    = "Target Hyper Protect platform (hpvs, hpcr-rhvs, or hpcc-peerpod)"
	CertFlagName                = "cert"
	CertFlagDescription         = "Path to encryption certificate file"
)

// ValidateInput - function to validate encrypt-string inputs
func ValidateInput(cmd *cobra.Command) (string, string, string, string, string, error) {
	inputData, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}
	if inputData == "" {
		_ = cmd.Help()
		return "", "", "", "", "", fmt.Errorf("Error: required flag '--in' is missing.")
	}

	inputFormat, err := cmd.Flags().GetString(FormatFlag)
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

	return inputData, inputFormat, hyperProtectVersion, encCertPath, outputPath, nil
}

// Process - function to generate encrypted string of plain or JSON text
func Process(inputData, inputFormat, hyperProtectVersion, encCertPath string) (string, error) {
	encCert, err := common.GetDataFromFile(encCertPath)
	if err != nil {
		return "", err
	}

	var encryptedString string
	if inputFormat == TextFormat {
		encryptedString, _, _, err = contract.HpcrTextEncrypted(inputData, hyperProtectVersion, encCert)
		if err != nil {
			return "", err
		}
	} else if inputFormat == JsonFormat {
		encryptedString, _, _, err = contract.HpcrJsonEncrypted(inputData, hyperProtectVersion, encCert)
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("invalid input format")
	}

	return encryptedString, nil
}

// Output - function to print encrypted data or redirect output to a file
func Output(outputPath, encryptedString string) error {
	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, encryptedString)
		if err != nil {
			return err
		}
		fmt.Println(successMessageEncryptString)
	} else {
		fmt.Println(encryptedString)
	}

	return nil
}
