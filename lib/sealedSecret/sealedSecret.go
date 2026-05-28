// Copyright (c) 2026 IBM Corp.
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

package sealedSecret

import (
	"fmt"
	"strings"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/v2/secrets"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "sealed-secret"
	ParameterShortDescription = "Generate sealed secret for CCCO"
	ParameterLongDescription  = `Generate a sealed secret for IBM Confidential Computing Containers for Red Hat Openshift Container Platform.

Creates a sealed secret that can be used in workload or environment sections
of your contract. The secret can be provided as a string or read from a file.`

	InputFlagName        = "in"
	InputFlagDescription = "Secret for sealing (provide as string or file path, use '-' for standard input)"

	TypeFlagName        = "type"
	TypeFlagDescription = "Type of secret: 'env' for env section of contract or 'workload' for workload section of contract"

	OutputFlagName        = "out"
	OutputFlagDescription = "Path to save sealed secret output (optional, prints to stdout if not specified)"

	EncryptionKeyFlagName        = "encryptionkey"
	EncryptionKeyFlagDescription = "Path to RSA private key for encryption (optional, generates new key if not provided)"

	SigningKeyFlagName        = "signingkey"
	SigningKeyFlagDescription = "Path to RSA private key for signing (optional, generates new key if not provided)"
)

// ValidateInput - function to validate inputs of sealed-secret
func ValidateInput(cmd *cobra.Command) (string, string, string, string, string, error) {
	inputData, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	if inputData == "" {
		err := fmt.Errorf("Error: required flag '--in' is missing")
		common.SetMandatoryFlagError(cmd, err)
	}

	// Validate stdin input conflicts
	common.ValidateStdinInput(cmd, inputData)

	secretType, err := cmd.Flags().GetString(TypeFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	if secretType == "" {
		err := fmt.Errorf("Error: required flag '--type' is missing")
		common.SetMandatoryFlagError(cmd, err)
	}

	// Validate type value
	if secretType != "env" && secretType != "workload" {
		err := fmt.Errorf("Error: invalid value for '--type'. Must be 'env' or 'workload'")
		common.SetMandatoryFlagError(cmd, err)
	}

	outputPath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	encryptionKeyPath, err := cmd.Flags().GetString(EncryptionKeyFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	signingKeyPath, err := cmd.Flags().GetString(SigningKeyFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	return inputData, secretType, outputPath, encryptionKeyPath, signingKeyPath, nil
}

// GenerateSealedSecret - function to generate sealed secret using contract-go SealSecret API
// Returns: sealedSecret, decryptionKeyPEM, verificationKeyPEM, inputSecretSha, encryptedSecretSha, error
func GenerateSealedSecret(inputDataPath, secretType, encryptionKeyPath, signingKeyPath string) (string, string, string, string, string, error) {
	var inputData string
	var err error

	// Handle stdin input
	if inputDataPath == "-" {
		inputData, err = common.ReadDataFromStdin()
		if err != nil {
			return "", "", "", "", "", fmt.Errorf("unable to read input from standard input: %w", err)
		}
	} else {
		// Check if input is a file path
		if common.CheckFileFolderExists(inputDataPath) {
			inputData, err = common.ReadDataFromFile(inputDataPath)
			if err != nil {
				return "", "", "", "", "", fmt.Errorf("unable to read input from file: %w", err)
			}
		} else {
			// Treat as direct string input
			inputData = inputDataPath
		}
	}

	// Read encryption key content if path is provided
	var encryptionKeyContent string
	if encryptionKeyPath != "" {
		encKeyStr, err := common.GetDataFromFile(encryptionKeyPath)
		if err != nil {
			return "", "", "", "", "", fmt.Errorf("failed to read encryption key from file: %w", err)
		}
		encryptionKeyContent = encKeyStr
	}

	// Read signing key content if path is provided
	var signingKeyContent string
	if signingKeyPath != "" {
		signKeyStr, err := common.GetDataFromFile(signingKeyPath)
		if err != nil {
			return "", "", "", "", "", fmt.Errorf("failed to read signing key from file: %w", err)
		}
		signingKeyContent = signKeyStr
	}

	// Call the SealSecret API from contract-go secrets package
	// The function now accepts PEM string content directly and returns multiple values
	sealedSecret, decryptionKey, verificationKey, inputSecretSha, encryptedSecretSha, err := secrets.HpccSealedSecret(
		inputData,
		secretType,
		encryptionKeyContent,
		signingKeyContent,
	)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to seal secret: %w", err)
	}

	return sealedSecret, decryptionKey, verificationKey, inputSecretSha, encryptedSecretSha, nil
}

// Output - function to output sealed secret and keys
func Output(sealedSecret, decryptionKeyPEM, verificationKeyPEM, outputPath string) error {
	// Format keys with escaped newlines (replace actual newlines with \n)
	decryptionKeyFormatted := formatKeyWithEscapedNewlines(decryptionKeyPEM)
	verificationKeyFormatted := formatKeyWithEscapedNewlines(verificationKeyPEM)

	// Format the output as shell variable assignments
	output := fmt.Sprintf("Sealed Secret:\n%s\n\nSECRET_DECRYPTION_KEY=%s\n\nSECRET_VERIFICATION_KEY=%s",
		sealedSecret,
		decryptionKeyFormatted,
		verificationKeyFormatted)

	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, output)
		if err != nil {
			return fmt.Errorf("failed to write sealed secret to file: %w", err)
		}
		fmt.Printf("Sealed secret and keys written to: %s\n", outputPath)
	} else {
		// Print to stdout if no output path specified
		fmt.Println(output)
	}

	return nil
}

// formatKeyWithEscapedNewlines converts actual newlines in a key to escaped newlines (\n)
// This matches the format: cat key.pem | tr '\n' '\\' | sed s/\\\\/\\\\n/g
func formatKeyWithEscapedNewlines(key string) string {
	// Replace actual newlines with the literal string "\n"
	return strings.ReplaceAll(key, "\n", "\\n")
}
