package cmd

import (
	"bytes"
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/stretchr/testify/assert"
)

const (
	sampleCertificatePath      = "../samples/certificate/certs.json"
	sampleCertificateVersion   = "1.0.21"
	sampleCertificateOuputPath = "../build/enc_cert.crt"
)

var (
	sampleGetCertificateCommand = []string{common.GetCertParamName, "--in", sampleCertificatePath, "--version", sampleCertificateVersion, "--out", sampleCertificateOuputPath}
)

func TestGetCertificateCmdSuccess(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleGetCertificateCommand)
	err := getCertificateCmd.Execute()

	assert.NoError(t, err)
}
