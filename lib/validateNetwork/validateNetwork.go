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
	"fmt"

	"github.com/spf13/cobra"
)

const (
	ParameterName             = "validate-network"
	ParameterShortDescription = "Validate network configuration schema"
	ParameterLongDescription  = `Validate network-config YAML file against the schema.

Validates network configuration for on-premise deployments, ensuring all required
fields are present and properly formatted.`
	InputFlagDescription = "Path to network-config YAML file"
	InputFlagName        = "in"
)

// ValidateInput - function to get network-config file
func ValidateInput(cmd *cobra.Command) (string, error) {
	networkConfig, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", err
	}
	if networkConfig == "" {
		_ = cmd.Help()
		return "", fmt.Errorf("Error: required flag '--in' is missing.")
	}

	return networkConfig, nil
}
