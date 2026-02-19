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

package validateEncryptionCertificate

import (
	"fmt"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/spf13/cobra"
)

const (
	// Validate Encryption certificates
	ParameterName             = "validate-encryption-certificate"
	ParameterShortDescription = "validate encryption certificate"
	ParameterLongDescription  = `validate encryption certificate for HPCR from IBM Hyper Protect Repository

Validates encryption certificate for on-premise, VPC deployment.
It will check encryption certificate validity`
	InputFlagDescription       = "Path to encryption certificate file (use '-' for standard input)"
	CertVersionFlagDescription = "Versions of Encryption Certificates to validate, Seperated by coma(,)"
	InputFlagName              = "in"
)

// ValidateInput - function to validate encryption certificate input
func ValidateInput(cmd *cobra.Command) (string, error) {
	encryptionCertsPath, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", err
	}
	if encryptionCertsPath == "" {
		err := fmt.Errorf("Error: required flag '--in' is missing")
		common.SetMandatoryFlagError(cmd, err)
	}

	// Validate stdin input
	common.ValidateStdinInput(cmd, encryptionCertsPath)

	return encryptionCertsPath, nil
}

// GetEncryptionCertfile - function to get encryption certificate
func GetEncryptionCertfile(encryptionCertsPath string) (string, error) {
	var encryptionCert string
	var err error

	// Handle stdin input
	if encryptionCertsPath == "-" {
		encryptionCert, err = common.ReadDataFromStdin()
		if err != nil {
			return "", fmt.Errorf("unable to read input from standard input: %w", err)
		}
	} else {
		if !common.CheckFileFolderExists(encryptionCertsPath) {
			return "", fmt.Errorf("The specified path does not contain the encryption certificates.")
		}
		encryptionCert, err = common.ReadDataFromFile(encryptionCertsPath)
		if err != nil {
			return "", err
		}
	}

	return encryptionCert, nil
}
