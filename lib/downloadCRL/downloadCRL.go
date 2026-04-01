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

package downloadCRL

import (
	"fmt"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "download-crl"
	ParameterShortDescription = "Download Certificate Revocation List"
	ParameterLongDescription  = `Download a Certificate Revocation List (CRL) from a URL.

Downloads CRL in PEM format for use with certificate revocation checking.
CRLs are used to verify that certificates have not been revoked by the issuing Certificate Authority.`
	URLFlagDescription        = "URL of the CRL to download"
	OutputPathDescription     = "Path to save downloaded CRL"
	successMessageDownloadCRL = "Successfully downloaded CRL"
	OutputFlagName            = "out"
	URLFlagName               = "url"
)

// ValidateInput - function to validate download-crl inputs
func ValidateInput(cmd *cobra.Command) (string, string, error) {
	crlURL, err := cmd.Flags().GetString(URLFlagName)
	if err != nil {
		return "", "", err
	}

	if crlURL == "" {
		err := fmt.Errorf("Error: required flag '--url' is missing")
		common.SetMandatoryFlagError(cmd, err)
	}

	outputPath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", err
	}

	return crlURL, outputPath, nil
}

// Output - function to print CRL or redirect output to a file
func Output(crlData, outputPath string) error {
	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, crlData)
		if err != nil {
			return err
		}
		fmt.Println(successMessageDownloadCRL)
	} else {
		fmt.Println(crlData)
	}

	return nil
}
