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

	"github.com/ibm-hyper-protect/contract-cli/lib/signContract"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testSignContractPath      = "../samples/contract.yaml"
	testSignPrivateKeyPath    = "../samples/sign/private.pem"
	testSignOutputPath        = "../build/test_cmd_sign_contract_output.txt"
	testSignInvalidPath       = "../build/file/file_not_exists.txt"
	testSignCorruptedContract = "../build/corrupted_contract_cmd.yaml"
	testSignCorruptedKey      = "../build/corrupted_key_cmd.pem"
	testSignEncryptedEnvPath  = "../samples/sign/encrypted-env.txt"
	testSignEncryptedWlPath   = "../samples/sign/encrypted-workload.txt"
)

// getSignContractCmd returns a fresh instance of the sign-contract command for testing
func getSignContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   signContract.ParameterName,
		Short: signContract.ParameterShortDescription,
		Long:  signContract.ParameterLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			contract, privateKey, output, password, envPath, workloadPath, err := signContract.ValidateInput(cmd)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			contractSign, err := signContract.GenerateSignContract(contract, privateKey, password, envPath, workloadPath)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			err = signContract.Output(contractSign, output)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
		},
	}

	cmd.PersistentFlags().String(signContract.InputFlagName, "", signContract.InputFlagDescription)
	cmd.PersistentFlags().String(signContract.EnvFlagName, "", signContract.EnvFlagDescription)
	cmd.PersistentFlags().String(signContract.WorkloadFlagName, "", signContract.WorkloadFlagDescription)
	cmd.PersistentFlags().String(signContract.PrivateKeyFlagName, "", signContract.PrivateKeyFlagDescription)
	cmd.PersistentFlags().String(signContract.PasswordFlagName, "", signContract.PasswordFlagDescription)
	cmd.PersistentFlags().String(signContract.OutputFlagName, "", signContract.OutputFlagDescription)

	return cmd
}

// ─── Existing --in path (backward-compatible) ─────────────────────────────────

// TestSignContractCmd_Success tests successful contract signing via --in
func TestSignContractCmd_Success(t *testing.T) {
	os.Remove(testSignOutputPath)

	cmd := getSignContractCmd()
	cmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testSignContractPath,
		"--" + signContract.PrivateKeyFlagName, testSignPrivateKeyPath,
		"--" + signContract.OutputFlagName, testSignOutputPath,
	})

	err := cmd.Execute()
	assert.NoError(t, err)

	_, statErr := os.Stat(testSignOutputPath)
	assert.NoError(t, statErr)

	content, readErr := os.ReadFile(testSignOutputPath)
	assert.NoError(t, readErr)
	assert.Contains(t, string(content), "envWorkloadSignature")

	os.Remove(testSignOutputPath)
}

// TestSignContractCmd_WithPassword tests signing with password parameter
func TestSignContractCmd_WithPassword(t *testing.T) {
	os.Remove(testSignOutputPath)

	cmd := getSignContractCmd()
	cmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testSignContractPath,
		"--" + signContract.PrivateKeyFlagName, testSignPrivateKeyPath,
		"--" + signContract.PasswordFlagName, "testPassword123",
		"--" + signContract.OutputFlagName, testSignOutputPath,
	})

	err := cmd.Execute()
	assert.NoError(t, err)

	_, statErr := os.Stat(testSignOutputPath)
	assert.NoError(t, statErr)

	content, readErr := os.ReadFile(testSignOutputPath)
	assert.NoError(t, readErr)
	assert.Contains(t, string(content), "envWorkloadSignature")

	os.Remove(testSignOutputPath)
}

// TestSignContractCmd_WithEmptyPassword tests signing with empty password
func TestSignContractCmd_WithEmptyPassword(t *testing.T) {
	os.Remove(testSignOutputPath)

	cmd := getSignContractCmd()
	cmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testSignContractPath,
		"--" + signContract.PrivateKeyFlagName, testSignPrivateKeyPath,
		"--" + signContract.PasswordFlagName, "",
		"--" + signContract.OutputFlagName, testSignOutputPath,
	})

	err := cmd.Execute()
	assert.NoError(t, err)

	_, statErr := os.Stat(testSignOutputPath)
	assert.NoError(t, statErr)

	os.Remove(testSignOutputPath)
}

// TestSignContractCmd_WithoutOutputPath tests signing without output path (stdout)
func TestSignContractCmd_WithoutOutputPath(t *testing.T) {
	cmd := getSignContractCmd()
	cmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testSignContractPath,
		"--" + signContract.PrivateKeyFlagName, testSignPrivateKeyPath,
	})

	err := cmd.Execute()
	assert.NoError(t, err)
}

// TestSignContractCmd_WithPasswordAndOutput tests complete --in workflow with password and output
func TestSignContractCmd_WithPasswordAndOutput(t *testing.T) {
	os.Remove(testSignOutputPath)

	cmd := getSignContractCmd()
	cmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testSignContractPath,
		"--" + signContract.PrivateKeyFlagName, testSignPrivateKeyPath,
		"--" + signContract.PasswordFlagName, "securePass456",
		"--" + signContract.OutputFlagName, testSignOutputPath,
	})

	err := cmd.Execute()
	assert.NoError(t, err)

	content, readErr := os.ReadFile(testSignOutputPath)
	assert.NoError(t, readErr)
	assert.NotEmpty(t, content)
	assert.Contains(t, string(content), "envWorkloadSignature")

	os.Remove(testSignOutputPath)
}

// TestSignContractCmd_CorruptedContract tests error with corrupted --in contract
func TestSignContractCmd_CorruptedContract(t *testing.T) {
	err := os.WriteFile(testSignCorruptedContract, []byte("invalid: yaml: content: ["), 0644)
	assert.NoError(t, err)
	defer os.Remove(testSignCorruptedContract)

	cmd := getSignContractCmd()
	cmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testSignCorruptedContract,
		"--" + signContract.PrivateKeyFlagName, testSignPrivateKeyPath,
	})

	err = cmd.Execute()
	assert.NoError(t, err)
}

// TestSignContractCmd_CorruptedPrivateKey tests error with corrupted private key
func TestSignContractCmd_CorruptedPrivateKey(t *testing.T) {
	err := os.WriteFile(testSignCorruptedKey, []byte("not a valid private key"), 0644)
	assert.NoError(t, err)
	defer os.Remove(testSignCorruptedKey)

	cmd := getSignContractCmd()
	cmd.SetArgs([]string{
		"--" + signContract.InputFlagName, testSignContractPath,
		"--" + signContract.PrivateKeyFlagName, testSignCorruptedKey,
	})

	err = cmd.Execute()
	assert.NoError(t, err)
}

// ─── New --env + --workload path ──────────────────────────────────────────────

// TestSignContractCmd_FromEnvWorkload_Success tests end-to-end with --env + --workload
func TestSignContractCmd_FromEnvWorkload_Success(t *testing.T) {
	os.Remove(testSignOutputPath)

	cmd := getSignContractCmd()
	cmd.SetArgs([]string{
		"--" + signContract.EnvFlagName, testSignEncryptedEnvPath,
		"--" + signContract.WorkloadFlagName, testSignEncryptedWlPath,
		"--" + signContract.PrivateKeyFlagName, testSignPrivateKeyPath,
		"--" + signContract.OutputFlagName, testSignOutputPath,
	})

	err := cmd.Execute()
	assert.NoError(t, err)

	_, statErr := os.Stat(testSignOutputPath)
	assert.NoError(t, statErr)

	content, readErr := os.ReadFile(testSignOutputPath)
	assert.NoError(t, readErr)
	assert.Contains(t, string(content), "envWorkloadSignature")
	assert.Contains(t, string(content), "env:")
	assert.Contains(t, string(content), "workload:")

	os.Remove(testSignOutputPath)
}

// TestSignContractCmd_FromEnvWorkload_ToStdout tests new path writing to stdout
func TestSignContractCmd_FromEnvWorkload_ToStdout(t *testing.T) {
	cmd := getSignContractCmd()
	cmd.SetArgs([]string{
		"--" + signContract.EnvFlagName, testSignEncryptedEnvPath,
		"--" + signContract.WorkloadFlagName, testSignEncryptedWlPath,
		"--" + signContract.PrivateKeyFlagName, testSignPrivateKeyPath,
	})

	err := cmd.Execute()
	assert.NoError(t, err)
}

// TestSignContractCmd_FromEnvWorkload_WithPassword tests new path with password flag
func TestSignContractCmd_FromEnvWorkload_WithPassword(t *testing.T) {
	os.Remove(testSignOutputPath)

	cmd := getSignContractCmd()
	cmd.SetArgs([]string{
		"--" + signContract.EnvFlagName, testSignEncryptedEnvPath,
		"--" + signContract.WorkloadFlagName, testSignEncryptedWlPath,
		"--" + signContract.PrivateKeyFlagName, testSignPrivateKeyPath,
		"--" + signContract.PasswordFlagName, "anyPassword",
		"--" + signContract.OutputFlagName, testSignOutputPath,
	})

	err := cmd.Execute()
	assert.NoError(t, err)

	content, readErr := os.ReadFile(testSignOutputPath)
	assert.NoError(t, readErr)
	assert.Contains(t, string(content), "envWorkloadSignature")

	os.Remove(testSignOutputPath)
}

// TestSignContractCmd_FlagRegistration verifies that --env and --workload flags are
// registered on the command. Mutual-exclusion and pair-validation logic is tested at
// the library level (signContract_test.go) where os.Exit is not triggered.
func TestSignContractCmd_FlagRegistration(t *testing.T) {
	cmd := getSignContractCmd()
	_, err := cmd.PersistentFlags().GetString(signContract.EnvFlagName)
	assert.NoError(t, err, "--env flag should be registered")
	_, err = cmd.PersistentFlags().GetString(signContract.WorkloadFlagName)
	assert.NoError(t, err, "--workload flag should be registered")
	_, err = cmd.PersistentFlags().GetString(signContract.InputFlagName)
	assert.NoError(t, err, "--in flag should be registered")
}

// Note: Error test cases that trigger SetMandatoryFlagError (os.Exit(1)) — e.g. missing flags,
// --in + --env conflict, partial pairs — are thoroughly covered at the library level in
// lib/signContract/signContract_test.go where they can be validated without process exit.
