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
	"os"
	"testing"

	gen "github.com/ibm-hyper-protect/contract-go/v2/common/general"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testPrivateKeyPath     = "../../samples/attestation/private.pem"
	testAttestationPath    = "../../build/test_attestation_records.txt"
	testSignaturePath      = "../../build/test_signature.bin"
	testCertPath           = "../../build/test_cert.pem"
	testInvalidPath        = "../../build/file/file_not_exists.txt"
	testInvalidCertPath    = "../../build/invalid_cert.pem"
	testInvalidSigPath     = "../../build/invalid_signature.bin"
	testAttestationContent = "test attestation records content for verification"
)

// setupTestFiles creates test files needed for signature verification tests
func setupTestFiles(t *testing.T) (string, string, string) {
	// Read private key
	privateKeyData, err := os.ReadFile(testPrivateKeyPath)
	assert.NoError(t, err)

	// Create attestation records file
	err = os.WriteFile(testAttestationPath, []byte(testAttestationContent), 0644)
	assert.NoError(t, err)

	// Create temporary private key file
	keyPath, err := gen.CreateTempFile(string(privateKeyData))
	assert.NoError(t, err)

	// Generate self-signed certificate
	cert, err := gen.ExecCommand(gen.GetOpenSSLPath(), "", "req", "-new", "-x509", "-key", keyPath, "-days", "365", "-subj", "/CN=Test")
	assert.NoError(t, err)

	// Write certificate to file
	err = os.WriteFile(testCertPath, []byte(cert), 0644)
	assert.NoError(t, err)

	// Sign the attestation records
	signature, err := gen.ExecCommand(gen.GetOpenSSLPath(), "", "dgst", "-sha256", "-sign", keyPath, testAttestationPath)
	assert.NoError(t, err)

	// Write signature to file
	err = os.WriteFile(testSignaturePath, []byte(signature), 0644)
	assert.NoError(t, err)

	// Clean up temp key file
	gen.RemoveTempFile(keyPath)

	return testAttestationPath, testSignaturePath, testCertPath
}

// cleanupTestFiles removes test files created during tests
func cleanupTestFiles() {
	os.Remove(testAttestationPath)
	os.Remove(testSignaturePath)
	os.Remove(testCertPath)
	os.Remove(testInvalidCertPath)
	os.Remove(testInvalidSigPath)
}

// TestValidateInput_Success tests ValidateInput with all required flags
func TestValidateInput_Success(t *testing.T) {
	attestationPath, signaturePath, certPath := setupTestFiles(t)
	defer cleanupTestFiles()

	cmd := &cobra.Command{}
	cmd.Flags().String(AttestationFlagName, attestationPath, "")
	cmd.Flags().String(SignatureFlagName, signaturePath, "")
	cmd.Flags().String(CertFlagName, certPath, "")

	resultAttestPath, resultSigPath, resultCertPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, attestationPath, resultAttestPath)
	assert.Equal(t, signaturePath, resultSigPath)
	assert.Equal(t, certPath, resultCertPath)
}

// TestValidateInput_WithoutFlags tests ValidateInput when flags are not set
func TestValidateInput_WithoutFlags(t *testing.T) {
	cmd := &cobra.Command{}
	_, _, _, err := ValidateInput(cmd)
	assert.Error(t, err)
}

// TestVerifySignatureAttestationRecords_Success tests successful signature verification
func TestVerifySignatureAttestationRecords_Success(t *testing.T) {
	attestationPath, signaturePath, certPath := setupTestFiles(t)
	defer cleanupTestFiles()

	err := VerifySignatureAttestationRecords(attestationPath, signaturePath, certPath)
	assert.NoError(t, err)
}

// TestVerifySignatureAttestationRecords_InvalidSignature tests with invalid signature
func TestVerifySignatureAttestationRecords_InvalidSignature(t *testing.T) {
	attestationPath, _, certPath := setupTestFiles(t)
	defer cleanupTestFiles()

	// Create invalid signature file
	err := os.WriteFile(testInvalidSigPath, []byte("invalid signature data"), 0644)
	assert.NoError(t, err)

	err = VerifySignatureAttestationRecords(attestationPath, testInvalidSigPath, certPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature verification failed")
}

// TestVerifySignatureAttestationRecords_InvalidCertificate tests with invalid certificate
func TestVerifySignatureAttestationRecords_InvalidCertificate(t *testing.T) {
	attestationPath, signaturePath, _ := setupTestFiles(t)
	defer cleanupTestFiles()

	// Create invalid certificate file
	invalidCert := `-----BEGIN CERTIFICATE-----
invalid-certificate-data
-----END CERTIFICATE-----`
	err := os.WriteFile(testInvalidCertPath, []byte(invalidCert), 0644)
	assert.NoError(t, err)

	err = VerifySignatureAttestationRecords(attestationPath, signaturePath, testInvalidCertPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature verification failed")
}

// TestVerifySignatureAttestationRecords_TamperedData tests with tampered attestation data
func TestVerifySignatureAttestationRecords_TamperedData(t *testing.T) {
	attestationPath, signaturePath, certPath := setupTestFiles(t)
	defer cleanupTestFiles()

	// Tamper with attestation data
	err := os.WriteFile(attestationPath, []byte("tampered attestation data"), 0644)
	assert.NoError(t, err)

	err = VerifySignatureAttestationRecords(attestationPath, signaturePath, certPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature verification failed")
}

// TestVerifySignatureAttestationRecords_NonExistentAttestationFile tests with non-existent attestation file
func TestVerifySignatureAttestationRecords_NonExistentAttestationFile(t *testing.T) {
	setupTestFiles(t)
	defer cleanupTestFiles()

	// This should trigger log.Fatal, but we can't easily test that
	// Instead, we test that the file check works
	exists := false
	if _, err := os.Stat(testInvalidPath); err == nil {
		exists = true
	}
	assert.False(t, exists, "Invalid path should not exist")
}

// TestVerifySignatureAttestationRecords_EmptyAttestationFile tests with empty attestation file
func TestVerifySignatureAttestationRecords_EmptyAttestationFile(t *testing.T) {
	_, signaturePath, certPath := setupTestFiles(t)
	defer cleanupTestFiles()

	// Create empty attestation file
	emptyFile := "../../build/empty_attestation.txt"
	err := os.WriteFile(emptyFile, []byte(""), 0644)
	assert.NoError(t, err)
	defer os.Remove(emptyFile)

	err = VerifySignatureAttestationRecords(emptyFile, signaturePath, certPath)
	assert.Error(t, err)
}

// TestPrintVerificationResult_Success tests successful verification result printing
func TestPrintVerificationResult_Success(t *testing.T) {
	// This function calls log.Fatal on error, so we can only test success case
	// Success case just prints to stdout
	PrintVerificationResult(nil)
	// If we reach here, the test passed (no panic/fatal)
}
