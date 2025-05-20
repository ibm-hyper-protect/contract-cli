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
	sampleComposeFolderPath      = "../samples/tgz"
	sampleTgzOutputPathPlain     = "../build/base64.tgz.txt"
	sampleTgzOutputPathEncrypted = "../build/base64.tgz.enc"
)

var (
	sampleValidBase64TgzPlainCommand     = []string{common.Base64TgzParamName, "--in", sampleComposeFolderPath, "--output", common.Base64TgzOutputFormatUnencrypted, "--out", sampleTgzOutputPathPlain}
	sampleValidBase64TgzEncryptedCommand = []string{common.Base64TgzParamName, "--in", sampleComposeFolderPath, "--output", common.Base64TgzOutputFormatencrypted, "--out", sampleTgzOutputPathEncrypted}
)

func TestBase64TgzCmdSuccessPlain(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidBase64TgzPlainCommand)
	err := base64TgzCmd.Execute()

	assert.NoError(t, err)
}

func TestBase64TgzCmdSuccessEncrypted(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidBase64TgzEncryptedCommand)
	err := base64TgzCmd.Execute()

	assert.NoError(t, err)
}
