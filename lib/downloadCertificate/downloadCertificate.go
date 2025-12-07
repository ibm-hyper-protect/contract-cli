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

package downloadCertificate

import (
	"fmt"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "download-certificate"
	ParameterShortDescription = "Download encryption certificates"
	ParameterLongDescription  = `Download encryption certificates from the IBM Hyper Protect Repository.

Retrieves the latest or specific versions of HPCR encryption certificates required
for contract encryption and workload deployment.`
	EncryptionCertVersionDescription  = "Specific certificate versions to download (comma-separated, e.g., 1.0.21,1.0.22)"
	OutputPathDescription             = "Path to save downloaded encryption certificates"
	successMessageDownloadCertificate = "Successfully stored certificates"
	OutputFlagName                    = "out"
	FormatFlag                        = "format"
	VersionFlag                       = "version"
	JsonFormat                        = "json"
	DataFormatFlag                    = "Output format for data (json, yaml, or text)"
)

// ValidateInput - function to validate download-certificate inputs
func ValidateInput(cmd *cobra.Command) (string, string, error) {
	formatType, err := cmd.Flags().GetString(FormatFlag)
	if err != nil {
		return "", "", err
	}

	certificatePath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", err
	}

	return formatType, certificatePath, nil
}

// Output - function to print encryption certificates or redirect output to a file
func Output(certificates, certificatePath string) error {
	if certificatePath != "" {
		err := common.WriteDataToFile(certificatePath, certificates)
		if err != nil {
			return err
		}
		fmt.Println(successMessageDownloadCertificate)
	} else {
		fmt.Println(certificates)
	}

	return nil
}
