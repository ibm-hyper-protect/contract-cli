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

package downloadCRL

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testURL        = "http://example.com/crl.pem"
	testOutputPath = "../../build/test-crl.pem"
)

// TestValidateInput_Success tests ValidateInput with valid input
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(URLFlagName, testURL, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")

	url, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testURL, url)
	assert.Equal(t, testOutputPath, outputPath)
}

// TestValidateInput_WithoutOutput tests ValidateInput without output path
func TestValidateInput_WithoutOutput(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(URLFlagName, testURL, "")
	cmd.Flags().String(OutputFlagName, "", "")

	url, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testURL, url)
	assert.Equal(t, "", outputPath)
}

// TestValidateInput_FlagErrors tests error handling for flag retrieval
func TestValidateInput_FlagErrors(t *testing.T) {
	cmd := &cobra.Command{}
	_, _, err := ValidateInput(cmd)
	assert.Error(t, err)
}

// TestOutput_ToFile tests writing CRL to file
func TestOutput_ToFile(t *testing.T) {
	testData := "-----BEGIN X509 CRL-----\ntest crl data\n-----END X509 CRL-----"
	testFile := "../../build/test-output-crl.pem"

	// Clean up before test
	os.Remove(testFile)

	err := Output(testData, testFile)
	assert.NoError(t, err)

	// Verify file was created
	_, statErr := os.Stat(testFile)
	assert.NoError(t, statErr)

	// Clean up after test
	os.Remove(testFile)
}

// TestOutput_ToStdout tests writing CRL to stdout (no output path)
func TestOutput_ToStdout(t *testing.T) {
	testData := "-----BEGIN X509 CRL-----\ntest crl data\n-----END X509 CRL-----"

	err := Output(testData, "")
	assert.NoError(t, err)
}
