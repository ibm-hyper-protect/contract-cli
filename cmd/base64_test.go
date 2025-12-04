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
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-cli/lib/base64"
	"github.com/stretchr/testify/assert"
)

const (
	sampleBase64InputText  = "hello"
	sampleBase64InputJson  = `{"dev":"sash"}`
	sampleBase64OutputPath = "../build/base64.txt"
)

var (
	sampleValidBase64Command     = []string{base64.ParameterName, "--in", sampleBase64InputText, "--out", sampleBase64OutputPath}
	sampleValidBase64JsonCommand = []string{base64.ParameterName, "--in", sampleBase64InputJson, "--format", common.DataFormatJson, "--out", sampleBase64OutputPath}
)

// Testcase to check if base64 command is working with plain text
func TestBase64CmdSuccess1(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidBase64Command)
	err := base64Cmd.Execute()

	assert.NoError(t, err)
}

// Testcase to check if base64 command is working with JSON input
func TestBase64CmdSuccess2(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidBase64JsonCommand)
	err := base64Cmd.Execute()

	assert.NoError(t, err)
}
