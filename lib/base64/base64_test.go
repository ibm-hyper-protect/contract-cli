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

package base64

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testInputText   = "hello world"
	testInputJson   = `{"key":"value"}`
	testOutputPath  = "../../build/test_base64_output.txt"
	testInvalidPath = "../../build/file/file_not_exists.txt"
)

// TestValidateInput_Success tests ValidateInput with valid input
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testInputText, "")
	cmd.Flags().String(FormatFlagName, TextFormat, "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")

	inputData, formatType, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testInputText, inputData)
	assert.Equal(t, TextFormat, formatType)
	assert.Equal(t, testOutputPath, outputPath)
}

// TestProcess_Success tests Process function with valid text input
func TestProcess_Success(t *testing.T) {
	result, err := Process(testInputText, TextFormat)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestProcess_InvalidFormat tests Process function with invalid format
func TestProcess_InvalidFormat(t *testing.T) {
	result, err := Process(testInputText, "invalid")

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), invalidInputMessageBase64)
}

// TestOutput_Success tests Output function with valid file path
func TestOutput_Success(t *testing.T) {
	testData := "dGVzdCBkYXRh"
	os.Remove(testOutputPath)
	err := Output(testData, testOutputPath)

	assert.NoError(t, err)
	_, statErr := os.Stat(testOutputPath)
	assert.NoError(t, statErr)
	os.Remove(testOutputPath)
}

// TestOutput_InvalidPath tests Output function with invalid file path
func TestOutput_InvalidPath(t *testing.T) {
	testData := "dGVzdCBkYXRh"
	err := Output(testData, testInvalidPath)
	assert.Error(t, err)
}
