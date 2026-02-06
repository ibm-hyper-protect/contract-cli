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

package decryptAttestation

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testEncAttestPath  = "../../samples/attestation/se-checksums.txt.enc"
	testPrivateKeyPath = "../../samples/attestation/private.pem"
	testPublicKeyPath  = "../../samples/attestation/public.pem"
	testOutputPath     = "../../build/test_decrypt_attestation_output.txt"
	testInvalidPath    = "../../build/file/file_not_exists.txt"
	testInvalidKeyPath = "../../build/key_not_exists.pem"
)

// TestValidateInput_Success tests ValidateInput with all required flags
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testEncAttestPath, "")
	cmd.Flags().String(PrivateKeyFlagName, testPrivateKeyPath, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")

	encAttestPath, privateKeyPath, decryptedAttestPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testEncAttestPath, encAttestPath)
	assert.Equal(t, testPrivateKeyPath, privateKeyPath)
	assert.Equal(t, testOutputPath, decryptedAttestPath)
}

// TestValidateInput_WithoutOutputPath tests ValidateInput without output path
func TestValidateInput_WithoutOutputPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testEncAttestPath, "")
	cmd.Flags().String(PrivateKeyFlagName, testPrivateKeyPath, "")
	cmd.Flags().String(OutputFlagName, "", "")

	encAttestPath, privateKeyPath, decryptedAttestPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testEncAttestPath, encAttestPath)
	assert.Equal(t, testPrivateKeyPath, privateKeyPath)
	assert.Equal(t, "", decryptedAttestPath)
}

// TestPrintDecryptAttestation_ToFile tests writing decrypted attestation to file
func TestPrintDecryptAttestation_ToFile(t *testing.T) {
	testData := "decrypted attestation data"
	os.Remove(testOutputPath)
	err := PrintDecryptAttestation(testData, testOutputPath)
	assert.NoError(t, err)

	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Equal(t, testData, string(content))

	os.Remove(testOutputPath)
}

// TestPrintDecryptAttestation_ToStdout tests printing to stdout (empty path)
func TestPrintDecryptAttestation_ToStdout(t *testing.T) {
	testData := "decrypted attestation data"
	err := PrintDecryptAttestation(testData, "")
	assert.NoError(t, err)
}

// TestPrintDecryptAttestation_InvalidPath tests with invalid output path
func TestPrintDecryptAttestation_InvalidPath(t *testing.T) {
	testData := "decrypted attestation data"
	err := PrintDecryptAttestation(testData, testInvalidPath)
	assert.Error(t, err)
}

// TestValidateInput_FlagErrors tests error handling for flag retrieval
func TestValidateInput_FlagErrors(t *testing.T) {
	cmd := &cobra.Command{}
	_, _, _, err := ValidateInput(cmd)
	assert.Error(t, err)
}
