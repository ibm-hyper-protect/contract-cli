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

package signContract

import (
	"fmt"
	"strings"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "sign-contract"
	ParameterShortDescription = "Sign an encrypted contract"
	ParameterLongDescription  = `Sign a contract with encrypted workload and env`
	InputFlagName             = "in"
	InputFlagDescription      = "Path to encrypted contract (mutually exclusive with --env/--workload)"
	EnvFlagName               = "env"
	EnvFlagDescription        = "Path to file containing the encrypted env section (must be used together with --workload)"
	WorkloadFlagName          = "workload"
	WorkloadFlagDescription   = "Path to file containing the encrypted workload section (must be used together with --env)"
	PrivateKeyFlagName        = "priv"
	PrivateKeyFlagDescription = "Path to private key file for signing"
	PasswordFlagName          = "password"
	PasswordFlagDescription   = "Password for encrypted private key"
	OutputFlagName            = "out"
	OutputFlagDescription     = "Path to save signed output"
)

// ValidateInput validates inputs of sign-contract and returns:
// inputPath, privateKeyPath, outputPath, password, envPath, workloadPath, error
func ValidateInput(cmd *cobra.Command) (string, string, string, string, string, string, error) {
	inputData, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", "", "", "", "", err
	}

	privateKeyPath, err := cmd.Flags().GetString(PrivateKeyFlagName)
	if err != nil {
		return "", "", "", "", "", "", err
	}

	envPath, err := cmd.Flags().GetString(EnvFlagName)
	if err != nil {
		return "", "", "", "", "", "", err
	}

	workloadPath, err := cmd.Flags().GetString(WorkloadFlagName)
	if err != nil {
		return "", "", "", "", "", "", err
	}

	// Mutual-exclusion: --in cannot be combined with --env or --workload
	if inputData != "" && (envPath != "" || workloadPath != "") {
		err := fmt.Errorf("Error: '--in' is mutually exclusive with '--env'/'--workload'. Provide either '--in' or both '--env' and '--workload'")
		common.SetMandatoryFlagError(cmd, err)
	}

	// --env and --workload must be a pair
	if (envPath != "" && workloadPath == "") || (envPath == "" && workloadPath != "") {
		err := fmt.Errorf("Error: '--env' and '--workload' must be provided together")
		common.SetMandatoryFlagError(cmd, err)
	}

	// At least one input path must be given
	if inputData == "" && envPath == "" && workloadPath == "" {
		err := fmt.Errorf("Error: required flag(s) missing — provide '--in' or both '--env' and '--workload'")
		common.SetMandatoryFlagError(cmd, err)
	}

	// --priv is always required
	if privateKeyPath == "" {
		err := fmt.Errorf("Error: required flag '--priv' is missing")
		common.SetMandatoryFlagError(cmd, err)
	}

	// Stdin conflict check applies only to the --in path
	if inputData != "" {
		common.ValidateStdinInput(cmd, inputData)
	}

	outputPath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", "", "", "", "", err
	}

	password, err := cmd.Flags().GetString(PasswordFlagName)
	if err != nil {
		return "", "", "", "", "", "", err
	}

	return inputData, privateKeyPath, outputPath, password, envPath, workloadPath, nil
}

// BuildContractFromParts reads the encrypted env and workload files and assembles
// an in-memory contract YAML string equivalent to a pre-assembled encrypted-contract.yaml.
// The result is never written to disk.
func BuildContractFromParts(envPath, workloadPath string) (string, error) {
	if !common.CheckFileFolderExists(envPath) {
		return "", fmt.Errorf("env file path doesn't exist: %s", envPath)
	}

	if !common.CheckFileFolderExists(workloadPath) {
		return "", fmt.Errorf("workload file path doesn't exist: %s", workloadPath)
	}

	envBlob, err := common.ReadDataFromFile(envPath)
	if err != nil {
		return "", fmt.Errorf("failed to read env file: %w", err)
	}

	workloadBlob, err := common.ReadDataFromFile(workloadPath)
	if err != nil {
		return "", fmt.Errorf("failed to read workload file: %w", err)
	}

	envBlob = strings.TrimSpace(envBlob)
	workloadBlob = strings.TrimSpace(workloadBlob)

	contractYAML := fmt.Sprintf("env: %s\nworkload: %s\n", envBlob, workloadBlob)
	return contractYAML, nil
}

// GenerateSignContract signs a contract. It accepts either:
//   - inputDataPath: path to a pre-assembled encrypted-contract.yaml (or "-" for stdin)
//   - envPath + workloadPath: paths to individual encrypted section files, stitched in-memory
func GenerateSignContract(inputDataPath, privateKeyPath, password, envPath, workloadPath string) (string, error) {
	var inputData string
	var err error

	if inputDataPath != "" {
		// --in path: read from file or stdin
		if inputDataPath == "-" {
			inputData, err = common.ReadDataFromStdin()
			if err != nil {
				return "", fmt.Errorf("unable to read input from standard input: %w", err)
			}
		} else {
			if !common.CheckFileFolderExists(inputDataPath) {
				return "", fmt.Errorf("the contract path doesn't exist")
			}
			inputData, err = common.ReadDataFromFile(inputDataPath)
			if err != nil {
				return "", err
			}
		}
	} else {
		// --env + --workload path: stitch in-memory
		inputData, err = BuildContractFromParts(envPath, workloadPath)
		if err != nil {
			return "", err
		}
	}

	privateKey, err := common.GetPrivateKey(privateKeyPath)
	if err != nil {
		return "", err
	}

	signedContract, _, _, err := contract.HpcrContractSign(inputData, privateKey, password)
	if err != nil {
		return "", err
	}

	return signedContract, nil
}

// Output prints the signed contract to stdout or writes it to a file.
func Output(signedContract, outputPath string) error {
	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, signedContract)
		if err != nil {
			return err
		}
		fmt.Println("Successfully generated signed contract")
	} else {
		fmt.Println(signedContract)
	}

	return nil
}
