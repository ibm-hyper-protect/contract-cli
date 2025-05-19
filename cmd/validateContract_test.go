package cmd

import (
	"bytes"
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/stretchr/testify/assert"
)

const (
	sampleValidateContractInput  = "../samples/contract.yaml"
	sampleValidateContractOsType = "hpvs"
)

var (
	sampleValidContractCommand = []string{common.ValidateContractParamName, "--in", sampleValidateContractInput, "--os", sampleValidateContractOsType}
)

func TestValidateContractCmdSucess(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidContractCommand)
	err := validateContractCmd.Execute()

	assert.NoError(t, err)
}
