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

package cmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/ibm-hyper-protect/contract-cli/lib/decryptString"
)

const (
	testDecryptEncryptedInput = "../samples/decrypt/encrypt.txt"
	testDecryptPrivateKey     = "../samples/decrypt/private.key"
	testDecryptOutputPath     = "../build/test_cmd_decrypt_output.txt"
)

// getDecryptStringCmd returns a fresh instance of the decrypt command for testing
func getDecryptStringCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   decryptString.ParameterName,
		Short: decryptString.ParameterShortDescription,
		Long:  decryptString.ParameterLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			inputData, privateKeyPath, password, outputPath, err := decryptString.ValidateInput(cmd)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			decryptedText, err := decryptString.Process(inputData, privateKeyPath, password)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			err = decryptString.Output(outputPath, decryptedText)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
		},
	}

	cmd.Flags().String(decryptString.InputFlagName, "", decryptString.InputFlagDescription)
	cmd.Flags().String(decryptString.PrivateKeyFlagName, "", decryptString.PrivateKeyFlagDescription)
	cmd.Flags().String(decryptString.PasswordFlagName, "", decryptString.PasswordFlagDescription)
	cmd.Flags().String(decryptString.OutputFlagName, "", decryptString.OutputFlagDescription)

	return cmd
}

// TestDecryptStringCmd_Success - decrypt with valid input and private key succeeds
func TestDecryptStringCmd_Success(t *testing.T) {
	defer os.Remove(testDecryptOutputPath)

	cmd := getDecryptStringCmd()
	cmd.SetArgs([]string{
		"--" + decryptString.InputFlagName, testDecryptEncryptedInput,
		"--" + decryptString.PrivateKeyFlagName, testDecryptPrivateKey,
		"--" + decryptString.OutputFlagName, testDecryptOutputPath,
	})

	err := cmd.Execute()
	assert.NoError(t, err)

	// Verify output file was created
	_, statErr := os.Stat(testDecryptOutputPath)
	assert.NoError(t, statErr)

	// Verify output is not empty
	content, readErr := os.ReadFile(testDecryptOutputPath)
	assert.NoError(t, readErr)
	assert.NotEmpty(t, content)
}

// TestDecryptStringCmd_WithoutOutputPath - decrypt to stdout succeeds
func TestDecryptStringCmd_WithoutOutputPath(t *testing.T) {
	cmd := getDecryptStringCmd()
	cmd.SetArgs([]string{
		"--" + decryptString.InputFlagName, testDecryptEncryptedInput,
		"--" + decryptString.PrivateKeyFlagName, testDecryptPrivateKey,
	})

	err := cmd.Execute()
	assert.NoError(t, err)
}
