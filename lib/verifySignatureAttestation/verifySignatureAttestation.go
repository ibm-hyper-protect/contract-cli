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

package verifySignatureAttestation

import (
	"fmt"
	"log"
	"strings"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/v2/attestation"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "verify-signature-attestation"
	ParameterShortDescription = "Verify signature of decrypted attestation records"
	ParameterLongDescription  = `Verify the signature of decrypted attestation records against an IBM attestation certificate.

This ensures the attestation records have not been modified and were signed by IBM.
The attestation records should be in decrypted form (output from decrypt-attestation command).

Typical workflow:
  1. Decrypt attestation: contract-cli decrypt-attestation --in se-checksums.txt.enc --priv private.pem --out se-checksums.txt
  2. Verify signature: contract-cli verify-signature-attestation --attestation se-checksums.txt --signature se-signature.bin --cert attestation-cert.pem`

	AttestationFlagName        = "attestation"
	AttestationFlagDescription = "Path to decrypted attestation records file (se-checksums.txt, use '-' for standard input)"
	SignatureFlagName          = "signature"
	SignatureFlagDescription   = "Path to signature file (se-signature.bin)"
	CertFlagName               = "cert"
	CertFlagDescription        = "Path to IBM attestation certificate file (PEM format)"
	successMessage             = "Signature verification successful - attestation records are valid and signed by IBM"
)

// ValidateInput - function to validate verify-signature-attestation inputs
func ValidateInput(cmd *cobra.Command) (string, string, string, error) {
	attestationPath, err := cmd.Flags().GetString(AttestationFlagName)
	if err != nil {
		return "", "", "", err
	}

	signaturePath, err := cmd.Flags().GetString(SignatureFlagName)
	if err != nil {
		return "", "", "", err
	}

	certPath, err := cmd.Flags().GetString(CertFlagName)
	if err != nil {
		return "", "", "", err
	}

	requiredFlags := map[string]string{
		"--attestation": attestationPath,
		"--signature":   signaturePath,
		"--cert":        certPath,
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
			err := fmt.Errorf("Error: required flags %s are missing.",
				strings.Join(missing, ", "))
			common.SetMandatoryFlagError(cmd, err)
		}
	}

	// Validate stdin input
	common.ValidateStdinInput(cmd, attestationPath)

	return attestationPath, signaturePath, certPath, nil
}

// VerifySignatureAttestationRecords - function to verify signature of attestation records
func VerifySignatureAttestationRecords(attestationPath, signaturePath, certPath string) error {
	var attestationRecords string
	var err error

	// Handle stdin input for attestation records
	if attestationPath == "-" {
		attestationRecords, err = common.ReadDataFromStdin()
		if err != nil {
			return fmt.Errorf("unable to read attestation records from standard input: %w", err)
		}
	} else {
		if !common.CheckFileFolderExists(attestationPath) {
			log.Fatal("The path to attestation records file doesn't exist")
		}
		attestationRecords, err = common.ReadDataFromFile(attestationPath)
		if err != nil {
			return err
		}
	}

	// Read signature file (binary data)
	if !common.CheckFileFolderExists(signaturePath) {
		log.Fatal("The path to signature file doesn't exist")
	}
	signature, err := common.ReadDataFromFile(signaturePath)
	if err != nil {
		return err
	}

	// Read certificate file
	if !common.CheckFileFolderExists(certPath) {
		log.Fatal("The path to certificate file doesn't exist")
	}
	cert, err := common.ReadDataFromFile(certPath)
	if err != nil {
		return err
	}

	// Verify signature using contract-go library
	err = attestation.HpcrVerifySignatureAttestationRecords(attestationRecords, signature, cert)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}

// PrintVerificationResult - print verification result
func PrintVerificationResult(err error) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(successMessage)
}
