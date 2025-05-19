package cmd

import (
	"bytes"
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/stretchr/testify/assert"
)

const (
	sampleEncryptedAttestationFilePath = "../samples/attestation/se-checksums.txt.enc"
	samplePrivateKeyFilePath           = "../samples/attestation/private.pem"
	sampleDecryptedAttestationFilePath = "../build/se-checksums.txt"
)

var (
	sampleValidCommand = []string{common.DecryptAttestParamName, "--in", sampleEncryptedAttestationFilePath, "--priv", samplePrivateKeyFilePath, "--out", sampleDecryptedAttestationFilePath}
)

func TestDecryptAttestationCmdSucess(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidCommand)
	err := decryptAttestationCmd.Execute()

	assert.NoError(t, err)
}
