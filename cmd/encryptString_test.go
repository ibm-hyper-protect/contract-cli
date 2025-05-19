package cmd

import (
	"bytes"
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/stretchr/testify/assert"
)

const (
	sampleEncryptStringInput = "testing"
	sampleEncryptStringJson  = `{"type":"workload"}`

	sampleEncryptStringFormatText = "text"
	sampleEncryptStringFormatJson = "json"

	sampleEncryptStringOutputYaml = "../build/encrypt_string.txt"
	sampleEncryptStringOutputJson = "../build/encrypt_string.json"
)

var (
	sampleEncryptStringValidCommandText = []string{common.EncryptStrParamName, "--in", sampleEncryptStringInput, "--format", sampleEncryptStringFormatText, "--out", sampleEncryptStringOutputYaml}
	sampleEncryptStringValidCommandJson = []string{common.EncryptStrParamName, "--in", sampleEncryptStringJson, "--format", sampleEncryptStringFormatJson, "--out", sampleEncryptStringOutputJson}
)

func TestEncryptStringCmdText(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleEncryptStringValidCommandText)
	err := encryptStringCmd.Execute()

	assert.NoError(t, err)
}

func TestEncryptStringCmdYaml(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleEncryptStringValidCommandJson)
	err := encryptStringCmd.Execute()

	assert.NoError(t, err)
}
