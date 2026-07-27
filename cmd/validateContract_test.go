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
	"fmt"
	"testing"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-cli/lib/validateContract"
	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	sampleValidateContractInput    = "../samples/contract.yaml"
	sampleValidateContractEnvInput = "../samples/env.yaml"
	sampleValidateContractWlInput  = "../samples/workload.yaml"
	sampleValidateContractOsType   = "ccrt"
	sampleValidateContractSection  = ""
)

var (
	sampleValidContractCommand         = []string{validateContract.ParameterName, "--in", sampleValidateContractInput, "--os", sampleValidateContractOsType}
	sampleValidContractWorkloadCommand = []string{validateContract.ParameterName, "--in", sampleValidateContractWlInput, "--os", sampleValidateContractOsType, "--type", "workload"}
	sampleValidContractEnvCommand      = []string{validateContract.ParameterName, "--in", sampleValidateContractEnvInput, "--os", sampleValidateContractOsType, "--type", "env"}
)

// getValidateContractCmd returns a fresh validate-contract command instance for isolated testing.
// Using a fresh command avoids shared-state issues with rootCmd across test cases.
//
// Note: error cases that trigger log.Fatal (os.Exit) cannot be tested at the cmd level.
// Those scenarios are tested at the lib level via HpcrVerifyContract directly.
func getValidateContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   validateContract.ParameterName,
		Short: validateContract.ParameterShortDescription,
		Long:  validateContract.ParameterLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			contractPath, version, section, err := validateContract.ValidateInput(cmd)
			if err != nil {
				return err
			}
			if !common.CheckFileFolderExists(contractPath) {
				return fmt.Errorf("the path to contract doesn't exist")
			}
			contractData, err := common.ReadDataFromFile(contractPath)
			if err != nil {
				return err
			}
			return contract.HpcrVerifyContract(contractData, version, section)
		},
	}
	cmd.Flags().String(validateContract.InputFlagName, "", validateContract.InputFlagDescription)
	cmd.Flags().String(validateContract.OsVersionFlagName, "", validateContract.OsVersionFlagDescription)
	cmd.Flags().String(validateContract.TypeSectionFlagName, "", validateContract.TypeSectionFlagDescription)
	return cmd
}

// TestValidateContractCmdSucess validates the full contract (both workload and env sections)
func TestValidateContractCmdSucess(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs(sampleValidContractCommand)
	err := validateContractCmd.Execute()

	assert.NoError(t, err)
}

// TestValidateContractCmdWorkloadSection validates only the workload section
func TestValidateContractCmdWorkloadSection(t *testing.T) {
	cmd := getValidateContractCmd()
	cmd.SetArgs([]string{"--in", sampleValidateContractWlInput, "--os", sampleValidateContractOsType, "--type", "workload"})

	err := cmd.Execute()

	assert.NoError(t, err)
}

// TestValidateContractCmdEnvSection validates only the env section
func TestValidateContractCmdEnvSection(t *testing.T) {
	cmd := getValidateContractCmd()
	cmd.SetArgs([]string{"--in", sampleValidateContractEnvInput, "--os", sampleValidateContractOsType, "--type", "env"})

	err := cmd.Execute()

	assert.NoError(t, err)
}

// TestValidateContractCmdAllOsVersions validates the full contract across OS versions that
// accept the compose-based sample (ccrt and ccco). ccrv requires play or confidential-containers.
func TestValidateContractCmdAllOsVersions(t *testing.T) {
	for _, osVer := range []string{"ccrt", "ccco"} {
		t.Run(osVer, func(t *testing.T) {
			cmd := getValidateContractCmd()
			cmd.SetArgs([]string{"--in", sampleValidateContractInput, "--os", osVer})

			err := cmd.Execute()

			assert.NoError(t, err)
		})
	}
}

// TestValidateContractCmdWorkloadAllOsVersions validates the workload section across OS versions that
// accept the compose-based sample (ccrt and ccco). ccrv requires play or confidential-containers.
func TestValidateContractCmdWorkloadAllOsVersions(t *testing.T) {
	for _, osVer := range []string{"ccrt", "ccco"} {
		t.Run(osVer, func(t *testing.T) {
			cmd := getValidateContractCmd()
			cmd.SetArgs([]string{"--in", sampleValidateContractWlInput, "--os", osVer, "--type", "workload"})

			err := cmd.Execute()

			assert.NoError(t, err)
		})
	}
}

// TestValidateContractCmdEnvAllOsVersions validates the env section across all supported OS versions
func TestValidateContractCmdEnvAllOsVersions(t *testing.T) {
	for _, osVer := range []string{"ccrt", "ccrv", "ccco"} {
		t.Run(osVer, func(t *testing.T) {
			cmd := getValidateContractCmd()
			cmd.SetArgs([]string{"--in", sampleValidateContractEnvInput, "--os", osVer, "--type", "env"})

			err := cmd.Execute()

			assert.NoError(t, err)
		})
	}
}

// TestValidateContractCmdWrapperRejectedForTypeEnv checks that a wrapper-format contract
// is rejected when --type env is specified (must use raw section format)
func TestValidateContractCmdWrapperRejectedForTypeEnv(t *testing.T) {
	cmd := getValidateContractCmd()
	cmd.SetArgs([]string{"--in", sampleValidateContractInput, "--os", sampleValidateContractOsType, "--type", "env"})

	err := cmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contract contains both env and workload sections")
}

// TestValidateContractCmdWrapperRejectedForTypeWorkload checks that a wrapper-format contract
// is rejected when --type workload is specified (must use raw section format)
func TestValidateContractCmdWrapperRejectedForTypeWorkload(t *testing.T) {
	cmd := getValidateContractCmd()
	cmd.SetArgs([]string{"--in", sampleValidateContractInput, "--os", sampleValidateContractOsType, "--type", "workload"})

	err := cmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contract contains both env and workload sections")
}
