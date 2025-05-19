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

func TestImageCmdSuccess1(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidCommandJson)
	err := imageCmd.Execute()

	assert.NoError(t, err)
}

func TestImageCmdSuccess2(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidCommandYaml)
	err := imageCmd.Execute()

	assert.NoError(t, err)
}

func TestImageCmdSuccessWithoutVersion(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidCommandJsonWithoutVersion)
	err := imageCmd.Execute()

	assert.NoError(t, err)
}
