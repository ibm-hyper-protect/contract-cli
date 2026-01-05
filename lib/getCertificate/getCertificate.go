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

package getCertificate

import (
	"fmt"
	"strings"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/v2/certificate"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "get-certificate"
	ParameterShortDescription = "Extract specific certificate version from download output"
	ParameterLongDescription  = `Extract a specific encryption certificate version from download-certificate output.

Parses the JSON output from download-certificate and extracts the certificate
for the specified version.`
	FileInFlagDescription  = "Path to download-certificate JSON output"
	VersionFlagDescription = "Certificate version to extract (e.g., 1.0.23)"
	FileOutFlagDescription = "Path to save extracted encryption certificate"
	InputFlagName          = "in"
	OutputFlagName         = "out"
	VersionFlagName        = "version"
)

// ValidateInput - function to validate get-certificate input
func ValidateInput(cmd *cobra.Command) (string, string, string, error) {
	encryptionCertsPath, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", "", err
	}

	version, err := cmd.Flags().GetString(VersionFlagName)
	if err != nil {
		return "", "", "", err
	}

	encryptionCertificatePath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", "", err
	}

	requiredFlags := map[string]string{
		"--in":      encryptionCertsPath,
		"--version": version,
	}

	var missing []string
	for flag, val := range requiredFlags {
		if val == "" {
			missing = append(missing, flag)
		}
	}

	if len(missing) > 0 {
		if len(missing) == 1 {
			err := fmt.Errorf("Error: required flag %s is missing.",
				strings.Join(missing, ", "))
			common.SetMandatoryFlagError(cmd, err)
		} else {
			err := fmt.Errorf("Error: required flag %s are missing.",
				strings.Join(missing, ", "))
			common.SetMandatoryFlagError(cmd, err)
		}
	}

	return encryptionCertsPath, version, encryptionCertificatePath, nil
}

// Process - function to get encryption certificate
func Process(encryptionCertsPath, version string) (string, error) {
	if !common.CheckFileFolderExists(encryptionCertsPath) {
		return "", fmt.Errorf("the path to encryption certificates doesn't exist")
	}

	encryptionCertsJson, err := common.ReadDataFromFile(encryptionCertsPath)
	if err != nil {
		return "", err
	}

	_, outputCertificate, err := certificate.HpcrGetEncryptionCertificateFromJson(encryptionCertsJson, version)
	if err != nil {
		return "", err
	}

	return outputCertificate, nil
}

// Output - function to print encryption certificate or redirect it to a file
func Output(cert, certPath string) error {
	if certPath != "" {
		err := common.WriteDataToFile(certPath, cert)
		if err != nil {
			return err
		}
		fmt.Println("Successfully added encryption certificate ")
	} else {
		fmt.Println(cert)
	}

	return nil
}
