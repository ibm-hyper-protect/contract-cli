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

	"github.com/ibm-hyper-protect/contract-cli/lib/listCertVersions"
	"github.com/stretchr/testify/assert"
)

var (
	sampleListAllCertVersionsCommand  = []string{listCertVersions.ParameterName}
	sampleListCcrtCertVersionsCommand = []string{listCertVersions.ParameterName, "--os", "ccrt"}
	sampleListCcrvCertVersionsCommand = []string{listCertVersions.ParameterName, "--os", "ccrv"}
	sampleListCccoCertVersionsCommand = []string{listCertVersions.ParameterName, "--os", "ccco"}
	sampleListInvalidPlatformCommand  = []string{listCertVersions.ParameterName, "--os", "invalid"}
	sampleListCaseInsensitiveCommand  = []string{listCertVersions.ParameterName, "--os", "CCRT"}
)

// TestListCertVersionsCmdSuccess tests listing all certificate versions
func TestListCertVersionsCmdSuccess(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleListAllCertVersionsCommand)
	err := listCertVersionsCmd.Execute()

	assert.NoError(t, err)
}

// TestListCertVersionsCmdCcrt tests listing CCRT certificate versions
func TestListCertVersionsCmdCcrt(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleListCcrtCertVersionsCommand)
	err := listCertVersionsCmd.Execute()

	assert.NoError(t, err)
}

// TestListCertVersionsCmdCcrv tests listing CCRV certificate versions
func TestListCertVersionsCmdCcrv(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleListCcrvCertVersionsCommand)
	err := listCertVersionsCmd.Execute()

	assert.NoError(t, err)
}

// TestListCertVersionsCmdCcco tests listing CCCO certificate versions
func TestListCertVersionsCmdCcco(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleListCccoCertVersionsCommand)
	err := listCertVersionsCmd.Execute()

	assert.NoError(t, err)
}

// TestListCertVersionsCmdCaseInsensitive tests case-insensitive platform names
func TestListCertVersionsCmdCaseInsensitive(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleListCaseInsensitiveCommand)
	err := listCertVersionsCmd.Execute()

	assert.NoError(t, err)
}

// Made with Bob
