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
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/lib/validateCertChain"
	"github.com/stretchr/testify/assert"
)

const (
	// Using existing certificate files for testing
	sampleLeafCert         = "../samples/certificate/active.crt"
	sampleIntermediateCert = "../samples/contract-expiry/personal_ca.pem"
	sampleRootCert         = "../samples/contract-expiry/personal_ca.crt"
)

var (
	sampleValidateCertChainCommand = []string{
		validateCertChain.ParameterName,
		"--cert", sampleLeafCert,
		"--intermediate", sampleIntermediateCert,
		"--root", sampleRootCert,
	}
)

// Testcase to check if validate-cert-chain command is registered
func TestValidateCertChainCmdRegistered(t *testing.T) {
	// Verify the command is registered with root
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == validateCertChain.ParameterName {
			found = true
			// Verify command has required flags
			assert.NotNil(t, cmd.PersistentFlags().Lookup("cert"), "cert flag should exist")
			assert.NotNil(t, cmd.PersistentFlags().Lookup("intermediate"), "intermediate flag should exist")
			assert.NotNil(t, cmd.PersistentFlags().Lookup("root"), "root flag should exist")
			break
		}
	}
	assert.True(t, found, "validate-cert-chain command should be registered")
}
