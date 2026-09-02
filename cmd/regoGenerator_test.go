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

// TestRegoGeneratorCmd_Success tests the rego-generator command — default format (base64)
func TestRegoGeneratorCmd_Success(t *testing.T) {
	tmpDir := t.TempDir()
	outputStem := filepath.Join(tmpDir, "policy")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{"rego-generator", "--in", testRegoGeneratorPodFile, "--out", outputStem})

	err := regoGeneratorCmd.Execute()
	assert.NoError(t, err)

	// Default format is base64: only _base64 file must exist
	b64File := filepath.Join(tmpDir, "policy_base64")
	_, err = os.Stat(b64File)
	assert.NoError(t, err, "_base64 file must be created with default format")

	// Plain .rego file must NOT be created
	_, err = os.Stat(filepath.Join(tmpDir, "policy.rego"))
	assert.True(t, os.IsNotExist(err), "plain .rego file must not be created with default format=base64")
}

// TestRegoGeneratorCmd_FormatText tests --format text — only .rego file written
func TestRegoGeneratorCmd_FormatText(t *testing.T) {
	tmpDir := t.TempDir()
	outputStem := filepath.Join(tmpDir, "policy")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{"rego-generator", "--in", testRegoGeneratorPodFile, "--out", outputStem, "--format", "text"})

	err := regoGeneratorCmd.Execute()
	assert.NoError(t, err)

	// Plain .rego file must exist
	regoFile := filepath.Join(tmpDir, "policy.rego")
	content, err := os.ReadFile(regoFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "package agent_policy")
	assert.Contains(t, string(content), "allow_image")
	assert.Contains(t, string(content), "docker\\\\.io/library/nginx:latest")
	assert.Contains(t, string(content), "docker\\\\.io/library/redis:alpine")
	assert.Contains(t, string(content), "quay\\\\.io/prometheus/busybox:latest")
	assert.Contains(t, string(content), "docker\\\\.io/library/postgres:15-alpine")

	// _base64 file must NOT be created
	_, err = os.Stat(filepath.Join(tmpDir, "policy_base64"))
	assert.True(t, os.IsNotExist(err), "_base64 file must not be created for format=text")
}

// TestRegoGeneratorCmd_FormatBase64 tests --format base64 — only _base64 file written
func TestRegoGeneratorCmd_FormatBase64(t *testing.T) {
	tmpDir := t.TempDir()
	outputStem := filepath.Join(tmpDir, "policy")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{"rego-generator", "--in", testRegoGeneratorPodFile, "--out", outputStem, "--format", "base64"})

	err := regoGeneratorCmd.Execute()
	assert.NoError(t, err)

	// _base64 file must exist and be non-empty
	b64File := filepath.Join(tmpDir, "policy_base64")
	b64Content, err := os.ReadFile(b64File)
	assert.NoError(t, err)
	assert.NotEmpty(t, string(b64Content))

	// Plain .rego file must NOT be created
	_, err = os.Stat(filepath.Join(tmpDir, "policy.rego"))
	assert.True(t, os.IsNotExist(err), "plain .rego file must not be created for format=base64")
}

// TestRegoGeneratorCmd_FormatBoth tests --format both — both files written
func TestRegoGeneratorCmd_FormatBoth(t *testing.T) {
	tmpDir := t.TempDir()
	outputStem := filepath.Join(tmpDir, "policy")

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{"rego-generator", "--in", testRegoGeneratorPodFile, "--out", outputStem, "--format", "both"})

	err := regoGeneratorCmd.Execute()
	assert.NoError(t, err)

	// Plain .rego file must exist with policy content
	regoFile := filepath.Join(tmpDir, "policy.rego")
	content, err := os.ReadFile(regoFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "package agent_policy")
	assert.Contains(t, string(content), "allow_image")

	// _base64 file must also exist and be non-empty
	b64File := filepath.Join(tmpDir, "policy_base64")
	b64Content, err := os.ReadFile(b64File)
	assert.NoError(t, err)
	assert.NotEmpty(t, string(b64Content))
}

// TestRegoGeneratorCmd_CommandProperties tests command properties and flag registration
func TestRegoGeneratorCmd_CommandProperties(t *testing.T) {
	assert.NotNil(t, regoGeneratorCmd)
	assert.Contains(t, regoGeneratorCmd.Use, "rego-generator")
	assert.NotEmpty(t, regoGeneratorCmd.Short)
	assert.NotEmpty(t, regoGeneratorCmd.Long)

	// Verify --format flag is registered with the correct default
	formatFlag := regoGeneratorCmd.PersistentFlags().Lookup("format")
	assert.NotNil(t, formatFlag, "--format flag must be registered")
	assert.Equal(t, "base64", formatFlag.DefValue, "--format default must be 'base64'")
}
