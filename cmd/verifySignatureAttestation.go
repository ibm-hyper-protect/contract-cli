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
	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-cli/lib/verifySignatureAttestation"
	"github.com/spf13/cobra"
)

// verifySignatureAttestationCmd represents the verify-signature-attestation command
var verifySignatureAttestationCmd = &cobra.Command{
	Use:   verifySignatureAttestation.ParameterName,
	Short: verifySignatureAttestation.ParameterShortDescription,
	Long:  verifySignatureAttestation.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		attestationPath, signaturePath, certPath, err := verifySignatureAttestation.ValidateInput(cmd)
		if err != nil {
			return
		}

		err = verifySignatureAttestation.VerifySignatureAttestationRecords(attestationPath, signaturePath, certPath)
		verifySignatureAttestation.PrintVerificationResult(err)
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(verifySignatureAttestationCmd)

	requiredFlags := map[string]bool{
		"attestation": true,
		"signature":   true,
		"cert":        true,
	}
	verifySignatureAttestationCmd.PersistentFlags().String(verifySignatureAttestation.AttestationFlagName, "", verifySignatureAttestation.AttestationFlagDescription)
	verifySignatureAttestationCmd.PersistentFlags().String(verifySignatureAttestation.SignatureFlagName, "", verifySignatureAttestation.SignatureFlagDescription)
	verifySignatureAttestationCmd.PersistentFlags().String(verifySignatureAttestation.CertFlagName, "", verifySignatureAttestation.CertFlagDescription)
	common.SetCustomHelpTemplate(verifySignatureAttestationCmd, requiredFlags)
	common.SetCustomErrorTemplate(verifySignatureAttestationCmd)
}
