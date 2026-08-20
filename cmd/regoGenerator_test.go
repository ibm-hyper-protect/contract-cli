// Copyright (c) 2026 IBM Corp.
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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testRegoGeneratorPodFile = "../samples/hpcc/sample-pod.yaml"
)

// TestRegoGeneratorCmd_Success tests the rego-generator command with valid input
func TestRegoGeneratorCmd_Success(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "policy.rego")

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Set command arguments using sample pod file
	rootCmd.SetArgs([]string{"rego-generator", "--in", testRegoGeneratorPodFile, "--out", outputFile})

	// Execute command
	err := regoGeneratorCmd.Execute()
	assert.NoError(t, err)

	// Verify output file was created
	_, err = os.Stat(outputFile)
	assert.NoError(t, err)

	// Verify output file contains expected content
	content, err := os.ReadFile(outputFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "package agent_policy")
	assert.Contains(t, string(content), "allow_image")
	assert.Contains(t, string(content), "docker\\\\.io/library/nginx:latest")
	assert.Contains(t, string(content), "docker\\\\.io/library/redis:alpine")
	assert.Contains(t, string(content), "quay\\\\.io/prometheus/busybox:latest")
	assert.Contains(t, string(content), "docker\\\\.io/library/postgres:15-alpine")
}

// TestRegoGeneratorCmd_CommandProperties tests command properties
func TestRegoGeneratorCmd_CommandProperties(t *testing.T) {
	assert.NotNil(t, regoGeneratorCmd)
	assert.Contains(t, regoGeneratorCmd.Use, "rego-generator")
	assert.NotEmpty(t, regoGeneratorCmd.Short)
	assert.NotEmpty(t, regoGeneratorCmd.Long)
}
