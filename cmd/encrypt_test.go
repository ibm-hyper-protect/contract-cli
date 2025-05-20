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

	sampleContractExpiryPrivatePath = "../samples/contract-expiry/private.pem"
	sampleContractExpiryCaCertPath  = "../samples/contract-expiry/personal_ca.crt"
	sampleContractExpiryCaKeyPath   = "../samples/contract-expiry/personal_ca.pem"
	sampleContractExpiryCsrPath     = "../samples/contract-expiry/csr.pem"
	sampleContractExpiryExpiry      = "100"
	sampleContractExpiryOutputPath  = "../build/signed-encrypt-contract-expiry.yaml"
)

var (
	sampleSignEncryptCommand    = []string{common.EncryptParamName, "--in", sampleSignEncryptInputPath, "--priv", sampleSignEncryptPrivatePath, "--out", sampleSignEncryptOutputPath}
	sampleContractExpiryCommand = []string{common.EncryptParamName, "--contract-expiry", "--in", sampleSignEncryptInputPath, "--priv", sampleContractExpiryPrivatePath, "--cacert", sampleContractExpiryCaCertPath, "--cakey", sampleContractExpiryCaKeyPath, "--csr", sampleContractExpiryCsrPath, "--expiry", sampleContractExpiryExpiry, "--out", sampleContractExpiryOutputPath}
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

func TestEncryptCmdTextContractExpiry(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleContractExpiryCommand)
	err := encryptCmd.Execute()

	assert.NoError(t, err)
}
