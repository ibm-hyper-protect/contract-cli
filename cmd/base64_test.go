package cmd

import (
	"bytes"
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/stretchr/testify/assert"
)

const (
	sampleBase64InputText  = "hello"
	sampleBase64InputJson  = `{"dev":"sash"}`
	sampleBase64OutputPath = "../build/base64.txt"
)

var (
	sampleValidBase64Command     = []string{common.Base64ParamName, "--in", sampleBase64InputText, "--out", sampleBase64OutputPath}
	sampleValidBase64JsonCommand = []string{common.Base64ParamName, "--in", sampleBase64InputJson, "--format", common.DataFormatJson, "--out", sampleBase64OutputPath}
)

func TestBase64CmdSuccess1(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidBase64Command)
	err := base64Cmd.Execute()

	assert.NoError(t, err)
}

func TestBase64CmdSuccess2(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidBase64JsonCommand)
	err := base64Cmd.Execute()

	assert.NoError(t, err)
}
