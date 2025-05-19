package cmd

import (
	"bytes"
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/stretchr/testify/assert"
)

const (
	sampleDownloadCertificateFormat = "yaml"
	sampleEncryptionCertPath        = "../build/enc-cert.yaml"
)

var (
	sampleValidDownloadCertCommand = []string{common.DownloadCertParamName, "--version", "1.0.21,1.0.22", "--format", sampleDownloadCertificateFormat, "--out", sampleEncryptionCertPath}
)

func TestDownloadCertificatesCmdSuccess(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidDownloadCertCommand)
	err := downloadCertificatesCmd.Execute()

	assert.NoError(t, err)
}
