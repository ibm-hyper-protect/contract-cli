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
	"encoding/json"
	"fmt"
	"path/filepath"
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

// SealedSecretOutput holds the JSON-serialisable output fields.
type SealedSecretOutput struct {
	SealedSecret    string `json:"sealed_secret"`
	DecryptionKey   string `json:"decryption_key"`
	VerificationKey string `json:"verification_key"`
}

// Output - function to output sealed secret and keys.
// When outputPath is provided the three values are written to separate files
// derived from the base name:
//
//	<base>_SealedValue<ext>   — the sealed secret
//	<base>_DecryptionKey<ext> — the decryption key (PEM, escaped newlines)
//	<base>_VerificationKey<ext> — the verification key (PEM, escaped newlines)
//
// When no outputPath is given, a JSON object containing all three values is
// printed to stdout.
func Output(sealedSecret, decryptionKeyPEM, verificationKeyPEM, outputPath string) error {
	// Format keys with escaped newlines (replace actual newlines with \n)
	decryptionKeyFormatted := formatKeyWithEscapedNewlines(decryptionKeyPEM)
	verificationKeyFormatted := formatKeyWithEscapedNewlines(verificationKeyPEM)

	if outputPath != "" {
		sealedValuePath, decryptionKeyPath, verificationKeyPath := splitOutputPaths(outputPath)

		if err := common.WriteDataToFile(sealedValuePath, sealedSecret); err != nil {
			return fmt.Errorf("failed to write sealed secret to file: %w", err)
		}
		if err := common.WriteDataToFile(decryptionKeyPath, decryptionKeyFormatted); err != nil {
			return fmt.Errorf("failed to write decryption key to file: %w", err)
		}
		if err := common.WriteDataToFile(verificationKeyPath, verificationKeyFormatted); err != nil {
			return fmt.Errorf("failed to write verification key to file: %w", err)
		}

		fmt.Printf("Sealed value written to:      %s\n", sealedValuePath)
		fmt.Printf("Decryption key written to:    %s\n", decryptionKeyPath)
		fmt.Printf("Verification key written to:  %s\n", verificationKeyPath)
	} else {
		// Print JSON to stdout when no output path is specified
		out := SealedSecretOutput{
			SealedSecret:    sealedSecret,
			DecryptionKey:   decryptionKeyFormatted,
			VerificationKey: verificationKeyFormatted,
		}
		jsonBytes, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal output to JSON: %w", err)
		}
		fmt.Println(string(jsonBytes))
	}

	return nil
}

// splitOutputPaths derives the three output file paths from the user-supplied base path.
// The extension (if any) is preserved; suffixes are inserted before it.
// Example: "sealed_secret.yaml" → "sealed_secret_SealedValue.yaml",
//
//	"sealed_secret_DecryptionKey.yaml", "sealed_secret_VerificationKey.yaml"
func splitOutputPaths(outputPath string) (sealedValuePath, decryptionKeyPath, verificationKeyPath string) {
	ext := filepath.Ext(outputPath)
	base := strings.TrimSuffix(outputPath, ext)

	sealedValuePath = base + "_SealedValue" + ext
	decryptionKeyPath = base + "_DecryptionKey" + ext
	verificationKeyPath = base + "_VerificationKey" + ext
	return
}

// formatKeyWithEscapedNewlines converts actual newlines in a key to escaped newlines (\n)
// This matches the format: cat key.pem | tr '\n' '\\' | sed s/\\\\/\\\\n/g
func formatKeyWithEscapedNewlines(key string) string {
	// Replace actual newlines with the literal string "\n"
	return strings.ReplaceAll(key, "\n", "\\n")
}
