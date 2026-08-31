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
	testPodFile = "../../samples/hpcc/sample-pod.yaml"
)

// TestValidateInput_Success tests ValidateInput with all required flags
func TestValidateInput_Success(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "policy.rego")

	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, "test-pod.yaml", "")
	cmd.Flags().String(OutputFlagName, outputPath, "")
	cmd.Flags().String(FormatFlagName, FormatBase64, "")

	inputPath, gotOutputPath, format, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, "test-pod.yaml", inputPath)
	assert.Equal(t, outputPath, gotOutputPath)
	assert.Equal(t, FormatBase64, format)
}

// TestValidateInput_WithoutOutputPath tests ValidateInput without output path (optional)
func TestValidateInput_WithoutOutputPath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, "test-pod.yaml", "")
	cmd.Flags().String(OutputFlagName, "", "")
	cmd.Flags().String(FormatFlagName, FormatBase64, "")

	inputPath, outputPath, format, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, "test-pod.yaml", inputPath)
	assert.Equal(t, "", outputPath)
	assert.Equal(t, FormatBase64, format)
}

// TestValidateInput_AllFormats tests ValidateInput accepts all valid format values
func TestValidateInput_AllFormats(t *testing.T) {
	for _, fmt := range []string{FormatBase64, FormatText, FormatBoth} {
		t.Run(fmt, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().String(InputFlagName, "test-pod.yaml", "")
			cmd.Flags().String(OutputFlagName, "", "")
			cmd.Flags().String(FormatFlagName, fmt, "")

			_, _, format, err := ValidateInput(cmd)

			assert.NoError(t, err)
			assert.Equal(t, fmt, format)
		})
	}
}

// TestValidateInput_InvalidFormat tests ValidateInput rejects unknown format values
func TestValidateInput_InvalidFormat(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, "test-pod.yaml", "")
	cmd.Flags().String(OutputFlagName, "", "")
	cmd.Flags().String(FormatFlagName, "json", "")

	_, _, _, err := ValidateInput(cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid --format value")
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

// TestOutput_Base64_ToStdout tests Output with format=base64 to stdout — only base64 printed
func TestOutput_Base64_ToStdout(t *testing.T) {
	policy := "package agent_policy\n\ndefault CreateContainerRequest := false"

	err := Output("", FormatBase64, policy, "hyper-protect-basic.abc123")

	assert.NoError(t, err)
}

// TestOutput_Text_ToStdout tests Output with format=text to stdout — only plain rego printed
func TestOutput_Text_ToStdout(t *testing.T) {
	policy := "package agent_policy\n\ndefault CreateContainerRequest := false"

	err := Output("", FormatText, policy, "hyper-protect-basic.abc123")

	assert.NoError(t, err)
}

// TestOutput_Both_ToStdout tests Output with format=both to stdout — both printed
func TestOutput_Both_ToStdout(t *testing.T) {
	policy := "package agent_policy\n\ndefault CreateContainerRequest := false"

	err := Output("", FormatBoth, policy, "hyper-protect-basic.abc123")

	assert.NoError(t, err)
}

// TestOutput_Base64_ToFile tests Output with format=base64 — only _base64 file written
func TestOutput_Base64_ToFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputStem := filepath.Join(tmpDir, "policy")
	policy := "package agent_policy\n\ndefault CreateContainerRequest := false"
	b64Policy := "hyper-protect-basic.abc123"

	err := Output(outputStem, FormatBase64, policy, b64Policy)

	assert.NoError(t, err)

	// Verify _base64 file is written
	b64File := filepath.Join(tmpDir, "policy_base64")
	b64Content, err := os.ReadFile(b64File)
	assert.NoError(t, err)
	assert.Equal(t, b64Policy, string(b64Content))

	// Verify plain rego file is NOT written
	_, err = os.Stat(filepath.Join(tmpDir, "policy.rego"))
	assert.True(t, os.IsNotExist(err), "plain .rego file must not be created for format=base64")
}

// TestOutput_Text_ToFile tests Output with format=text — only .rego file written
func TestOutput_Text_ToFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputStem := filepath.Join(tmpDir, "policy")
	policy := "package agent_policy\n\ndefault CreateContainerRequest := false"
	b64Policy := "hyper-protect-basic.abc123"

	err := Output(outputStem, FormatText, policy, b64Policy)

	assert.NoError(t, err)

	// Verify plain rego file is written
	regoFile := filepath.Join(tmpDir, "policy.rego")
	content, err := os.ReadFile(regoFile)
	assert.NoError(t, err)
	assert.Equal(t, policy, string(content))

	// Verify _base64 file is NOT written
	_, err = os.Stat(filepath.Join(tmpDir, "policy_base64"))
	assert.True(t, os.IsNotExist(err), "_base64 file must not be created for format=text")
}

// TestOutput_Both_ToFile tests Output with format=both — both files written
func TestOutput_Both_ToFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputStem := filepath.Join(tmpDir, "policy")
	policy := "package agent_policy\n\ndefault CreateContainerRequest := false"
	b64Policy := "hyper-protect-basic.abc123"

	err := Output(outputStem, FormatBoth, policy, b64Policy)

	assert.NoError(t, err)

	// Verify plain rego file
	regoFile := filepath.Join(tmpDir, "policy.rego")
	content, err := os.ReadFile(regoFile)
	assert.NoError(t, err)
	assert.Equal(t, policy, string(content))

	// Verify _base64 file
	b64File := filepath.Join(tmpDir, "policy_base64")
	b64Content, err := os.ReadFile(b64File)
	assert.NoError(t, err)
	assert.Equal(t, b64Policy, string(b64Content))
}

// TestOutput_Text_ToFile_WithRegoExtension tests that passing a stem already ending in .rego works
func TestOutput_Text_ToFile_WithRegoExtension(t *testing.T) {
	tmpDir := t.TempDir()
	outputStem := filepath.Join(tmpDir, "policy.rego")
	policy := "package agent_policy\n\ndefault CreateContainerRequest := false"

	err := Output(outputStem, FormatText, policy, "")

	assert.NoError(t, err)

	content, err := os.ReadFile(outputStem)
	assert.NoError(t, err)
	assert.Equal(t, policy, string(content))
}

// TestTextOutputPath tests the .rego path derivation
func TestTextOutputPath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"policy", "policy.rego"},
		{"out/policy", "out/policy.rego"},
		{"/abs/path/policy", "/abs/path/policy.rego"},
		{"policy.rego", "policy.rego"}, // already has extension — unchanged
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, textOutputPath(tt.input))
		})
	}
}

// TestBase64OutputPath tests the companion base64 path derivation
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

// TestOutput_InvalidPath tests Output with an invalid file path (format=base64 reaches the base64 write)
func TestOutput_InvalidPath(t *testing.T) {
	policy := "package agent_policy"

	err := Output("/invalid/path/that/does/not/exist/policy", FormatBase64, policy, "hyper-protect-basic.abc")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to write base64 Rego policy to file")
}

// TestOutput_InvalidPath_Text tests Output with an invalid file path (format=text reaches the rego write)
func TestOutput_InvalidPath_Text(t *testing.T) {
	policy := "package agent_policy"

	err := Output("/invalid/path/that/does/not/exist/policy", FormatText, policy, "hyper-protect-basic.abc")

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

// TestProcess_InvalidKind tests Process with an unsupported Kubernetes resource kind
func TestProcess_InvalidKind(t *testing.T) {
	invalidKindYAML := `apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
data:
  key: value
`
	tmpDir := t.TempDir()
	invalidFile := filepath.Join(tmpDir, "configmap.yaml")
	err := os.WriteFile(invalidFile, []byte(invalidKindYAML), 0644)
	assert.NoError(t, err)

	_, _, err = Process(invalidFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to extract pod spec")
	assert.Contains(t, err.Error(), "unsupported resource kind: ConfigMap")
}

// TestIntegration_EndToEnd tests the complete workflow with format=both (writes both files)
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

	// Output with format=both — should create policy.rego and policy_base64
	outputStem := filepath.Join(tmpDir, "policy")
	err = Output(outputStem, FormatBoth, policy, policyBase64)
	assert.NoError(t, err)

	// Verify the plain policy file
	regoFile := filepath.Join(tmpDir, "policy.rego")
	content, err := os.ReadFile(regoFile)
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
