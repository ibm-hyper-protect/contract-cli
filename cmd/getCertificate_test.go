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
