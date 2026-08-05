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

package decryptString

import (
	"fmt"
	"strings"

	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/spf13/cobra"

	"github.com/ibm-hyper-protect/contract-cli/common"
)

const (
	ParameterName             = "decrypt"
	ParameterShortDescription = "Decrypt encrypted text in IBM Confidential Computing format"
	ParameterLongDescription  = `Decrypt encrypted strings in IBM Confidential Computing format using an RSA private key.

Supports both encryption formats:
  - contract-basic.<encrypted-password>.<encrypted-data>  (CCRT/CCRV)
  - hyper-protect-basic.<encrypted-password>.<encrypted-data>  (CCCO/HPVS)

The private key must correspond to the key used during encryption.`

	InputFlagName             = "in"
	InputFlagDescription      = "Path to encrypted input file or encrypted string (use '-' for standard input)"
	PrivateKeyFlagName        = "priv"
	PrivateKeyFlagDescription = "Path to RSA private key file (PEM format)"
	PasswordFlagName          = "password"
	PasswordFlagDescription   = "Password for encrypted private key"
	OutputFlagName            = "out"
	OutputFlagDescription     = "Path to save decrypted output (prints to stdout if omitted)"

	successMessage = "Successfully stored decrypted text"
)

// ValidateInput - function to validate decrypt inputs
func ValidateInput(cmd *cobra.Command) (string, string, string, string, error) {
	inputData, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", "", "", err
	}

	privateKeyPath, err := cmd.Flags().GetString(PrivateKeyFlagName)
	if err != nil {
		return "", "", "", "", err
	}

	password, err := cmd.Flags().GetString(PasswordFlagName)
	if err != nil {
		return "", "", "", "", err
	}

	outputPath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", "", "", err
	}

	requiredFlags := map[string]string{
		"--in":   inputData,
		"--priv": privateKeyPath,
	}

	var missing []string
	for flag, val := range requiredFlags {
		if val == "" {
			missing = append(missing, flag)
		}
	}

	if len(missing) > 0 {
		if len(missing) == 1 {
			err := fmt.Errorf("Error: required flag %s is missing.", strings.Join(missing, ", "))
			common.SetMandatoryFlagError(cmd, err)
		} else {
			err := fmt.Errorf("Error: required flags %s are missing.", strings.Join(missing, ", "))
			common.SetMandatoryFlagError(cmd, err)
		}
	}

	// Validate stdin input
	common.ValidateStdinInput(cmd, inputData)

	return inputData, privateKeyPath, password, outputPath, nil
}

// Process - function to decrypt encrypted text
func Process(inputData, privateKeyPath, password string) (string, error) {
	var encryptedText string
	var err error

	// Handle stdin input
	if inputData == "-" {
		encryptedText, err = common.ReadDataFromStdin()
		if err != nil {
			return "", fmt.Errorf("unable to read input from standard input: %w", err)
		}
	} else {
		// Try as file path first, fall back to treating as raw string
		if common.CheckFileFolderExists(inputData) {
			encryptedText, err = common.ReadDataFromFile(inputData)
			if err != nil {
				return "", err
			}
		} else {
			encryptedText = inputData
		}
	}

	if privateKeyPath == "" {
		return "", fmt.Errorf("private key path is required")
	}

	if !common.CheckFileFolderExists(privateKeyPath) {
		return "", fmt.Errorf("the path to private key doesn't exist: %s", privateKeyPath)
	}

	privateKey, err := common.ReadDataFromFile(privateKeyPath)
	if err != nil {
		return "", err
	}

	decryptedText, _, _, err := contract.HpcrTextDecrypted(encryptedText, privateKey, password)
	if err != nil {
		return "", err
	}

	return decryptedText, nil
}

// Output - function to print decrypted data or redirect output to a file
func Output(outputPath, decryptedText string) error {
	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, decryptedText)
		if err != nil {
			return err
		}
		fmt.Println(successMessage)
	} else {
		fmt.Println(decryptedText)
	}

	return nil
}
