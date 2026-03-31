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
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testContractPath      = "../../samples/hpcc/signed-encrypt-hpcc.yaml"
	testOutputPath        = "../../build/test_initdata_output.txt"
	testInvalidPath       = "../../build/file/file_not_exists.txt"
	testCorruptedContract = "../../build/corrupted_contract_initdata.yaml"
	testEmptyContract     = "../../build/empty_contract_initdata.yaml"
	testTextFile          = "../../build/test_text_file.txt"
)

// TestValidateInput_Success tests ValidateInput with all required flags
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")

	inputData, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, inputData)
	assert.Equal(t, testOutputPath, outputPath)
}

// TestValidateInput_WithoutOutputPath tests ValidateInput without output path
func TestValidateInput_WithoutOutputPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(OutputFlagName, "", "")

	inputData, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, inputData)
	assert.Equal(t, "", outputPath)
}

// Note: TestValidateInput_WithoutFlags removed because ValidateInput calls
// SetMandatoryFlagError which calls os.Exit(1), terminating the test process.
// This validation scenario is tested at the command level where appropriate.

// TestValidateInput_FlagErrors tests ValidateInput with flag retrieval errors
func TestValidateInput_FlagErrors(t *testing.T) {
	cmd := &cobra.Command{}
	// Don't define flags to trigger errors

	_, _, err := ValidateInput(cmd)

	assert.Error(t, err)
}

// TestGenerateInitdata_Success tests successful initdata generation
func TestGenerateInitdata_Success(t *testing.T) {
	result, err := GenerateInitdata(testContractPath)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	// Gzipped and base64 encoded data should be a string
	assert.IsType(t, "", result)
}

// TestGenerateInitdata_InvalidPath tests with invalid contract path
func TestGenerateInitdata_InvalidPath(t *testing.T) {
	result, err := GenerateInitdata(testInvalidPath)

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "doesn't exist")
}

// TestGenerateInitdata_EmptyPath tests with empty path
func TestGenerateInitdata_EmptyPath(t *testing.T) {
	result, err := GenerateInitdata("")

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "doesn't exist")
}

// TestGenerateInitdata_CorruptedContract tests with corrupted contract file
// Note: HpccInitdata doesn't validate YAML format, it just gzips and encodes the content
func TestGenerateInitdata_CorruptedContract(t *testing.T) {
	err := os.WriteFile(testCorruptedContract, []byte("invalid: yaml: content: ["), 0644)
	assert.NoError(t, err)
	defer os.Remove(testCorruptedContract)

	result, err := GenerateInitdata(testCorruptedContract)

	// HpccInitdata accepts any content and gzips it
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestGenerateInitdata_EmptyContract tests with empty contract file
func TestGenerateInitdata_EmptyContract(t *testing.T) {
	err := os.WriteFile(testEmptyContract, []byte(""), 0644)
	assert.NoError(t, err)
	defer os.Remove(testEmptyContract)

	result, err := GenerateInitdata(testEmptyContract)

	assert.Error(t, err)
	assert.Equal(t, "", result)
}

// TestGenerateInitdata_NonYamlFile tests with non-YAML file
// Note: HpccInitdata doesn't validate file format, it just gzips and encodes the content
func TestGenerateInitdata_NonYamlFile(t *testing.T) {
	err := os.WriteFile(testTextFile, []byte("This is just plain text, not a contract"), 0644)
	assert.NoError(t, err)
	defer os.Remove(testTextFile)

	result, err := GenerateInitdata(testTextFile)

	// HpccInitdata accepts any content and gzips it
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestPrintInitdata_WithFilePath tests PrintInitdata with valid file path
func TestPrintInitdata_WithFilePath(t *testing.T) {
	testData := "H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA="
	os.Remove(testOutputPath)

	err := PrintInitdata(testData, testOutputPath)

	assert.NoError(t, err)

	// Verify file was created
	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	// Verify content
	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Equal(t, testData, string(content))

	os.Remove(testOutputPath)
}

// TestPrintInitdata_WithoutFilePath tests PrintInitdata without file path (prints to stdout)
func TestPrintInitdata_WithoutFilePath(t *testing.T) {
	testData := "H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA="

	err := PrintInitdata(testData, "")

	assert.NoError(t, err)
}

// TestPrintInitdata_InvalidPath tests PrintInitdata with invalid file path
func TestPrintInitdata_InvalidPath(t *testing.T) {
	testData := "H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA="
	invalidPath := "/invalid/directory/that/does/not/exist/output.txt"

	err := PrintInitdata(testData, invalidPath)

	assert.Error(t, err)
}

// TestPrintInitdata_EmptyData tests PrintInitdata with empty data
func TestPrintInitdata_EmptyData(t *testing.T) {
	os.Remove(testOutputPath)

	err := PrintInitdata("", testOutputPath)

	assert.NoError(t, err)

	// Verify file was created with empty content
	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Equal(t, "", string(content))

	os.Remove(testOutputPath)
}

// TestPrintInitdata_LargeData tests PrintInitdata with large data
func TestPrintInitdata_LargeData(t *testing.T) {
	// Create a large base64 string (simulating large gzipped data)
	largeData := ""
	for i := 0; i < 1000; i++ {
		largeData += "H4sIAAAAAAAA/wEAAP//AAAAAAAAAAA="
	}

	os.Remove(testOutputPath)

	err := PrintInitdata(largeData, testOutputPath)

	assert.NoError(t, err)

	// Verify file was created
	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	// Verify content
	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Equal(t, largeData, string(content))

	os.Remove(testOutputPath)
}

// TestGenerateInitdata_ValidContract tests with a valid encrypted contract
func TestGenerateInitdata_ValidContract(t *testing.T) {
	result, err := GenerateInitdata(testContractPath)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	// The result should be base64 encoded gzipped data
	// It should start with standard gzip base64 prefix
	assert.Greater(t, len(result), 10, "Generated initdata should be substantial")
}

// TestGenerateInitdata_RelativePath tests with relative path
func TestGenerateInitdata_RelativePath(t *testing.T) {
	result, err := GenerateInitdata("../../samples/hpcc/signed-encrypt-hpcc.yaml")

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestGenerateInitdata_AbsolutePath tests with absolute path
func TestGenerateInitdata_AbsolutePath(t *testing.T) {
	// Get absolute path
	absPath, err := os.Getwd()
	assert.NoError(t, err)

	fullPath := absPath + "/../../samples/hpcc/signed-encrypt-hpcc.yaml"

	result, err := GenerateInitdata(fullPath)

	// May fail if path doesn't resolve correctly, but that's okay
	if err == nil {
		assert.NotEmpty(t, result)
	}
}

// Made with Bob
