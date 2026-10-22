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

package common

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	sampleFile = "../samples/attestation/se-checksums.txt.enc"

	simpleSampleTextPath = "../samples/simple_file.txt"
	simpleSampleText     = "Testing"

	simpleSampleWritePath = "../build/simple_file.txt"

	samplePrivateKeyPath = "../samples/sign/private.pem"

	sampleCertPath = "../samples/contract-expiry/personal_ca.crt"
)

// Testcase to check if CheckFileFolderExists() is able to check if a folder exists
func TestCheckFileFolderExists(t *testing.T) {
	result := CheckFileFolderExists(sampleFile)

	assert.True(t, result)
}

// Testcase to check if ReadDataFromFile() is able to read data from file
func TestReadDataFromFile(t *testing.T) {
	content, err := ReadDataFromFile(simpleSampleTextPath)
	if err != nil {
		t.Errorf("failed to read text from file - %v", err)
	}

	assert.Equal(t, content, simpleSampleText)
}

// Testcase to check if WriteDataToFile() is able to write data to file
func TestWriteDataToFile(t *testing.T) {
	err := WriteDataToFile(simpleSampleWritePath, simpleSampleText)
	if err != nil {
		t.Errorf("failed to write data to file - %v", err)
	}
}

// Testcase to verify WriteDataToFile() resets permissions on a pre-existing
// permissive file. A file pre-created with mode 0666 must be replaced by a
// 0600 file after the write, so that other OS users cannot read sensitive content.
func TestWriteDataToFile_PreExistingPermissiveFile(t *testing.T) {
	path := "../build/permissive_overwrite_test.txt"

	// Simulate a pre-existing file with mode 0666.
	if err := os.WriteFile(path, []byte("existing"), 0666); err != nil {
		t.Fatalf("failed to pre-create file: %v", err)
	}
	t.Cleanup(func() { os.Remove(path) })

	if err := WriteDataToFile(path, simpleSampleText); err != nil {
		t.Fatalf("WriteDataToFile failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat written file: %v", err)
	}

	got := info.Mode().Perm()
	if got != 0600 {
		t.Errorf("expected mode 0600 after overwrite, got %04o", got)
	}
}

// Testcase to check if ExecCommand() is able to execute command
func TestExecCommand(t *testing.T) {
	_, err := ExecCommand("openssl", "", "version")
	if err != nil {
		t.Errorf("failed to execute command - %v", err)
	}
}

// Testcase to check if OpensslCheck() is able to check if openssl is present
func TestOpensslCheck(t *testing.T) {
	err := OpensslCheck()
	if err != nil {
		t.Errorf("openssl check failed - %v", err)
	}
}

// Testcase to check if GetPrivateKey() is able to get key from file
func TestGetPrivateKeyNoKey(t *testing.T) {
	result, err := GetPrivateKey(samplePrivateKeyPath)
	if err != nil {
		t.Errorf("failed to get private key - %v", err)
	}

	assert.NotEmpty(t, result)
}

// Testcase to check if GetPrivateKey() is able generate private key
func TestGetPrivateKey(t *testing.T) {
	result, err := GetPrivateKey("")
	if err != nil {
		t.Errorf("failed to get private key - %v", err)
	}

	assert.NotEmpty(t, result)
}

// Testcase to check if generatePrivateKey() is able to generate private key
func TestGeneratePrivateKey(t *testing.T) {
	result, err := generatePrivateKey()
	if err != nil {
		t.Errorf("failed to generate private key - %v", err)
	}

	assert.NotEmpty(t, result)
}

// Testcase to check if GetDataFromFile() is able to get data form file
func TestGetDataFromFileWithData(t *testing.T) {
	result, err := GetDataFromFile(sampleCertPath)
	if err != nil {
		t.Errorf("failed to get data from file - %v", err)
	}

	assert.NotEmpty(t, result)
}

// Testcase to check if GetDataFromFile() is throwing error for invalid file path
func TestGetDataFromFileWithoutData(t *testing.T) {
	result, err := GetDataFromFile("")
	if err != nil {
		t.Errorf("failed to get data from file - %v", err)
	}

	assert.Empty(t, result)
}

// Testcase to verify that SetCustomHelpTemplate() prints the help message for 'contract-cli --help' in the correct format.
func TestSetCustomHelpTemplate(t *testing.T) {
	testCmd := &cobra.Command{
		Use: "testcmd",
	}

	testCmd.Flags().String("required1", "", "Required flag 1")
	testCmd.Flags().String("optional1", "", "Optional flag 1")

	requiredFlags := map[string]bool{
		"required1": true,
	}

	var buf bytes.Buffer
	testCmd.SetOut(&buf)
	testCmd.SetErr(&buf)

	testCmd.SetIn(strings.NewReader(""))
	SetCustomHelpTemplate(testCmd, requiredFlags)

	testCmd.SetArgs([]string{"--help"})
	_ = testCmd.Execute()

	output := buf.String()
	if output == "" {
		t.Fatalf("contract cli --help is empty")
	}

	if !strings.Contains(output, "Usage:") {
		t.Fatalf("Expected 'Usage:' in help output, got:\n%s", output)
	}

	if !strings.Contains(output, "Mandatory Flags:") {
		t.Errorf("expected Mandatory Flags section, got: %s", output)
	}

	if !strings.Contains(output, "Optional Flags:") {
		t.Errorf("expected Optional Flags section, got: %s", output)
	}

	mandatoryIdx := strings.Index(output, "Mandatory Flags:")
	requiredIdx := strings.Index(output, "--required1")

	if !(mandatoryIdx < requiredIdx && requiredIdx != -1) {
		t.Errorf("required flag not found under Mandatory Flags")
	}

	optionalIdx := strings.Index(output, "--optional1")
	optionalSectionIdx := strings.Index(output, "Optional Flags:")

	if !(optionalSectionIdx < optionalIdx && optionalIdx != -1) {
		t.Errorf("optional flag not found under Optional Flags")
	}

	if !strings.Contains(output, "Required flag 1") {
		t.Errorf("missing usage text for required1 flag")
	}

	if !strings.Contains(output, "Optional flag 1") {
		t.Errorf("missing usage text for optional1 flag")
	}
}
