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

package cmd

import (
	"log"

	"github.com/ibm-hyper-protect/contract-cli/lib/downloadCertificate"
	"github.com/ibm-hyper-protect/contract-go/v2/certificate"
	"github.com/spf13/cobra"
)

var (
	// downloadCertificatesCmd represents the download-certificate command
	downloadCertificatesCmd = &cobra.Command{
		Use:   downloadCertificate.ParameterName,
		Short: downloadCertificate.ParameterShortDescription,
		Long:  downloadCertificate.ParameterLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			formatType, certificatePath, err := downloadCertificate.ValidateInput(cmd)
			if err != nil {
				log.Fatal(err)
			}

			certificates, err := certificate.HpcrDownloadEncryptionCertificates(versions, formatType, "")
			if err != nil {
				log.Fatal(err)
			}

			err = downloadCertificate.Output(certificates, certificatePath)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	versions []string
)

// init - cobra init function
func init() {
	rootCmd.AddCommand(downloadCertificatesCmd)

	downloadCertificatesCmd.PersistentFlags().StringSliceVarP(&versions, downloadCertificate.VersionFlag, "", []string{}, downloadCertificate.EncryptionCertVersionDescription)
	downloadCertificatesCmd.PersistentFlags().String(downloadCertificate.FormatFlag, downloadCertificate.JsonFormat, downloadCertificate.DataFormatFlag)
	downloadCertificatesCmd.PersistentFlags().String(downloadCertificate.OutputFlagName, "", downloadCertificate.OutputPathDescription)
}
