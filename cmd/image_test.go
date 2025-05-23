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
	"github.com/stretchr/testify/assert"
)

const (
	sampleTerraformImageJson = "../samples/images/terraform_image.json"
	sampleCliImageJson       = "../samples/images/cli_image.json"
	sampleVersion            = "1.0.22"
	sampleFormatJson         = "json"
	sampleFormatYaml         = "yaml"
	sampleImageOutputJson    = "../build/hpcr_image.json"
	sampleImageOutputYaml    = "../build/hpcr_image.yaml"
)

var (
	sampleValidCommandJson               = []string{common.ImageParamName, "--in", sampleTerraformImageJson, "--version", sampleVersion, "--format", sampleFormatJson, "--out", sampleImageOutputJson}
	sampleValidCommandYaml               = []string{common.ImageParamName, "--in", sampleCliImageJson, "--version", sampleVersion, "--format", sampleFormatYaml, "--out", sampleImageOutputYaml}
	sampleValidCommandJsonWithoutVersion = []string{common.ImageParamName, "--in", sampleTerraformImageJson, "--format", sampleFormatJson, "--out", sampleImageOutputJson}
)

// Testcase to check if image is able to get image details from Terraform output
func TestImageCmdSuccess1(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidCommandJson)
	err := imageCmd.Execute()

	assert.NoError(t, err)
}

// Testcase to check if image is able to get image details from IBM Cloud CLI output
func TestImageCmdSuccess2(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidCommandYaml)
	err := imageCmd.Execute()

	assert.NoError(t, err)
}

// Testcase to check if image is able to get latest image details
func TestImageCmdSuccessWithoutVersion(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidCommandJsonWithoutVersion)
	err := imageCmd.Execute()

	assert.NoError(t, err)
}
