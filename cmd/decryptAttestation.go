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
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-cli/lib/decryptAttestation"
	"github.com/spf13/cobra"
)

// decryptAttestationCmd represents the decrypt-attestation command
var decryptAttestationCmd = &cobra.Command{
	Use:   decryptAttestation.ParameterName,
	Short: decryptAttestation.ParameterShortDescription,
	Long:  decryptAttestation.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		encAttestPath, privateKeyPath, decryptedAttestPath, signaturePath, certPath, err := decryptAttestation.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		decryptedAttestationRecords, err := decryptAttestation.DecryptAttestationRecords(encAttestPath, privateKeyPath)
		if err != nil {
			log.Fatal(err)
		}

		// If signature and cert are provided, verify the signature
		if signaturePath != "" && certPath != "" {
			err = decryptAttestation.VerifySignatureAttestationRecords(decryptedAttestationRecords, signaturePath, certPath)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Signature verification successful - attestation records are valid and signed by IBM")
			fmt.Println()
		}

		err = decryptAttestation.PrintDecryptAttestation(decryptedAttestationRecords, decryptedAttestPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(decryptAttestationCmd)

	requiredFlags := map[string]bool{
		"in":   true,
		"priv": true,
	}
	decryptAttestationCmd.PersistentFlags().String(decryptAttestation.InputFlagName, "", decryptAttestation.DecryptAttestFileInDescription)
	decryptAttestationCmd.PersistentFlags().String(decryptAttestation.PrivateKeyFlagName, "", decryptAttestation.PrivateKeyFlagDescription)
	decryptAttestationCmd.PersistentFlags().String(decryptAttestation.OutputFlagName, "", decryptAttestation.DecryptAttestFlagDescription)
	decryptAttestationCmd.PersistentFlags().String(decryptAttestation.SignatureFlagName, "", decryptAttestation.SignatureFlagDescription)
	decryptAttestationCmd.PersistentFlags().String(decryptAttestation.AttestationCertFlagName, "", decryptAttestation.AttestationCertFlagDescription)
	common.SetCustomHelpTemplate(decryptAttestationCmd, requiredFlags)
	common.SetCustomErrorTemplate(decryptAttestationCmd)
}
