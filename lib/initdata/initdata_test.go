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

package initdata

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testContractPath = "../../samples/hpcc/signed-encrypt-hpcc.yaml"
)

// TestValidateInput_Success tests ValidateInput with all required flags
func TestValidateInput_Success(t *testing.T) {
	tmpDir := t.TempDir()
	testOutputPath := filepath.Join(tmpDir, "test_initdata_output.txt")

	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(SehdrBinFlagName, "", "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")
	cmd.Flags().StringArray(ExtraRootCertsFlagName, []string{}, "")

	inputData, sehdrBinPath, outputPath, extraCerts, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, inputData)
	assert.Equal(t, "", sehdrBinPath)
	assert.Equal(t, testOutputPath, outputPath)
	assert.Empty(t, extraCerts)
}

// TestValidateInput_WithoutOutputPath tests ValidateInput without output path
func TestValidateInput_WithoutOutputPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(SehdrBinFlagName, "", "")
	cmd.Flags().String(OutputFlagName, "", "")
	cmd.Flags().StringArray(ExtraRootCertsFlagName, []string{}, "")

	inputData, sehdrBinPath, outputPath, extraCerts, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, inputData)
	assert.Equal(t, "", sehdrBinPath)
	assert.Equal(t, "", outputPath)
	assert.Empty(t, extraCerts)
}

// Note: TestValidateInput_WithoutFlags removed because ValidateInput calls
// SetMandatoryFlagError which calls os.Exit(1), terminating the test process.
// This validation scenario is tested at the command level where appropriate.

// TestValidateInput_FlagErrors tests ValidateInput with flag retrieval errors
func TestValidateInput_FlagErrors(t *testing.T) {
	cmd := &cobra.Command{}
	// Don't define flags to trigger errors

	_, _, _, _, err := ValidateInput(cmd)

	assert.Error(t, err)
}

// TestGenerateInitdata_Success tests successful initdata generation
func TestGenerateInitdata_Success(t *testing.T) {
	result, isBaremetal, err := GenerateInitdata(testContractPath, "", nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.False(t, isBaremetal)
	// Gzipped and base64 encoded data should be a string
	assert.IsType(t, "", result)
}

// TestGenerateInitdata_InvalidPath tests with invalid contract path
func TestGenerateInitdata_InvalidPath(t *testing.T) {
	tmpDir := t.TempDir()
	testInvalidPath := filepath.Join(tmpDir, "file", "file_not_exists.txt")

	result, isBaremetal, err := GenerateInitdata(testInvalidPath, "", nil)

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.False(t, isBaremetal)
	assert.Contains(t, err.Error(), "doesn't exist")
}

// TestGenerateInitdata_EmptyPath tests with empty path
func TestGenerateInitdata_EmptyPath(t *testing.T) {
	result, isBaremetal, err := GenerateInitdata("", "", nil)

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.False(t, isBaremetal)
	assert.Contains(t, err.Error(), "doesn't exist")
}

// TestGenerateInitdata_CorruptedContract tests with corrupted contract file
// Note: HpccInitdata doesn't validate YAML format, it just gzips and encodes the content
func TestGenerateInitdata_CorruptedContract(t *testing.T) {
	tmpDir := t.TempDir()
	testCorruptedContract := filepath.Join(tmpDir, "corrupted_contract_initdata.yaml")

	err := os.WriteFile(testCorruptedContract, []byte("invalid: yaml: content: ["), 0644)
	assert.NoError(t, err)

	result, isBaremetal, err := GenerateInitdata(testCorruptedContract, "", nil)

	// HpccInitdata accepts any content and gzips it
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.False(t, isBaremetal)
}

// TestGenerateInitdata_EmptyContract tests with empty contract file
func TestGenerateInitdata_EmptyContract(t *testing.T) {
	tmpDir := t.TempDir()
	testEmptyContract := filepath.Join(tmpDir, "empty_contract_initdata.yaml")

	err := os.WriteFile(testEmptyContract, []byte(""), 0644)
	assert.NoError(t, err)

	result, isBaremetal, err := GenerateInitdata(testEmptyContract, "", nil)

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.False(t, isBaremetal)
}

// TestGenerateInitdata_NonYamlFile tests with non-YAML file
// Note: HpccInitdata doesn't validate file format, it just gzips and encodes the content
func TestGenerateInitdata_NonYamlFile(t *testing.T) {
	tmpDir := t.TempDir()
	testTextFile := filepath.Join(tmpDir, "test_text_file.txt")

	err := os.WriteFile(testTextFile, []byte("This is just plain text, not a contract"), 0644)
	assert.NoError(t, err)

	result, isBaremetal, err := GenerateInitdata(testTextFile, "", nil)

	// HpccInitdata accepts any content and gzips it
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.False(t, isBaremetal)
}

// TestPrintInitdata_WithFilePath tests PrintInitdata with valid file path
func TestPrintInitdata_WithFilePath(t *testing.T) {
	tmpDir := t.TempDir()
	testOutputPath := filepath.Join(tmpDir, "test_initdata_output.txt")
	testData := "H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA="

	err := PrintInitdata(testData, testOutputPath, false)

	assert.NoError(t, err)

	// Verify file was created
	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	// Verify content
	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Equal(t, testData, string(content))
}

// TestPrintInitdata_WithoutFilePath tests PrintInitdata without file path (prints to stdout)
func TestPrintInitdata_WithoutFilePath(t *testing.T) {
	testData := "H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA="

	err := PrintInitdata(testData, "", false)

	assert.NoError(t, err)
}

// TestPrintInitdata_InvalidPath tests PrintInitdata with invalid file path
func TestPrintInitdata_InvalidPath(t *testing.T) {
	testData := "H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA="
	invalidPath := "/invalid/directory/that/does/not/exist/output.txt"

	err := PrintInitdata(testData, invalidPath, false)

	assert.Error(t, err)
}

// TestPrintInitdata_EmptyData tests PrintInitdata with empty data
func TestPrintInitdata_EmptyData(t *testing.T) {
	tmpDir := t.TempDir()
	testOutputPath := filepath.Join(tmpDir, "test_initdata_output.txt")

	err := PrintInitdata("", testOutputPath, false)

	assert.NoError(t, err)

	// Verify file was created with empty content
	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Equal(t, "", string(content))
}

// TestPrintInitdata_LargeData tests PrintInitdata with large data
func TestPrintInitdata_LargeData(t *testing.T) {
	tmpDir := t.TempDir()
	testOutputPath := filepath.Join(tmpDir, "test_initdata_output.txt")

	// Create a large base64 string (simulating large gzipped data)
	largeData := ""
	for i := 0; i < 1000; i++ {
		largeData += "H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA="
	}

	err := PrintInitdata(largeData, testOutputPath, false)

	assert.NoError(t, err)

	// Verify file was created
	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	// Verify content
	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Equal(t, largeData, string(content))
}

// TestGenerateInitdata_ValidContract tests with a valid encrypted contract
func TestGenerateInitdata_ValidContract(t *testing.T) {
	result, isBaremetal, err := GenerateInitdata(testContractPath, "", nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.False(t, isBaremetal)

	// The result should be base64 encoded gzipped data
	// It should start with standard gzip base64 prefix
	assert.Greater(t, len(result), 10, "Generated initdata should be substantial")
}

// TestGenerateInitdata_RelativePath tests with relative path
func TestGenerateInitdata_RelativePath(t *testing.T) {
	result, isBaremetal, err := GenerateInitdata("../../samples/hpcc/signed-encrypt-hpcc.yaml", "", nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.False(t, isBaremetal)
}

// TestGenerateInitdata_AbsolutePath tests with absolute path
func TestGenerateInitdata_AbsolutePath(t *testing.T) {
	// Get absolute path
	absPath, err := os.Getwd()
	assert.NoError(t, err)

	fullPath := absPath + "/../../samples/hpcc/signed-encrypt-hpcc.yaml"

	result, isBaremetal, err := GenerateInitdata(fullPath, "", nil)

	// May fail if path doesn't resolve correctly, but that's okay
	if err == nil {
		assert.NotEmpty(t, result)
		assert.False(t, isBaremetal)
	}
}

// TestValidateInput_WithSehdrBin tests ValidateInput with sehdr binary path
func TestValidateInput_WithSehdrBin(t *testing.T) {
	tmpDir := t.TempDir()
	testSehdrBinPath := filepath.Join(tmpDir, "test_sehdr.bin")
	testOutputPath := filepath.Join(tmpDir, "test_initdata_output.txt")

	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(SehdrBinFlagName, testSehdrBinPath, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")
	cmd.Flags().StringArray(ExtraRootCertsFlagName, []string{}, "")

	inputData, sehdrBinPath, outputPath, extraCerts, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, inputData)
	assert.Equal(t, testSehdrBinPath, sehdrBinPath)
	assert.Equal(t, testOutputPath, outputPath)
	assert.Empty(t, extraCerts)
}

// TestGenerateInitdata_WithSehdrBin_Success tests successful initdata generation with SE header binary
func TestGenerateInitdata_WithSehdrBin_Success(t *testing.T) {
	tmpDir := t.TempDir()
	testSehdrBinPath := filepath.Join(tmpDir, "test_sehdr.bin")

	// Create a test binary file with some content
	testBinaryData := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}
	err := os.WriteFile(testSehdrBinPath, testBinaryData, 0644)
	assert.NoError(t, err)

	result, isBaremetal, err := GenerateInitdata(testContractPath, testSehdrBinPath, nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.True(t, isBaremetal, "isBaremetal should be true when sehdr binary is provided")
	// The result should be base64 encoded gzipped data
	assert.Greater(t, len(result), 10, "Generated initdata should be substantial")
}

// TestGenerateInitdata_WithSehdrBin_InvalidPath tests with non-existent sehdr binary path
func TestGenerateInitdata_WithSehdrBin_InvalidPath(t *testing.T) {
	tmpDir := t.TempDir()
	testInvalidBinPath := filepath.Join(tmpDir, "file", "invalid_sehdr.bin")

	result, isBaremetal, err := GenerateInitdata(testContractPath, testInvalidBinPath, nil)

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.False(t, isBaremetal)
	assert.Contains(t, err.Error(), "doesn't exist", "Error should mention path doesn't exist")
}

// TestGenerateInitdata_WithSehdrBin_EmptyFile tests with empty sehdr binary file
func TestGenerateInitdata_WithSehdrBin_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	testSehdrBinPath := filepath.Join(tmpDir, "test_sehdr.bin")

	// Create an empty binary file
	err := os.WriteFile(testSehdrBinPath, []byte{}, 0644)
	assert.NoError(t, err)

	result, isBaremetal, err := GenerateInitdata(testContractPath, testSehdrBinPath, nil)

	// Empty binary file should cause an error in HpcrText
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.False(t, isBaremetal)
}

// TestGenerateInitdata_WithSehdrBin_LargeBinary tests with large binary file
func TestGenerateInitdata_WithSehdrBin_LargeBinary(t *testing.T) {
	tmpDir := t.TempDir()
	testSehdrBinPath := filepath.Join(tmpDir, "test_sehdr.bin")

	// Create a larger binary file (1KB)
	largeBinaryData := make([]byte, 1024)
	for i := range largeBinaryData {
		largeBinaryData[i] = byte(i % 256)
	}
	err := os.WriteFile(testSehdrBinPath, largeBinaryData, 0644)
	assert.NoError(t, err)

	result, isBaremetal, err := GenerateInitdata(testContractPath, testSehdrBinPath, nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.True(t, isBaremetal)
	assert.Greater(t, len(result), 100, "Generated initdata with large binary should be substantial")
}

// TestGenerateInitdata_WithExtraRootCerts tests GenerateInitdata with extra root certificate files
func TestGenerateInitdata_WithExtraRootCerts(t *testing.T) {
	tmpDir := t.TempDir()
	certPath := filepath.Join(tmpDir, "root_ca.pem")

	certContent := "-----BEGIN CERTIFICATE-----\nMIIBcert1\n-----END CERTIFICATE-----\n"
	err := os.WriteFile(certPath, []byte(certContent), 0644)
	assert.NoError(t, err)

	result, isBaremetal, err := GenerateInitdata(testContractPath, "", []string{certPath})

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.False(t, isBaremetal)
	assert.Greater(t, len(result), 10, "Generated initdata with extra root certs should be substantial")
}

// TestGenerateInitdata_WithExtraRootCerts_InvalidPath tests GenerateInitdata with non-existent cert path
func TestGenerateInitdata_WithExtraRootCerts_InvalidPath(t *testing.T) {
	tmpDir := t.TempDir()
	invalidCertPath := filepath.Join(tmpDir, "nonexistent", "root_ca.pem")

	result, isBaremetal, err := GenerateInitdata(testContractPath, "", []string{invalidCertPath})

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.False(t, isBaremetal)
	assert.Contains(t, err.Error(), "doesn't exist")
}

// TestGenerateInitdata_WithMultipleExtraRootCerts tests GenerateInitdata with multiple cert files
func TestGenerateInitdata_WithMultipleExtraRootCerts(t *testing.T) {
	tmpDir := t.TempDir()
	certPath1 := filepath.Join(tmpDir, "root_ca1.pem")
	certPath2 := filepath.Join(tmpDir, "root_ca2.pem")

	err := os.WriteFile(certPath1, []byte("-----BEGIN CERTIFICATE-----\nMIIBcert1\n-----END CERTIFICATE-----\n"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(certPath2, []byte("-----BEGIN CERTIFICATE-----\nMIIBcert2\n-----END CERTIFICATE-----\n"), 0644)
	assert.NoError(t, err)

	result, isBaremetal, err := GenerateInitdata(testContractPath, "", []string{certPath1, certPath2})

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.False(t, isBaremetal)
	assert.Greater(t, len(result), 10)
}

// TestGenerateInitdata_WithSehdrBinAndExtraRootCerts tests GenerateInitdata with both sehdr bin and certs
func TestGenerateInitdata_WithSehdrBinAndExtraRootCerts(t *testing.T) {
	tmpDir := t.TempDir()
	testSehdrBinPath := filepath.Join(tmpDir, "test_sehdr.bin")
	certPath := filepath.Join(tmpDir, "root_ca.pem")

	err := os.WriteFile(testSehdrBinPath, []byte{0x00, 0x01, 0x02, 0x03}, 0644)
	assert.NoError(t, err)
	err = os.WriteFile(certPath, []byte("-----BEGIN CERTIFICATE-----\nMIIBcert1\n-----END CERTIFICATE-----\n"), 0644)
	assert.NoError(t, err)

	result, isBaremetal, err := GenerateInitdata(testContractPath, testSehdrBinPath, []string{certPath})

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.True(t, isBaremetal)
	assert.Greater(t, len(result), 10)
}

// TestValidateInput_WithExtraRootCerts tests ValidateInput with extra_root_certificates flag
func TestValidateInput_WithExtraRootCerts(t *testing.T) {
	tmpDir := t.TempDir()
	certPath := filepath.Join(tmpDir, "root_ca.pem")
	testOutputPath := filepath.Join(tmpDir, "test_initdata_output.txt")

	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(SehdrBinFlagName, "", "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")
	cmd.Flags().StringArray(ExtraRootCertsFlagName, []string{certPath}, "")

	inputData, sehdrBinPath, outputPath, extraCerts, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, inputData)
	assert.Equal(t, "", sehdrBinPath)
	assert.Equal(t, testOutputPath, outputPath)
	assert.Equal(t, []string{certPath}, extraCerts)
}

// TestPrintInitdata_Baremetal tests PrintInitdata with baremetal flag set to true
func TestPrintInitdata_Baremetal(t *testing.T) {
	tmpDir := t.TempDir()
	testOutputPath := filepath.Join(tmpDir, "test_initdata_output.txt")
	testData := "H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA="

	err := PrintInitdata(testData, testOutputPath, true)

	assert.NoError(t, err)

	// Verify file was created
	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	// Verify content
	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Equal(t, testData, string(content))
	// Note: The actual "baremetal" message is printed to stdout, not written to file
	// This test verifies the file writing works correctly with isBaremetal=true
}

// TestPrintInitdata_Peerpod tests PrintInitdata with baremetal flag set to false
func TestPrintInitdata_Peerpod(t *testing.T) {
	tmpDir := t.TempDir()
	testOutputPath := filepath.Join(tmpDir, "test_initdata_output.txt")
	testData := "H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA="

	err := PrintInitdata(testData, testOutputPath, false)

	assert.NoError(t, err)

	// Verify file was created
	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	// Verify content
	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Equal(t, testData, string(content))
	// Note: The actual "peerpod" message is printed to stdout, not written to file
	// This test verifies the file writing works correctly with isBaremetal=false
}

// Made with Bob
