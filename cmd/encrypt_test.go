package cmd

import (
	"bytes"
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/stretchr/testify/assert"
)

const (
	sampleSignEncryptInputPath   = "../samples/contract.yaml"
	sampleSignEncryptPrivatePath = "../samples/sign/private.pem"
	sampleSignEncryptOutputPath  = "../build/signed-encrypt.yaml"
)

var (
	sampleSignEncryptCommand = []string{common.EncryptParamName, "--in", sampleSignEncryptInputPath, "--priv", sampleSignEncryptPrivatePath, "--out", sampleSignEncryptOutputPath}
)

func TestEncryptCmdText(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleSignEncryptCommand)
	err := encryptCmd.Execute()

	assert.NoError(t, err)
}
