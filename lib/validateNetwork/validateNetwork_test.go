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

package validateNetwork

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testNetworkConfigPath = "../../samples/network/network-config.yaml"
)

// TestValidateInput_Success tests ValidateInput with valid input
func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, testNetworkConfigPath, "")

	networkConfig, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, testNetworkConfigPath, networkConfig)
}

// TestValidateInput_FlagErrors tests error handling for flag retrieval
func TestValidateInput_FlagErrors(t *testing.T) {
	cmd := &cobra.Command{}
	_, err := ValidateInput(cmd)
	assert.Error(t, err)
}

// TestValidateInput_WithParentDirectory tests ValidateInput with parent directory reference
func TestValidateInput_WithParentDirectory(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, "../samples/network/network-config.yaml", "")

	networkConfig, err := ValidateInput(cmd)

	assert.NoError(t, err)
	assert.Equal(t, "../samples/network/network-config.yaml", networkConfig)
}
