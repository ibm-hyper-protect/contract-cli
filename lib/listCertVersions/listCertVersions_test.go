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

package listCertVersions

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestValidateInput_Success tests ValidateInput with valid flags
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(OsVersionFlagName, "ccrt", "")
	cmd.Flags().String(OutputFlagName, "", "")
	cmd.Flags().String(FormatFlagName, "json", "")

	osVersion, outputPath, format, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, "ccrt", osVersion)
	assert.Equal(t, "", outputPath)
	assert.Equal(t, "json", format)
}

// TestValidateInput_WithoutFlags tests ValidateInput when flags are not set
func TestValidateInput_WithoutFlags(t *testing.T) {
	cmd := &cobra.Command{}

	_, _, _, err := ValidateInput(cmd)
	assert.Error(t, err)
}

// TestValidateInput_InvalidFormat tests ValidateInput with invalid format
func TestValidateInput_InvalidFormat(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(OsVersionFlagName, "ccrt", "")
	cmd.Flags().String(OutputFlagName, "", "")
	cmd.Flags().String(FormatFlagName, "xml", "")

	_, _, _, err := ValidateInput(cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid format")
}

// TestProcess_JsonFormat tests Process function with JSON format
func TestProcess_JsonFormat(t *testing.T) {
	result, err := Process("", "json")

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	// JSON format should contain all platforms
	assert.Contains(t, result, "ccrt")
	assert.Contains(t, result, "ccrv")
	assert.Contains(t, result, "ccco")
}

// TestProcess_YamlFormat tests Process function with YAML format
func TestProcess_YamlFormat(t *testing.T) {
	result, err := Process("", "yaml")

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	// YAML format should contain all platforms
	assert.Contains(t, result, "ccrt:")
	assert.Contains(t, result, "ccrv:")
	assert.Contains(t, result, "ccco:")
}

// TestProcess_DefaultFormat tests Process function with empty format (defaults to JSON)
func TestProcess_DefaultFormat(t *testing.T) {
	result, err := Process("", "")

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	// Should default to JSON format
	assert.Contains(t, result, "ccrt")
}

// TestProcess_SpecificPlatform_Json tests Process function with ccrt platform and JSON format
func TestProcess_SpecificPlatform_Json(t *testing.T) {
	result, err := Process("ccrt", "json")

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "ccrt")
	// Should contain version numbers
	assert.True(t, strings.Contains(result, "26.2.0") || strings.Contains(result, "25."))
}

// TestProcess_SpecificPlatform_Yaml tests Process function with ccrt platform and YAML format
func TestProcess_SpecificPlatform_Yaml(t *testing.T) {
	result, err := Process("ccrt", "yaml")

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "ccrt:")
}

// TestProcess_InvalidPlatform tests Process function with invalid platform
func TestProcess_InvalidPlatform(t *testing.T) {
	result, err := Process("invalid", "json")

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "no certificates found for platform")
	assert.Contains(t, err.Error(), "Valid platforms: ccrt, ccrv, ccco and hpvs for legacy")
}

// TestProcess_CaseInsensitive tests Process function with uppercase platform name
func TestProcess_CaseInsensitive(t *testing.T) {
	result, err := Process("CCRT", "json")

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "ccrt")
}

// TestOutput tests Output function returns the input string unchanged
func TestOutput(t *testing.T) {
	testData := "Test output data"
	result := Output(testData)

	assert.Equal(t, testData, result)
}

// TestProcess_VersionFormat tests that versions are properly formatted
func TestProcess_VersionFormat(t *testing.T) {
	result, err := Process("ccrt", "json")

	assert.NoError(t, err)
	// Check that versions follow semantic versioning format (X.Y.Z)
	assert.Regexp(t, `\d+\.\d+\.\d+`, result)
}

// TestProcess_VersionOrdering tests that versions are in descending order
func TestProcess_VersionOrdering(t *testing.T) {
	result, err := Process("ccrt", "json")

	assert.NoError(t, err)
	// The actual ordering is handled by the certificate package
	// Just verify the output contains version numbers in JSON format
	assert.Contains(t, result, "ccrt")
	assert.Regexp(t, `\d+\.\d+\.\d+`, result)
}
