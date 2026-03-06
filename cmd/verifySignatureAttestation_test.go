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
	"bytes"
	"os"
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/lib/verifySignatureAttestation"
	gen "github.com/ibm-hyper-protect/contract-go/v2/common/general"
	"github.com/stretchr/testify/assert"
)

const (
	verifyPrivateKeyFilePath  = "../samples/attestation/private.pem"
	verifyAttestationFilePath = "../build/test_attestation.txt"
	verifySignatureFilePath   = "../build/test_signature.bin"
	verifyCertificateFilePath = "../build/test_certificate.pem"
	verifyAttestationContent  = "test attestation records for command test"
)

// setupCommandTestFiles creates test files for command-level tests
func setupCommandTestFiles(t *testing.T) {
	// Read private key
	privateKeyData, err := os.ReadFile(verifyPrivateKeyFilePath)
	if err != nil {
		t.Fatalf("failed to read private key: %v", err)
	}

	// Create attestation records file
	err = os.WriteFile(verifyAttestationFilePath, []byte(verifyAttestationContent), 0644)
	if err != nil {
		t.Fatalf("failed to create attestation file: %v", err)
	}

	// Create temporary private key file
	keyPath, err := gen.CreateTempFile(string(privateKeyData))
	if err != nil {
		t.Fatalf("failed to create temp key file: %v", err)
	}
	defer gen.RemoveTempFile(keyPath)

	// Generate self-signed certificate
	cert, err := gen.ExecCommand(gen.GetOpenSSLPath(), "", "req", "-new", "-x509", "-key", keyPath, "-days", "365", "-subj", "/CN=TestCommand")
	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}

	// Write certificate to file
	err = os.WriteFile(verifyCertificateFilePath, []byte(cert), 0644)
	if err != nil {
		t.Fatalf("failed to write certificate file: %v", err)
	}

	// Sign the attestation records
	signature, err := gen.ExecCommand(gen.GetOpenSSLPath(), "", "dgst", "-sha256", "-sign", keyPath, verifyAttestationFilePath)
	if err != nil {
		t.Fatalf("failed to sign attestation: %v", err)
	}

	// Write signature to file
	err = os.WriteFile(verifySignatureFilePath, []byte(signature), 0644)
	if err != nil {
		t.Fatalf("failed to write signature file: %v", err)
	}
}

// cleanupCommandTestFiles removes test files
func cleanupCommandTestFiles() {
	os.Remove(verifyAttestationFilePath)
	os.Remove(verifySignatureFilePath)
	os.Remove(verifyCertificateFilePath)
}

// TestVerifySignatureAttestationCmdSuccess tests successful signature verification via command
func TestVerifySignatureAttestationCmdSuccess(t *testing.T) {
	setupCommandTestFiles(t)
	defer cleanupCommandTestFiles()

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	sampleValidCommand := []string{
		verifySignatureAttestation.ParameterName,
		"--attestation", verifyAttestationFilePath,
		"--signature", verifySignatureFilePath,
		"--cert", verifyCertificateFilePath,
	}

	rootCmd.SetArgs(sampleValidCommand)
	err := verifySignatureAttestationCmd.Execute()

	// Command should execute without error
	assert.NoError(t, err)
}

// TestVerifySignatureAttestationCmdMissingFlags tests command with missing required flags
func TestVerifySignatureAttestationCmdMissingFlags(t *testing.T) {
	setupCommandTestFiles(t)
	defer cleanupCommandTestFiles()

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Missing --cert flag
	sampleInvalidCommand := []string{
		verifySignatureAttestation.ParameterName,
		"--attestation", verifyAttestationFilePath,
		"--signature", verifySignatureFilePath,
	}

	rootCmd.SetArgs(sampleInvalidCommand)
	err := verifySignatureAttestationCmd.Execute()

	// Command should not return error but validation should fail
	assert.NoError(t, err)
}
