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

package image

import (
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

const (
	testApiImagePath  = "../../samples/images/api_image.json"
	testCliImagePath  = "../../samples/images/cli_image.json"
	testTerraformPath = "../../samples/images/terraform_image.json"
	testInvalidPath   = "../../build/file/file_not_exist.json"
	testVersion       = "1.0.0"
)

// TestValidateInput_Success tests ValidateInput with all required flags
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testApiImagePath, "")
	cmd.Flags().String(VersionFlagName, testVersion, "")
	cmd.Flags().String(FormatFlag, JsonFormat, "")
	cmd.Flags().String(OutputFlagName, "output.json", "")

	imageListJsonPath, versionName, formatType, hpcrImagePath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testApiImagePath, imageListJsonPath)
	assert.Equal(t, testVersion, versionName)
	assert.Equal(t, JsonFormat, formatType)
	assert.Equal(t, "output.json", hpcrImagePath)
}

// TestValidateInput_WithYamlFormat tests ValidateInput with YAML format
func TestValidateInput_WithYamlFormat(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testApiImagePath, "")
	cmd.Flags().String(VersionFlagName, testVersion, "")
	cmd.Flags().String(FormatFlag, YamlFormat, "")
	cmd.Flags().String(OutputFlagName, "output.yaml", "")

	imageListJsonPath, versionName, formatType, hpcrImagePath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testApiImagePath, imageListJsonPath)
	assert.Equal(t, testVersion, versionName)
	assert.Equal(t, YamlFormat, formatType)
	assert.Equal(t, "output.yaml", hpcrImagePath)
}

// TestValidateInput_WithoutVersion tests ValidateInput without version (optional)
func TestValidateInput_WithoutVersion(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testApiImagePath, "")
	cmd.Flags().String(VersionFlagName, "", "")
	cmd.Flags().String(FormatFlag, JsonFormat, "")
	cmd.Flags().String(OutputFlagName, "output.json", "")

	imageListJsonPath, versionName, formatType, hpcrImagePath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testApiImagePath, imageListJsonPath)
	assert.Equal(t, "", versionName)
	assert.Equal(t, JsonFormat, formatType)
	assert.Equal(t, "output.json", hpcrImagePath)
}

// TestValidateInput_WithoutOutputPath tests ValidateInput without output path (optional)
func TestValidateInput_WithoutOutputPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testApiImagePath, "")
	cmd.Flags().String(VersionFlagName, testVersion, "")
	cmd.Flags().String(FormatFlag, JsonFormat, "")
	cmd.Flags().String(OutputFlagName, "", "")

	imageListJsonPath, versionName, formatType, hpcrImagePath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testApiImagePath, imageListJsonPath)
	assert.Equal(t, testVersion, versionName)
	assert.Equal(t, JsonFormat, formatType)
	assert.Equal(t, "", hpcrImagePath)
}

// TestValidateInput_WithoutFlags tests ValidateInput when flags are not set
func TestValidateInput_WithoutFlags(t *testing.T) {
	cmd := &cobra.Command{}
	_, _, _, _, err := ValidateInput(cmd)
	assert.Error(t, err)
}

// TestValidateInput_EmptyFormat tests ValidateInput with empty format
func TestValidateInput_EmptyFormat(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testApiImagePath, "")
	cmd.Flags().String(VersionFlagName, testVersion, "")
	cmd.Flags().String(FormatFlag, "", "")
	cmd.Flags().String(OutputFlagName, "output.json", "")

	imageListJsonPath, versionName, formatType, hpcrImagePath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testApiImagePath, imageListJsonPath)
	assert.Equal(t, testVersion, versionName)
	assert.Equal(t, "", formatType)
	assert.Equal(t, "output.json", hpcrImagePath)
}

// TestProcess_Success tests successful image processing
func TestProcess_Success(t *testing.T) {
	result, err := Process(testApiImagePath, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, result.Id)
	assert.NotEmpty(t, result.Name)
	assert.NotEmpty(t, result.Checksum)
}

// TestProcess_WithCliImage tests processing CLI image format
func TestProcess_WithCliImage(t *testing.T) {
	result, err := Process(testCliImagePath, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, result.Id)
	assert.NotEmpty(t, result.Name)
	assert.NotEmpty(t, result.Checksum)
}

// TestProcess_WithTerraformImage tests processing Terraform image format
func TestProcess_WithTerraformImage(t *testing.T) {
	result, err := Process(testTerraformPath, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, result.Id)
	assert.NotEmpty(t, result.Name)
	assert.NotEmpty(t, result.Checksum)
}

// TestOutput_JsonFormat tests Output function with JSON format
func TestOutput_JsonFormat(t *testing.T) {
	testImage := ImageDetails{
		Id:       "r006-test-id",
		Name:     "test-image",
		Checksum: "abc123",
		Version:  "1.0.0",
	}

	result, err := Output(testImage, JsonFormat)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "r006-test-id")
	assert.Contains(t, result, "test-image")
	assert.Contains(t, result, "abc123")
	assert.Contains(t, result, "1.0.0")

	var parsed ImageDetails
	jsonErr := json.Unmarshal([]byte(result), &parsed)
	assert.NoError(t, jsonErr)
	assert.Equal(t, testImage.Id, parsed.Id)
	assert.Equal(t, testImage.Name, parsed.Name)
	assert.Equal(t, testImage.Checksum, parsed.Checksum)
	assert.Equal(t, testImage.Version, parsed.Version)
}

// TestOutput_YamlFormat tests Output function with YAML format
func TestOutput_YamlFormat(t *testing.T) {
	testImage := ImageDetails{
		Id:       "r006-test-id",
		Name:     "test-image",
		Checksum: "abc123",
		Version:  "1.0.0",
	}

	result, err := Output(testImage, YamlFormat)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "r006-test-id")
	assert.Contains(t, result, "test-image")
	assert.Contains(t, result, "abc123")
	assert.Contains(t, result, "1.0.0")

	var parsed ImageDetails
	yamlErr := yaml.Unmarshal([]byte(result), &parsed)
	assert.NoError(t, yamlErr)
	assert.Equal(t, testImage.Id, parsed.Id)
	assert.Equal(t, testImage.Name, parsed.Name)
	assert.Equal(t, testImage.Checksum, parsed.Checksum)
	assert.Equal(t, testImage.Version, parsed.Version)
}

// TestOutput_InvalidFormat tests Output function with invalid format
func TestOutput_InvalidFormat(t *testing.T) {
	testImage := ImageDetails{
		Id:       "r006-test-id",
		Name:     "test-image",
		Checksum: "abc123",
		Version:  "1.0.0",
	}

	result, err := Output(testImage, "invalid")
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), invalidFormatMessage)
}

// TestOutput_EmptyFormat tests Output function with empty format
func TestOutput_EmptyFormat(t *testing.T) {
	testImage := ImageDetails{
		Id:       "r006-test-id",
		Name:     "test-image",
		Checksum: "abc123",
		Version:  "1.0.0",
	}

	result, err := Output(testImage, "")
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), invalidFormatMessage)
}

// TestOutput_EmptyImageDetails tests Output with empty image details
func TestOutput_EmptyImageDetails(t *testing.T) {
	testImage := ImageDetails{}
	result, err := Output(testImage, JsonFormat)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	var parsed ImageDetails
	jsonErr := json.Unmarshal([]byte(result), &parsed)
	assert.NoError(t, jsonErr)
}

// TestOutput_YamlFormatStructure tests YAML output structure
func TestOutput_YamlFormatStructure(t *testing.T) {
	testImage := ImageDetails{
		Id:       "r006-test-id",
		Name:     "test-image",
		Checksum: "abc123",
		Version:  "1.0.0",
	}

	result, err := Output(testImage, YamlFormat)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "id: r006-test-id")
	assert.Contains(t, result, "name: test-image")
	assert.Contains(t, result, "checksum: abc123")
	assert.Contains(t, result, "version: 1.0.0")
}

// TestImageDetails_EmptyFieldsJsonMarshaling tests JSON marshaling with empty fields
func TestImageDetails_EmptyFieldsJsonMarshaling(t *testing.T) {
	testImage := ImageDetails{
		Id: "test-id",
	}

	jsonData, err := json.Marshal(testImage)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var parsed ImageDetails
	err = json.Unmarshal(jsonData, &parsed)
	assert.NoError(t, err)
	assert.Equal(t, testImage.Id, parsed.Id)
	assert.Equal(t, "", parsed.Name)
	assert.Equal(t, "", parsed.Checksum)
	assert.Equal(t, "", parsed.Version)
}
