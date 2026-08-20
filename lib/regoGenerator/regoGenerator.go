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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/ibm-hyper-protect/contract-go/v2/rego"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

const (
	ParameterName             = "rego-generator"
	ParameterShortDescription = "Generate OPA Rego policy from Kubernetes pod YAML"
	ParameterLongDescription  = `Generate Open Policy Agent (OPA) Rego policy from a Kubernetes pod specification.

This command reads a Kubernetes pod YAML file and generates an OPA Rego policy
that includes allow_image() and allow_command() rules for validating container
images and commands in IBM Confidential Computing environments.

The generated policy can be used with the Kata Agent Policy to enforce security
policies for confidential containers.`
	InputFlagName         = "in"
	InputFlagDescription  = "Path to Kubernetes pod YAML file (use '-' for standard input)"
	OutputFlagName        = "out"
	OutputFlagDescription = "Path to save generated Rego policy (optional, prints to stdout if not specified)"
	successMessage        = "Successfully generated Rego policy"
	successMessageBase64  = "Successfully generated Rego policy (base64)"
)

// ValidateInput validates the input flags for rego-generator command
func ValidateInput(cmd *cobra.Command) (string, string, error) {
	inputPath, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", err
	}
	if inputPath == "" {
		err := fmt.Errorf("Error: required flag '--in' is missing")
		common.SetMandatoryFlagError(cmd, err)
		return "", "", err
	}

	// Validate stdin input
	common.ValidateStdinInput(cmd, inputPath)

	outputPath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", err
	}

	return inputPath, outputPath, nil
}

// Process reads the Kubernetes resource YAML and generates the Rego policy.
// Returns the plain policy string and its IBM CC base64 representation via contract.HpcrText.
func Process(inputPath string) (policy, policyBase64 string, err error) {
	var resourceYAML string

	// Handle stdin input
	if inputPath == "-" {
		resourceYAML, err = common.ReadDataFromStdin()
		if err != nil {
			return "", "", fmt.Errorf("failed to read from stdin: %w", err)
		}
	} else {
		// Read from file
		data, err := os.ReadFile(inputPath)
		if err != nil {
			return "", "", fmt.Errorf("failed to read YAML file: %w", err)
		}
		resourceYAML = string(data)
	}

	// Extract Pod spec from the resource
	podYAML, err := extractPodSpec(resourceYAML)
	if err != nil {
		return "", "", fmt.Errorf("failed to extract pod spec: %w", err)
	}

	// Generate Rego policy using the contract-go library
	policy, _, _, err = rego.GenerateRegoPolicy(podYAML, "")
	if err != nil {
		return "", "", fmt.Errorf("failed to generate Rego policy: %w", err)
	}

	// Wrap policy in IBM CC base64 format
	policyBase64, _, _, err = contract.HpcrText(policy)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate base64 policy: %w", err)
	}

	return policy, policyBase64, nil
}

// extractPodSpec extracts the Pod spec from various Kubernetes resources
func extractPodSpec(resourceYAML string) (string, error) {
	// Parse the YAML to determine the kind
	var resource map[string]interface{}
	if err := yaml.Unmarshal([]byte(resourceYAML), &resource); err != nil {
		return "", fmt.Errorf("failed to parse YAML: %w", err)
	}

	kind, ok := resource["kind"].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid 'kind' field in YAML")
	}

	// Handle different resource types
	switch strings.ToLower(kind) {
	case "pod":
		// Already a Pod, return as-is
		return resourceYAML, nil

	case "deployment", "statefulset", "daemonset":
		// Extract spec.template for these resources
		spec, ok := resource["spec"].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("missing or invalid 'spec' field in %s", kind)
		}

		template, ok := spec["template"].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("missing or invalid 'spec.template' field in %s", kind)
		}

		// Create a Pod from the template
		pod := map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata":   template["metadata"],
			"spec":       template["spec"],
		}

		// Convert back to YAML
		podYAML, err := yaml.Marshal(pod)
		if err != nil {
			return "", fmt.Errorf("failed to marshal pod spec: %w", err)
		}

		return string(podYAML), nil

	case "cronjob":
		// Extract spec.jobTemplate.spec.template for CronJob
		spec, ok := resource["spec"].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("missing or invalid 'spec' field in CronJob")
		}

		jobTemplate, ok := spec["jobTemplate"].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("missing or invalid 'spec.jobTemplate' field in CronJob")
		}

		jobSpec, ok := jobTemplate["spec"].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("missing or invalid 'spec.jobTemplate.spec' field in CronJob")
		}

		template, ok := jobSpec["template"].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("missing or invalid 'spec.jobTemplate.spec.template' field in CronJob")
		}

		// Create a Pod from the template
		pod := map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata":   template["metadata"],
			"spec":       template["spec"],
		}

		// Convert back to YAML
		podYAML, err := yaml.Marshal(pod)
		if err != nil {
			return "", fmt.Errorf("failed to marshal pod spec: %w", err)
		}

		return string(podYAML), nil

	default:
		return "", fmt.Errorf("unsupported resource kind: %s (supported: Pod, Deployment, StatefulSet, DaemonSet, CronJob)", kind)
	}
}

// base64OutputPath derives the companion base64 output path from the policy output path.
// e.g. "policy.rego" -> "policy_base64", "out/policy.rego" -> "out/policy_base64"
func base64OutputPath(outputPath string) string {
	dir := filepath.Dir(outputPath)
	base := filepath.Base(outputPath)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	return filepath.Join(dir, stem+"_base64")
}

// Output writes the plain policy to outputPath (or stdout if empty).
// When outputPath is given, also writes the IBM CC base64 policy to <stem>_base64.
func Output(outputPath, policy, policyBase64 string) error {
	if outputPath == "" {
		// Print plain policy to stdout only
		fmt.Println(policy)
		return nil
	}

	// Write plain policy
	if err := os.WriteFile(outputPath, []byte(policy), 0644); err != nil {
		return fmt.Errorf("failed to write Rego policy to file: %w", err)
	}
	fmt.Printf("%s: %s\n", successMessage, outputPath)

	// Write IBM CC base64 policy to companion file
	b64Path := base64OutputPath(outputPath)
	if err := os.WriteFile(b64Path, []byte(policyBase64), 0644); err != nil {
		return fmt.Errorf("failed to write base64 Rego policy to file: %w", err)
	}
	fmt.Printf("%s: %s\n", successMessageBase64, b64Path)

	return nil
}
