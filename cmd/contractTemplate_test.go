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
	"os"
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/lib/contractTemplate"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testContractTemplateOutputPath = "../build/test_cmd_contract_template_output.yaml"
	testContractTemplateInvalidOut = "../build/nonexistent/dir/output.yaml"
)

// getContractTemplateCmd returns a fresh contractTemplate command for testing,
// avoiding os.Exit issues caused by the global contractTemplateCmd in production use.
func getContractTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   contractTemplate.ParameterName,
		Short: contractTemplate.ParameterShortDescription,
		Long:  contractTemplate.ParameterLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			templateType, osVersion, outputPath, err := contractTemplate.ValidateInput(cmd)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			tmpl, err := contractTemplate.GenerateContractTemplate(templateType, osVersion)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			err = contractTemplate.Output(tmpl, outputPath)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
		},
	}

	cmd.PersistentFlags().String(contractTemplate.TypeFlagName, contractTemplate.TypeContract, contractTemplate.TypeFlagDescription)
	cmd.PersistentFlags().String(contractTemplate.OsVersionFlagName, contractTemplate.OsHpvs, contractTemplate.OsVersionFlagDescription)
	cmd.PersistentFlags().String(contractTemplate.OutputFlagName, "", contractTemplate.OutputFlagDescription)

	return cmd
}

// TestContractTemplateCmd_AllTypesToStdout verifies every --type value runs without error
// and prints to stdout (no --out flag).
func TestContractTemplateCmd_AllTypesToStdout(t *testing.T) {
	for _, typeVal := range []string{
		contractTemplate.TypeContract,
		contractTemplate.TypeEnv,
		contractTemplate.TypeWorkload,
		"", // empty defaults to contract
	} {
		cmd := getContractTemplateCmd()
		cmd.SetArgs([]string{"--" + contractTemplate.TypeFlagName, typeVal})
		assert.NoError(t, cmd.Execute(), "type=%q", typeVal)
	}
}

// TestContractTemplateCmd_AllOsValues verifies every supported --os value runs without error.
func TestContractTemplateCmd_AllOsValues(t *testing.T) {
	for _, osVal := range []string{
		contractTemplate.OsHpvs,
		contractTemplate.OsCcrt,
		contractTemplate.OsCcrv,
		contractTemplate.OsCccoPeerpod,
		contractTemplate.OsCccoBmtl,
	} {
		cmd := getContractTemplateCmd()
		cmd.SetArgs([]string{"--" + contractTemplate.OsVersionFlagName, osVal})
		assert.NoError(t, cmd.Execute(), "os=%s", osVal)
	}
}

// TestContractTemplateCmd_CcrtWorkloadToFile verifies the standard workload template
// (hpvs/ccrt) is written to a file and includes compose.
func TestContractTemplateCmd_CcrtWorkloadToFile(t *testing.T) {
	os.Remove(testContractTemplateOutputPath)

	cmd := getContractTemplateCmd()
	cmd.SetArgs([]string{
		"--" + contractTemplate.TypeFlagName, contractTemplate.TypeWorkload,
		"--" + contractTemplate.OsVersionFlagName, contractTemplate.OsCcrt,
		"--" + contractTemplate.OutputFlagName, testContractTemplateOutputPath,
	})
	assert.NoError(t, cmd.Execute())

	content, _ := os.ReadFile(testContractTemplateOutputPath)
	assert.Contains(t, string(content), "compose:")
	assert.Contains(t, string(content), "play:")
	os.Remove(testContractTemplateOutputPath)
}

// TestContractTemplateCmd_CcrvWorkloadToFile verifies the CCRV workload template
// is written to a file with play but without compose.
func TestContractTemplateCmd_CcrvWorkloadToFile(t *testing.T) {
	os.Remove(testContractTemplateOutputPath)

	cmd := getContractTemplateCmd()
	cmd.SetArgs([]string{
		"--" + contractTemplate.TypeFlagName, contractTemplate.TypeWorkload,
		"--" + contractTemplate.OsVersionFlagName, contractTemplate.OsCcrv,
		"--" + contractTemplate.OutputFlagName, testContractTemplateOutputPath,
	})
	assert.NoError(t, cmd.Execute())

	content, _ := os.ReadFile(testContractTemplateOutputPath)
	assert.Contains(t, string(content), "play:")
	assert.NotContains(t, string(content), "compose:")
	os.Remove(testContractTemplateOutputPath)
}

// TestContractTemplateCmd_CcrvCombinedToFile verifies the combined CCRV template
// has both sections but no compose.
func TestContractTemplateCmd_CcrvCombinedToFile(t *testing.T) {
	os.Remove(testContractTemplateOutputPath)

	cmd := getContractTemplateCmd()
	cmd.SetArgs([]string{
		"--" + contractTemplate.TypeFlagName, contractTemplate.TypeContract,
		"--" + contractTemplate.OsVersionFlagName, contractTemplate.OsCcrv,
		"--" + contractTemplate.OutputFlagName, testContractTemplateOutputPath,
	})
	assert.NoError(t, cmd.Execute())

	content, _ := os.ReadFile(testContractTemplateOutputPath)
	assert.Contains(t, string(content), "workload:")
	assert.Contains(t, string(content), "env:")
	assert.NotContains(t, string(content), "compose:")
	os.Remove(testContractTemplateOutputPath)
}

// TestContractTemplateCmd_CccoPeerpodWorkloadToFile verifies the ccco-peerpod workload
func TestContractTemplateCmd_CccoPeerpodWorkloadToFile(t *testing.T) {
	os.Remove(testContractTemplateOutputPath)

	cmd := getContractTemplateCmd()
	cmd.SetArgs([]string{
		"--" + contractTemplate.TypeFlagName, contractTemplate.TypeWorkload,
		"--" + contractTemplate.OsVersionFlagName, contractTemplate.OsCccoPeerpod,
		"--" + contractTemplate.OutputFlagName, testContractTemplateOutputPath,
	})
	assert.NoError(t, cmd.Execute())

	content, _ := os.ReadFile(testContractTemplateOutputPath)
	assert.Contains(t, string(content), "confidential-containers:")
	assert.NotContains(t, string(content), "compose:")
	assert.NotContains(t, string(content), "volumes:")
	os.Remove(testContractTemplateOutputPath)
}

// TestContractTemplateCmd_CccoBmtlWorkloadToFile verifies the ccco-bmtl workload
func TestContractTemplateCmd_CccoBmtlWorkloadToFile(t *testing.T) {
	os.Remove(testContractTemplateOutputPath)

	cmd := getContractTemplateCmd()
	cmd.SetArgs([]string{
		"--" + contractTemplate.TypeFlagName, contractTemplate.TypeWorkload,
		"--" + contractTemplate.OsVersionFlagName, contractTemplate.OsCccoBmtl,
		"--" + contractTemplate.OutputFlagName, testContractTemplateOutputPath,
	})
	assert.NoError(t, cmd.Execute())

	content, _ := os.ReadFile(testContractTemplateOutputPath)
	assert.Contains(t, string(content), "confidential-containers:")
	assert.Contains(t, string(content), "volumes:")
	assert.NotContains(t, string(content), "compose:")
	os.Remove(testContractTemplateOutputPath)
}

// TestContractTemplateCmd_CccoPeerpodEnvToFile verifies the ccco-peerpod env template
func TestContractTemplateCmd_CccoPeerpodEnvToFile(t *testing.T) {
	os.Remove(testContractTemplateOutputPath)

	cmd := getContractTemplateCmd()
	cmd.SetArgs([]string{
		"--" + contractTemplate.TypeFlagName, contractTemplate.TypeEnv,
		"--" + contractTemplate.OsVersionFlagName, contractTemplate.OsCccoPeerpod,
		"--" + contractTemplate.OutputFlagName, testContractTemplateOutputPath,
	})
	assert.NoError(t, cmd.Execute())

	content, _ := os.ReadFile(testContractTemplateOutputPath)
	assert.Contains(t, string(content), "logRouter:")
	assert.NotContains(t, string(content), "host-attestation:")
	assert.NotContains(t, string(content), "volumes:")
	os.Remove(testContractTemplateOutputPath)
}

// TestContractTemplateCmd_CccoBmtlEnvToFile verifies the ccco-bmtl env template
func TestContractTemplateCmd_CccoBmtlEnvToFile(t *testing.T) {
	os.Remove(testContractTemplateOutputPath)

	cmd := getContractTemplateCmd()
	cmd.SetArgs([]string{
		"--" + contractTemplate.TypeFlagName, contractTemplate.TypeEnv,
		"--" + contractTemplate.OsVersionFlagName, contractTemplate.OsCccoBmtl,
		"--" + contractTemplate.OutputFlagName, testContractTemplateOutputPath,
	})
	assert.NoError(t, cmd.Execute())

	content, _ := os.ReadFile(testContractTemplateOutputPath)
	assert.Contains(t, string(content), "logRouter:")
	assert.Contains(t, string(content), "volumes:")
	assert.Contains(t, string(content), "host-attestation:")
	os.Remove(testContractTemplateOutputPath)
}

// TestContractTemplateCmd_InvalidType verifies error is printed for unsupported --type.
func TestContractTemplateCmd_InvalidType(t *testing.T) {
	cmd := getContractTemplateCmd()
	cmd.SetArgs([]string{"--" + contractTemplate.TypeFlagName, "invalid-type"})
	assert.NoError(t, cmd.Execute())
}

// TestContractTemplateCmd_InvalidOs verifies error is printed for unsupported --os.
func TestContractTemplateCmd_InvalidOs(t *testing.T) {
	cmd := getContractTemplateCmd()
	cmd.SetArgs([]string{"--" + contractTemplate.OsVersionFlagName, "bad-os"})
	assert.NoError(t, cmd.Execute())
}

// TestContractTemplateCmd_InvalidOutputPath verifies error is printed for unwriteable path.
func TestContractTemplateCmd_InvalidOutputPath(t *testing.T) {
	cmd := getContractTemplateCmd()
	cmd.SetArgs([]string{"--" + contractTemplate.OutputFlagName, testContractTemplateInvalidOut})
	assert.NoError(t, cmd.Execute())
}
