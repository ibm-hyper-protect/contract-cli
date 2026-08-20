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

package regoGenerator

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testPodFile    = "../../samples/hpcc/sample-pod.yaml"
	testOutputPath = "../../build/test_rego_output.rego"
)

// TestValidateInput_Success tests ValidateInput with all required flags
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, "test-pod.yaml", "")
	cmd.Flags().String(OutputFlagName, testOutputPath, "")

	inputPath, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, "test-pod.yaml", inputPath)
	assert.Equal(t, testOutputPath, outputPath)
}

// TestValidateInput_WithoutOutputPath tests ValidateInput without output path (optional)
func TestValidateInput_WithoutOutputPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, "test-pod.yaml", "")
	cmd.Flags().String(OutputFlagName, "", "")

	inputPath, outputPath, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, "test-pod.yaml", inputPath)
	assert.Equal(t, "", outputPath)
}

// TestProcess_Success tests Process with valid pod YAML
func TestProcess_Success(t *testing.T) {
	policy, policyBase64, err := Process(testPodFile)

	assert.NoError(t, err)
	assert.NotEmpty(t, policy)
	assert.NotEmpty(t, policyBase64)
	assert.Contains(t, policy, "package agent_policy")
	assert.Contains(t, policy, `docker\\.io/library/nginx:latest`)
	assert.Contains(t, policy, `docker\\.io/library/redis:alpine`)
	assert.Contains(t, policy, `docker\\.io/library/postgres:15-alpine`)
	assert.Contains(t, policy, `quay\\.io/prometheus/busybox:latest`)
	assert.Contains(t, policy, "allow_image(image_name)")
	assert.Contains(t, policy, "allow_command(image_name, args)")
	// HpcrText returns standard base64 (not the encrypted hyper-protect-basic prefix)
	_, err = base64.StdEncoding.DecodeString(policyBase64)
	assert.NoError(t, err, "policyBase64 must be valid standard base64")
}

// TestProcess_InvalidFile tests Process with non-existent file
func TestProcess_InvalidFile(t *testing.T) {
	_, _, err := Process("/non/existent/file.yaml")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read YAML file")
}

// TestProcess_InvalidYAML tests Process with invalid YAML content
func TestProcess_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	invalidFile := filepath.Join(tmpDir, "invalid.yaml")
	err := os.WriteFile(invalidFile, []byte("invalid: yaml: content:"), 0644)
	assert.NoError(t, err)

	_, _, err = Process(invalidFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to extract pod spec")
}

// TestOutput_ToStdout tests Output to stdout — only plain policy printed, no base64 file
func TestOutput_ToStdout(t *testing.T) {
	policy := "package agent_policy\n\ndefault CreateContainerRequest := false"

	err := Output("", policy, "hyper-protect-basic.abc123")

	assert.NoError(t, err)
}

// TestOutput_ToFile tests Output to file — plain policy + companion _base64 file
func TestOutput_ToFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.rego")
	policy := "package agent_policy\n\ndefault CreateContainerRequest := false"
	b64Policy := "hyper-protect-basic.abc123"

	err := Output(outputFile, policy, b64Policy)

	assert.NoError(t, err)

	// Verify plain policy file
	content, err := os.ReadFile(outputFile)
	assert.NoError(t, err)
	assert.Equal(t, policy, string(content))

	// Verify companion base64 file: output_base64 (no extension)
	b64File := filepath.Join(tmpDir, "output_base64")
	b64Content, err := os.ReadFile(b64File)
	assert.NoError(t, err)
	assert.Equal(t, b64Policy, string(b64Content))
}

// TestBase64OutputPath tests the companion path derivation
func TestBase64OutputPath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"policy.rego", "policy_base64"},
		{"out/policy.rego", "out/policy_base64"},
		{"/abs/path/policy.rego", "/abs/path/policy_base64"},
		{"policy", "policy_base64"},
		{"policy.tar.gz", "policy.tar_base64"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, base64OutputPath(tt.input))
		})
	}
}

// TestOutput_InvalidPath tests Output with invalid file path
func TestOutput_InvalidPath(t *testing.T) {
	policy := "package agent_policy"

	err := Output("/invalid/path/that/does/not/exist/output.rego", policy, "hyper-protect-basic.abc")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to write Rego policy to file")
}

// TestProcess_WithInitContainer tests Process with pod containing init containers
func TestProcess_WithInitContainer(t *testing.T) {
	podYAML := `apiVersion: v1
kind: Pod
metadata:
  name: app-pod
spec:
  initContainers:
  - name: init
    image: quay.io/prometheus/busybox:latest
    command: ["/bin/sh"]
    args: ["-c", "echo 'Initializing...'"]
  containers:
  - name: app
    image: docker.io/library/redis:alpine
`
	tmpDir := t.TempDir()
	podFile := filepath.Join(tmpDir, "test-pod.yaml")
	err := os.WriteFile(podFile, []byte(podYAML), 0644)
	assert.NoError(t, err)

	policy, _, err := Process(podFile)

	assert.NoError(t, err)
	assert.Contains(t, policy, `quay\\.io/prometheus/busybox:latest`)
	assert.Contains(t, policy, `docker\\.io/library/redis:alpine`)
}

// TestProcess_EmptyPodYAML tests Process with empty YAML
func TestProcess_EmptyPodYAML(t *testing.T) {
	tmpDir := t.TempDir()
	emptyFile := filepath.Join(tmpDir, "empty.yaml")
	err := os.WriteFile(emptyFile, []byte(""), 0644)
	assert.NoError(t, err)

	_, _, err = Process(emptyFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to extract pod spec")
}

// TestProcess_WithDeployment tests Process with Kubernetes Deployment
func TestProcess_WithDeployment(t *testing.T) {
	deploymentYAML := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
spec:
  replicas: 2
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: nginx
        image: docker.io/library/nginx:alpine
        command: ["/bin/sh"]
        args: ["-c", "nginx -g 'daemon off;'"]
`
	tmpDir := t.TempDir()
	deploymentFile := filepath.Join(tmpDir, "deployment.yaml")
	err := os.WriteFile(deploymentFile, []byte(deploymentYAML), 0644)
	assert.NoError(t, err)

	policy, _, err := Process(deploymentFile)

	assert.NoError(t, err)
	assert.NotEmpty(t, policy)
	assert.Contains(t, policy, "package agent_policy")
	assert.Contains(t, policy, `docker\\.io/library/nginx:alpine`)
	assert.Contains(t, policy, "allow_image")
}

// TestProcess_WithStatefulSet tests Process with Kubernetes StatefulSet
func TestProcess_WithStatefulSet(t *testing.T) {
	statefulSetYAML := `apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: test-statefulset
spec:
  serviceName: test
  replicas: 1
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: redis
        image: docker.io/library/redis:alpine
`
	tmpDir := t.TempDir()
	statefulSetFile := filepath.Join(tmpDir, "statefulset.yaml")
	err := os.WriteFile(statefulSetFile, []byte(statefulSetYAML), 0644)
	assert.NoError(t, err)

	policy, _, err := Process(statefulSetFile)

	assert.NoError(t, err)
	assert.NotEmpty(t, policy)
	assert.Contains(t, policy, `docker\\.io/library/redis:alpine`)
}

// TestIntegration_EndToEnd tests the complete workflow
func TestIntegration_EndToEnd(t *testing.T) {
	tmpDir := t.TempDir()

	// Process the pod YAML from samples
	policy, policyBase64, err := Process(testPodFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, policy)
	assert.NotEmpty(t, policyBase64)
	// HpcrText returns standard base64 (not the encrypted hyper-protect-basic prefix)
	_, decErr := base64.StdEncoding.DecodeString(policyBase64)
	assert.NoError(t, decErr, "policyBase64 must be valid standard base64")

	// Output to file — should create policy.rego and policy_base64
	outputFile := filepath.Join(tmpDir, "policy.rego")
	err = Output(outputFile, policy, policyBase64)
	assert.NoError(t, err)

	// Verify the plain policy file
	content, err := os.ReadFile(outputFile)
	assert.NoError(t, err)
	assert.Equal(t, policy, string(content))

	// Verify the companion base64 file
	b64File := filepath.Join(tmpDir, "policy_base64")
	b64Content, err := os.ReadFile(b64File)
	assert.NoError(t, err)
	assert.Equal(t, policyBase64, string(b64Content))

	// Verify policy contains expected content
	policyStr := string(content)
	assert.True(t, strings.Contains(policyStr, "package agent_policy"))
	assert.True(t, strings.Contains(policyStr, "allow_image"))
	assert.True(t, strings.Contains(policyStr, "allow_command"))
	assert.True(t, strings.Contains(policyStr, `docker\\.io/library/nginx:latest`))
	assert.True(t, strings.Contains(policyStr, `docker\\.io/library/redis:alpine`))
	assert.True(t, strings.Contains(policyStr, `docker\\.io/library/postgres:15-alpine`))
	assert.True(t, strings.Contains(policyStr, `quay\\.io/prometheus/busybox:latest`))
}
