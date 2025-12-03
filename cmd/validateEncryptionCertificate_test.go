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
	// active.crt will be valid up to November 9, 2030
	sampleEncryptionCertificate = "../samples/certificate/active.crt"
)

var (
	sampleValidEncryptionCertificateCommand = []string{common.ValidateEncryptionCertParamName, "--in", sampleEncryptionCertificate}
)

// Testcase to check if validate-encryption-certificate is able to validate encryption certificate
func TestValidateEncryptionCertCmdSucess(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidEncryptionCertificateCommand)
	err := validateNetworkConfigCmd.Execute()
	assert.NoError(t, err)
}
