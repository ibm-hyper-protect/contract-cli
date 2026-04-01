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

package validateCertChain

import (
	"fmt"
	"strings"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "validate-cert-chain"
	ParameterShortDescription = "Validate certificate chain of trust"
	ParameterLongDescription  = `Validate that an encryption certificate has a valid chain of trust.

Verifies the complete certificate chain including leaf certificate, intermediate CA certificate,
and root CA certificate using OpenSSL verification. This ensures the certificate is properly
signed and trusted before using it for encryption operations.`
	CertFlagDescription         = "Path to encryption certificate file (use '-' for standard input)"
	IntermediateFlagDescription = "Path to intermediate CA certificate file"
	RootFlagDescription         = "Path to root CA certificate file"
	InputFlagName               = "cert"
	IntermediateFlagName        = "intermediate"
	RootFlagName                = "root"
)

// ValidateInput - function to validate cert chain inputs
func ValidateInput(cmd *cobra.Command) (string, string, string, error) {
	certPath, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", "", err
	}

	intermediatePath, err := cmd.Flags().GetString(IntermediateFlagName)
	if err != nil {
		return "", "", "", err
	}

	rootPath, err := cmd.Flags().GetString(RootFlagName)
	if err != nil {
		return "", "", "", err
	}

	requiredFlags := map[string]string{
		"--cert":         certPath,
		"--intermediate": intermediatePath,
		"--root":         rootPath,
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

	return certPath, intermediatePath, rootPath, nil
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

// GetIntermediateCertData - function to read intermediate certificate from file
func GetIntermediateCertData(intermediatePath string) (string, error) {
	if !common.CheckFileFolderExists(intermediatePath) {
		return "", fmt.Errorf("intermediate certificate file doesn't exist: %s", intermediatePath)
	}

	intermediate, err := common.ReadDataFromFile(intermediatePath)
	if err != nil {
		return "", err
	}

	return intermediate, nil
}

// GetRootCertData - function to read root certificate from file
func GetRootCertData(rootPath string) (string, error) {
	if !common.CheckFileFolderExists(rootPath) {
		return "", fmt.Errorf("root certificate file doesn't exist: %s", rootPath)
	}

	root, err := common.ReadDataFromFile(rootPath)
	if err != nil {
		return "", err
	}

	return root, nil
}
