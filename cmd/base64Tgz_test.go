package cmd

import (
	"bytes"
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/stretchr/testify/assert"
)

const (
	sampleComposeFolderPath      = "../samples/tgz"
	sampleTgzOutputPathPlain     = "../build/base64.tgz.txt"
	sampleTgzOutputPathEncrypted = "../build/base64.tgz.enc"
)

var (
	sampleValidBase64TgzPlainCommand     = []string{common.Base64TgzParamName, "--in", sampleComposeFolderPath, "--output", common.Base64TgzOutputFormatUnencrypted, "--out", sampleTgzOutputPathPlain}
	sampleValidBase64TgzEncryptedCommand = []string{common.Base64TgzParamName, "--in", sampleComposeFolderPath, "--output", common.Base64TgzOutputFormatencrypted, "--out", sampleTgzOutputPathEncrypted}
)

func TestBase64TgzCmdSuccessPlain(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidBase64TgzPlainCommand)
	err := base64TgzCmd.Execute()

	assert.NoError(t, err)
}

func TestBase64TgzCmdSuccessEncrypted(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidBase64TgzEncryptedCommand)
	err := base64TgzCmd.Execute()

	assert.NoError(t, err)
}
