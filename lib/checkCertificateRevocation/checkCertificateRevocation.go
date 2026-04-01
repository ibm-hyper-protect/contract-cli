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

package checkCertificateRevocation

import (
	"fmt"
	"strings"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "check-crl"
	ParameterShortDescription = "Check if certificate has been revoked"
	ParameterLongDescription  = `Check if an encryption certificate has been revoked using a Certificate Revocation List (CRL).

Verifies that the certificate has not been revoked by the issuing Certificate Authority.
This is a critical security check before using certificates for encryption operations.`
	CertFlagDescription = "Path to encryption certificate file (use '-' for standard input)"
	CRLFlagDescription  = "Path to Certificate Revocation List (CRL) file"
	InputFlagName       = "cert"
	CRLFlagName         = "crl"
)

// ValidateInput - function to validate revocation check inputs
func ValidateInput(cmd *cobra.Command) (string, string, error) {
	certPath, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", err
	}

	crlPath, err := cmd.Flags().GetString(CRLFlagName)
	if err != nil {
		return "", "", err
	}

	requiredFlags := map[string]string{
		"--cert": certPath,
		"--crl":  crlPath,
	}

	var missing []string
	for flag, val := range requiredFlags {
		if val == "" {
			missing = append(missing, flag)
		}
	}

	if len(missing) > 0 {
		if len(missing) == 1 {
			err := fmt.Errorf("Error: required flag %s is missing",
				strings.Join(missing, ", "))
			common.SetMandatoryFlagError(cmd, err)
		} else {
			err := fmt.Errorf("Error: required flags %s are missing",
				strings.Join(missing, ", "))
			common.SetMandatoryFlagError(cmd, err)
		}
	}

	// Validate stdin input for cert
	common.ValidateStdinInput(cmd, certPath)

	return certPath, crlPath, nil
}

// GetCertificateData - function to read certificate from file or stdin
func GetCertificateData(certPath string) (string, error) {
	var cert string
	var err error

	// Handle stdin input for certificate
	if certPath == "-" {
		cert, err = common.ReadDataFromStdin()
		if err != nil {
			return "", fmt.Errorf("unable to read certificate from standard input: %w", err)
		}
	} else {
		if !common.CheckFileFolderExists(certPath) {
			return "", fmt.Errorf("certificate file doesn't exist: %s", certPath)
		}
		cert, err = common.ReadDataFromFile(certPath)
		if err != nil {
			return "", err
		}
	}

	return cert, nil
}

// GetCRLData - function to read CRL from file
func GetCRLData(crlPath string) (string, error) {
	if !common.CheckFileFolderExists(crlPath) {
		return "", fmt.Errorf("CRL file doesn't exist: %s", crlPath)
	}

	crl, err := common.ReadDataFromFile(crlPath)
	if err != nil {
		return "", err
	}

	return crl, nil
}
