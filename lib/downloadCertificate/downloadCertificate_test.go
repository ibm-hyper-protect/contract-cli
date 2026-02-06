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

package downloadCertificate

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testOutputPath  = "../../build/test_download_certificate_output.json"
	testInvalidPath = "../../build/file/file_not_exists.json"
)

// TestValidateInput_Success tests ValidateInput with valid input and versions
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(FormatFlag, JsonFormat, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")
	versions := []string{"1.0.21", "1.0.22"}

	formatType, certificatePath, err := ValidateInput(cmd, versions)

	assert.NoError(t, err)
	assert.Equal(t, JsonFormat, formatType)
	assert.Equal(t, testOutputPath, certificatePath)
}

// TestValidateInput_WithoutOutputPath tests ValidateInput without output path (optional)
func TestValidateInput_WithoutOutputPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(FormatFlag, JsonFormat, "")
	cmd.Flags().String(OutputFlagName, "", "")
	versions := []string{"1.0.21"}

	formatType, certificatePath, err := ValidateInput(cmd, versions)

	assert.NoError(t, err)
	assert.Equal(t, JsonFormat, formatType)
	assert.Equal(t, "", certificatePath)
}

// TestValidateInput_WithYamlFormat tests ValidateInput with yaml format
func TestValidateInput_WithYamlFormat(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(FormatFlag, "yaml", "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")
	versions := []string{"1.0.21"}

	formatType, certificatePath, err := ValidateInput(cmd, versions)

	assert.NoError(t, err)
	assert.Equal(t, "yaml", formatType)
	assert.Equal(t, testOutputPath, certificatePath)
}

// TestValidateInput_WithTextFormat tests ValidateInput with text format
func TestValidateInput_WithTextFormat(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(FormatFlag, "text", "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")
	versions := []string{"1.0.21"}

	formatType, certificatePath, err := ValidateInput(cmd, versions)

	assert.NoError(t, err)
	assert.Equal(t, "text", formatType)
	assert.Equal(t, testOutputPath, certificatePath)
}

// TestValidateInput_WithMultipleVersions tests ValidateInput with multiple versions
func TestValidateInput_WithMultipleVersions(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(FormatFlag, JsonFormat, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")
	versions := []string{"1.0.21", "1.0.22", "1.0.23"}

	formatType, certificatePath, err := ValidateInput(cmd, versions)

	assert.NoError(t, err)
	assert.Equal(t, JsonFormat, formatType)
	assert.Equal(t, testOutputPath, certificatePath)
}

// TestValidateInput_FlagErrors tests error handling for flag retrieval
func TestValidateInput_FlagErrors(t *testing.T) {
	cmd := &cobra.Command{}
	versions := []string{"1.0.21"}
	_, _, err := ValidateInput(cmd, versions)
	assert.Error(t, err)
}

// TestOutput_ToFile tests writing certificate data to file
func TestOutput_ToFile(t *testing.T) {
	testData := `{"certificates": [{"version": "1.0.21", "data": "test"}]}`
	os.Remove(testOutputPath)

	err := Output(testData, testOutputPath)
	assert.NoError(t, err)

	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Equal(t, testData, string(content))
	os.Remove(testOutputPath)
}

// TestOutput_ToStdout tests printing to stdout (empty path)
func TestOutput_ToStdout(t *testing.T) {
	testData := `{"certificates": [{"version": "1.0.21", "data": "test"}]}`
	err := Output(testData, "")
	assert.NoError(t, err)
}

// TestOutput_InvalidPath tests with invalid output path
func TestOutput_InvalidPath(t *testing.T) {
	testData := `{"certificates": [{"version": "1.0.21", "data": "test"}]}`
	err := Output(testData, testInvalidPath)
	assert.Error(t, err)
}

// TestOutput_EmptyData tests with empty certificate data
func TestOutput_EmptyData(t *testing.T) {
	testData := ""
	os.Remove(testOutputPath)
	err := Output(testData, testOutputPath)
	assert.NoError(t, err)

	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)

	content, readErr := os.ReadFile(testOutputPath)
	assert.NoError(t, readErr)
	assert.Equal(t, "", string(content))
	os.Remove(testOutputPath)
}
