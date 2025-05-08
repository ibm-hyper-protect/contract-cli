package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	sampleTerraformImageJson = "../samples/images/terraform_image.json"
	sampleCliImageJson       = "../samples/images/cli_image.json"
	sampleVersion            = "1.0.22"
	sampleFormatJson         = "json"
	sampleFormatYaml         = "yaml"
	sampleOutputJson         = "../build/hpcr_image.json"
	sampleOutputYaml         = "../build/hpcr_image.yaml"
)

var (
	sampleValidCommandJson = []string{
		"image", "--in", sampleTerraformImageJson,
		"--version", sampleVersion, "--format", sampleFormatJson,
		"--out", sampleOutputJson}
	sampleValidCommandYaml = []string{
		"image", "--in", sampleCliImageJson,
		"--version", sampleVersion, "--format", sampleFormatYaml,
		"--out", sampleOutputYaml}
	sampleValidCommandJsonWithoutVersion = []string{
		"image", "--in", sampleTerraformImageJson,
		"--format", sampleFormatJson, "--out", sampleOutputJson}
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
