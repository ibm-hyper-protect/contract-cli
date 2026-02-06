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

package validateContract

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testContractPath = "../../samples/contract.yaml"
	testOsVersion    = "hpvs"
)

// TestValidateInput_Success tests ValidateInput with all required flags
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(OsVersionFlagName, testOsVersion, "")

	contract, version, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, contract)
	assert.Equal(t, testOsVersion, version)
}

// TestValidateInput_WithHpcrRhvs tests ValidateInput with hpcr-rhvs OS version
func TestValidateInput_WithHpcrRhvs(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(OsVersionFlagName, "hpcr-rhvs", "")

	contract, version, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, contract)
	assert.Equal(t, "hpcr-rhvs", version)
}

// TestValidateInput_WithHpccPeerpod tests ValidateInput with hpcc-peerpod OS version
func TestValidateInput_WithHpccPeerpod(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(OsVersionFlagName, "hpcc-peerpod", "")

	contract, version, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, contract)
	assert.Equal(t, "hpcc-peerpod", version)
}

// TestValidateInput_WithoutOsVersion tests ValidateInput without OS version (optional)
func TestValidateInput_WithoutOsVersion(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testContractPath, "")
	cmd.Flags().String(OsVersionFlagName, "", "")

	contract, version, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testContractPath, contract)
	assert.Equal(t, "", version)
}

// TestValidateInput_FlagErrors tests error handling for flag retrieval
func TestValidateInput_FlagErrors(t *testing.T) {
	cmd := &cobra.Command{}
	_, _, err := ValidateInput(cmd)
	assert.Error(t, err)
}

// TestValidateInput_WithRelativePath tests ValidateInput with relative path
func TestValidateInput_WithRelativePath(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, "./contract.yaml", "")
	cmd.Flags().String(OsVersionFlagName, testOsVersion, "")

	contract, version, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, "./contract.yaml", contract)
	assert.Equal(t, testOsVersion, version)
}

// TestValidateInput_AllOsVersions tests ValidateInput with all supported OS versions
func TestValidateInput_AllOsVersions(t *testing.T) {
	osVersions := []string{"hpvs", "hpcr-rhvs", "hpcc-peerpod"}

	for _, osVer := range osVersions {
		cmd := &cobra.Command{}
		cmd.Flags().String(InputFlagName, testContractPath, "")
		cmd.Flags().String(OsVersionFlagName, osVer, "")

		contract, version, err := ValidateInput(cmd)

		assert.NoError(t, err)
		assert.Equal(t, testContractPath, contract)
		assert.Equal(t, osVer, version)
	}
}
