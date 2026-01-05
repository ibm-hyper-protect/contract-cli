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
	"fmt"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "validate-contract"
	ParameterShortDescription = "Validate contract schema"
	ParameterLongDescription  = `Validate an unencrypted contract against the Hyper Protect schema.

Checks contract structure, required fields, and data types before encryption.
Helps catch errors early in the development process.`
	InputFlagDescription     = "Path to unencrypted Hyper Protect contract YAML file"
	InputFlagName            = "in"
	OsVersionFlagName        = "os"
	OsVersionFlagDescription = "Target Hyper Protect platform (hpvs, hpcr-rhvs, or hpcc-peerpod)"
)

// ValidateInput - function to validate plain contract
func ValidateInput(cmd *cobra.Command) (string, string, error) {
	contract, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", err
	}
	if contract == "" {
		err := fmt.Errorf("Error: required flag '--in' is missing")
		common.SetMandatoryFlagError(cmd, err)
	}

	version, err := cmd.Flags().GetString(OsVersionFlagName)
	if err != nil {
		return "", "", err
	}

	return contract, version, nil
}
