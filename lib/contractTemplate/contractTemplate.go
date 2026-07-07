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

package contractTemplate

import (
	"fmt"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "contract-template"
	ParameterShortDescription = "Generate a contract template"
	ParameterLongDescription  = `Generate a contract YAML template for IBM Confidential Computing deployments.

Returns a pre-filled YAML scaffold for the workload section, env section, or a
combined contract containing both sections. Use this as a starting point when
authoring a new contract.`

	TypeFlagName        = "type"
	TypeFlagDescription = "Template type to generate: env, workload, or contract (default: contract)"

	OsVersionFlagName        = "os"
	OsVersionFlagDescription = "Target IBM Confidential Computing platform (hpvs, ccrt, ccrv, or ccco) (default: hpvs)"

	OutputFlagName        = "out"
	OutputFlagDescription = "Path to save the generated template (prints to terminal if not specified)"

	// ValidTypes lists all accepted --type values.
	TypeEnv      = "env"
	TypeWorkload = "workload"
	TypeContract = "contract"

	// ValidOs lists all accepted --os values.
	OsHpvs = "hpvs"
	OsCcrt = "ccrt"
	OsCcrv = "ccrv"
	OsCcco = "ccco"
)

// ValidateInput parses and validates flags from the cobra command.
// Returns templateType, osVersion, outputPath, error.
func ValidateInput(cmd *cobra.Command) (string, string, string, error) {
	templateType, err := cmd.Flags().GetString(TypeFlagName)
	if err != nil {
		return "", "", "", err
	}

	// Validate --type value
	switch templateType {
	case TypeEnv, TypeWorkload, TypeContract, "":
		// empty string treated as default (contract)
	default:
		return "", "", "", fmt.Errorf("invalid value for --type: %q. Allowed values: env, workload, contract", templateType)
	}

	osVersion, err := cmd.Flags().GetString(OsVersionFlagName)
	if err != nil {
		return "", "", "", err
	}

	// Validate --os value if provided
	if osVersion != "" {
		switch osVersion {
		case OsHpvs, OsCcrt, OsCcrv, OsCcco:
			// valid
		default:
			return "", "", "", fmt.Errorf("invalid value for --os: %q. Allowed values: hpvs, ccrt, ccrv, ccco", osVersion)
		}
	}

	outputPath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", "", err
	}

	return templateType, osVersion, outputPath, nil
}

// GenerateContractTemplate calls contract-go to produce the requested template.
// templateType may be "env", "workload", or "" / "contract" (combined).
// os selects the platform-specific workload template ("ccrv" differs from all others).
func GenerateContractTemplate(templateType, os string) (string, error) {
	// Map "contract" and "" to "" (combined template in contract-go)
	contractGoType := templateType
	if templateType == TypeContract || templateType == "" {
		contractGoType = ""
	}

	result, err := contract.HpcrContractTemplate(contractGoType, os)
	if err != nil {
		return "", err
	}

	return result, nil
}

// Output prints the template to stdout or writes it to a file.
func Output(template, outputPath string) error {
	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, template)
		if err != nil {
			return err
		}
		fmt.Println("Successfully generated contract template")
	} else {
		fmt.Print(template)
	}

	return nil
}
