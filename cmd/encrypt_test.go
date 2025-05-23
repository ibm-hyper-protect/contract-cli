// Copyright (c) 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// Testcase to check if encrypt is able to generate signed and encrypted contract
func TestEncryptCmdText(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleSignEncryptCommand)
	err := encryptCmd.Execute()

	assert.NoError(t, err)
}

// Testcase to check if encrypt is able to generate signed and encrypted contract with contract expiry
func TestEncryptCmdTextContractExpiry(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleContractExpiryCommand)
	err := encryptCmd.Execute()

	assert.NoError(t, err)
}
