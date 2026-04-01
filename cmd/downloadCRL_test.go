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

	"github.com/ibm-hyper-protect/contract-cli/lib/downloadCRL"
	"github.com/stretchr/testify/assert"
)

const (
	// Note: This is a test URL and will fail in actual execution
	// Real tests would need a valid CRL distribution point
	sampleCRLURL        = "http://example.com/test.crl"
	sampleCRLOutputPath = "../build/test-crl.pem"
)

var (
	sampleDownloadCRLCommand = []string{
		downloadCRL.ParameterName,
		"--url", sampleCRLURL,
		"--out", sampleCRLOutputPath,
	}
)

// Testcase to check if download-crl command is registered
func TestDownloadCRLCmdRegistered(t *testing.T) {
	// Verify the command is registered with root
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == downloadCRL.ParameterName {
			found = true
			// Verify command has required flags
			assert.NotNil(t, cmd.PersistentFlags().Lookup("url"), "url flag should exist")
			assert.NotNil(t, cmd.PersistentFlags().Lookup("out"), "out flag should exist")
			break
		}
	}
	assert.True(t, found, "download-crl command should be registered")
}
