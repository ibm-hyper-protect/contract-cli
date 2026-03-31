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
	testEncAttestPath   = "../../samples/attestation/se-checksums.txt.enc"
	testPrivateKeyPath  = "../../samples/attestation/private.pem"
	testPublicKeyPath   = "../../samples/attestation/public.pem"
	testOutputPath      = "../../build/test_decrypt_attestation_output.txt"
	testInvalidPath     = "../../build/file/file_not_exists.txt"
	testInvalidKeyPath  = "../../build/key_not_exists.pem"
	testCorruptedAttest = "../../build/corrupted_attest.enc"
	testEmptyAttest     = "../../build/empty_attest.enc"
	testInvalidKeyFile  = "../../build/invalid_key.pem"
)

// TestValidateInput_Success tests ValidateInput with all required flags
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testEncAttestPath, "")
	cmd.Flags().String(PrivateKeyFlagName, testPrivateKeyPath, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")
	cmd.Flags().String(SignatureFlagName, "", "")
	cmd.Flags().String(AttestationCertFlagName, "", "")
	cmd.Flags().String(PasswordFlagName, "", "")

	encAttestPath, privateKeyPath, decryptedAttestPath, signaturePath, certPath, password, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testEncAttestPath, encAttestPath)
	assert.Equal(t, testPrivateKeyPath, privateKeyPath)
	assert.Equal(t, testOutputPath, decryptedAttestPath)
	assert.Equal(t, "", signaturePath)
	assert.Equal(t, "", certPath)
	assert.Equal(t, "", password)
}

// TestValidateInput_WithoutOutputPath tests ValidateInput without output path
func TestValidateInput_WithoutOutputPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testEncAttestPath, "")
	cmd.Flags().String(PrivateKeyFlagName, testPrivateKeyPath, "")
	cmd.Flags().String(OutputFlagName, "", "")
	cmd.Flags().String(SignatureFlagName, "", "")
	cmd.Flags().String(AttestationCertFlagName, "", "")
	cmd.Flags().String(PasswordFlagName, "", "")

	encAttestPath, privateKeyPath, decryptedAttestPath, signaturePath, certPath, password, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testEncAttestPath, encAttestPath)
	assert.Equal(t, testPrivateKeyPath, privateKeyPath)
	assert.Equal(t, "", decryptedAttestPath)
	assert.Equal(t, "", signaturePath)
	assert.Equal(t, "", certPath)
	assert.Equal(t, "", password)
}

// TestValidateInput_WithBothSignatureAndCert tests ValidateInput with both signature and cert flags
func TestValidateInput_WithBothSignatureAndCert(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testEncAttestPath, "")
	cmd.Flags().String(PrivateKeyFlagName, testPrivateKeyPath, "")
	cmd.Flags().String(OutputFlagName, "", "")
	cmd.Flags().String(SignatureFlagName, "signature.bin", "")
	cmd.Flags().String(AttestationCertFlagName, "cert.pem", "")
	cmd.Flags().String(PasswordFlagName, "", "")

	encAttestPath, privateKeyPath, decryptedAttestPath, signaturePath, certPath, password, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testEncAttestPath, encAttestPath)
	assert.Equal(t, testPrivateKeyPath, privateKeyPath)
	assert.Equal(t, "", decryptedAttestPath)
	assert.Equal(t, "signature.bin", signaturePath)
	assert.Equal(t, "cert.pem", certPath)
	assert.Equal(t, "", password)
}

// TestDecryptAttestationRecords_Success tests successful decryption
func TestDecryptAttestationRecords_Success(t *testing.T) {
	result, err := DecryptAttestationRecords(testEncAttestPath, testPrivateKeyPath, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestDecryptAttestationRecords_CorruptedEncryptedData tests with corrupted encrypted data
func TestDecryptAttestationRecords_CorruptedEncryptedData(t *testing.T) {
	// Create a temporary file with corrupted data that looks like encrypted format
	corruptedFile := "../../build/corrupted_attest.enc"
	err := os.WriteFile(corruptedFile, []byte("hyper-protect-basic.corrupted.data"), 0644)
	assert.NoError(t, err)
	defer os.Remove(corruptedFile)

	result, err := DecryptAttestationRecords(corruptedFile, testPrivateKeyPath, "")

	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestDecryptAttestationRecords_WrongPrivateKey tests with wrong private key
func TestDecryptAttestationRecords_WrongPrivateKey(t *testing.T) {
	// Using public key instead of private key should fail
	result, err := DecryptAttestationRecords(testEncAttestPath, testPublicKeyPath, "")

	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestDecryptAttestationRecords_EmptyEncryptedFile tests with empty encrypted file
func TestDecryptAttestationRecords_EmptyEncryptedFile(t *testing.T) {
	// Create a temporary empty file
	emptyFile := "../../build/empty_attest.enc"
	err := os.WriteFile(emptyFile, []byte(""), 0644)
	assert.NoError(t, err)
	defer os.Remove(emptyFile)

	result, err := DecryptAttestationRecords(emptyFile, testPrivateKeyPath, "")

	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestDecryptAttestationRecords_InvalidKeyFormat tests with invalid key format
func TestDecryptAttestationRecords_InvalidKeyFormat(t *testing.T) {
	// Create a temporary file with invalid key format
	invalidKeyFile := "../../build/invalid_key.pem"
	err := os.WriteFile(invalidKeyFile, []byte("not a valid key"), 0644)
	assert.NoError(t, err)
	defer os.Remove(invalidKeyFile)

	result, err := DecryptAttestationRecords(testEncAttestPath, invalidKeyFile, "")

	assert.Error(t, err)
	assert.Equal(t, "", result)
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

// TestPrintDecryptAttestation_ToStdout tests printing to stdout
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

// TestValidateInput_WithPassword tests ValidateInput with password flag
func TestValidateInput_WithPassword(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testEncAttestPath, "")
	cmd.Flags().String(PrivateKeyFlagName, testPrivateKeyPath, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")
	cmd.Flags().String(SignatureFlagName, "", "")
	cmd.Flags().String(AttestationCertFlagName, "", "")
	cmd.Flags().String(PasswordFlagName, "testPassword123", "")

	encAttestPath, privateKeyPath, decryptedAttestPath, signaturePath, certPath, password, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testEncAttestPath, encAttestPath)
	assert.Equal(t, testPrivateKeyPath, privateKeyPath)
	assert.Equal(t, testOutputPath, decryptedAttestPath)
	assert.Equal(t, "", signaturePath)
	assert.Equal(t, "", certPath)
	assert.Equal(t, "testPassword123", password)
}

// TestDecryptAttestationRecords_WithEmptyPassword tests decryption with empty password (unencrypted key)
func TestDecryptAttestationRecords_WithEmptyPassword(t *testing.T) {
	result, err := DecryptAttestationRecords(testEncAttestPath, testPrivateKeyPath, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	// Verify result contains attestation data (version, checksums, etc.)
	assert.Contains(t, result, "Machine Type")
}

// TestDecryptAttestationRecords_WithWrongPassword tests decryption with wrong password
// Note: This test expects an unencrypted key, so providing a password should still work
// For a truly encrypted key with wrong password, contract-go would return an error
func TestDecryptAttestationRecords_PasswordHandling(t *testing.T) {
	// Using unencrypted key with password parameter should work (password ignored)
	result, err := DecryptAttestationRecords(testEncAttestPath, testPrivateKeyPath, "anyPassword")

	// Should succeed because the key is not encrypted
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}
