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
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testOutputPath  = "../../build/test_contract_template_output.yaml"
	testInvalidPath = "../../build/nonexistent/dir/output.yaml"
)

// newCmd creates a cobra.Command with all flags registered for use in tests.
func newCmd(typeFlagVal, osFlagVal, outFlagVal string) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().String(TypeFlagName, typeFlagVal, "")
	cmd.Flags().String(OsVersionFlagName, osFlagVal, "")
	cmd.Flags().String(OutputFlagName, outFlagVal, "")
	return cmd
}

// ---------------------------------------------------------------------------
// ValidateInput
// ---------------------------------------------------------------------------

// TestValidateInput_ValidTypes verifies all accepted --type values are parsed correctly.
func TestValidateInput_ValidTypes(t *testing.T) {
	cases := []struct{ input, want string }{
		{TypeEnv, TypeEnv},
		{TypeWorkload, TypeWorkload},
		{TypeContract, TypeContract},
		{"", ""},
	}
	for _, tc := range cases {
		cmd := newCmd(tc.input, "", "")
		got, _, _, err := ValidateInput(cmd)
		assert.NoError(t, err, "type=%q", tc.input)
		assert.Equal(t, tc.want, got, "type=%q", tc.input)
	}
}

// TestValidateInput_InvalidType verifies unsupported --type returns an error.
func TestValidateInput_InvalidType(t *testing.T) {
	_, _, _, err := ValidateInput(newCmd("invalid-type", "", ""))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid value for --type")
}

// TestValidateInput_ValidOs verifies all accepted --os values are parsed correctly.
func TestValidateInput_ValidOs(t *testing.T) {
	for _, osVal := range []string{OsHpvs, OsCcrt, OsCcrv, OsCcco} {
		_, got, _, err := ValidateInput(newCmd(TypeContract, osVal, ""))
		assert.NoError(t, err, "os=%s", osVal)
		assert.Equal(t, osVal, got, "os=%s", osVal)
	}
}

// TestValidateInput_InvalidOs verifies unsupported --os returns an error.
func TestValidateInput_InvalidOs(t *testing.T) {
	_, _, _, err := ValidateInput(newCmd(TypeContract, "bad-os", ""))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid value for --os")
}

// TestValidateInput_OptionalFlags verifies empty --os and --out are accepted and forwarded.
func TestValidateInput_OptionalFlags(t *testing.T) {
	templateType, osVersion, outputPath, err := ValidateInput(newCmd(TypeContract, "", testOutputPath))
	assert.NoError(t, err)
	assert.Equal(t, TypeContract, templateType)
	assert.Equal(t, "", osVersion)
	assert.Equal(t, testOutputPath, outputPath)
}

// TestValidateInput_FlagNotRegistered verifies missing flag registration returns an error.
func TestValidateInput_FlagNotRegistered(t *testing.T) {
	_, _, _, err := ValidateInput(&cobra.Command{})
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// GenerateContractTemplate
// ---------------------------------------------------------------------------

// TestGenerateContractTemplate_WorkloadPerOs verifies the correct workload template is
// returned for each OS: ccrt/hpvs/ccco/"" include compose; ccrv does not.
func TestGenerateContractTemplate_WorkloadPerOs(t *testing.T) {
	cases := []struct {
		os         string
		hasCompose bool
	}{
		{"", true},
		{OsHpvs, true},
		{OsCcrt, true},
		{OsCcco, true},
		{OsCcrv, false},
	}
	for _, tc := range cases {
		result, err := GenerateContractTemplate(TypeWorkload, tc.os)
		assert.NoError(t, err, "os=%s", tc.os)
		assert.Contains(t, result, "type: workload", "os=%s", tc.os)
		assert.Contains(t, result, "play:", "os=%s", tc.os)
		if tc.hasCompose {
			assert.Contains(t, result, "compose:", "os=%s should have compose", tc.os)
		} else {
			assert.NotContains(t, result, "compose:", "os=%s should not have compose", tc.os)
		}
	}
}

// TestGenerateContractTemplate_EnvSameForAllOs verifies env template is identical for all OS values.
func TestGenerateContractTemplate_EnvSameForAllOs(t *testing.T) {
	base, err := GenerateContractTemplate(TypeEnv, "")
	assert.NoError(t, err)
	assert.Contains(t, base, "type: env")

	for _, osVal := range []string{OsHpvs, OsCcrt, OsCcrv, OsCcco} {
		result, err := GenerateContractTemplate(TypeEnv, osVal)
		assert.NoError(t, err, "os=%s", osVal)
		assert.Equal(t, base, result, "env template must be identical for os=%s", osVal)
	}
}

// TestGenerateContractTemplate_CombinedContainsBothSections verifies combined output has both sections.
// "contract" and "" are treated identically; ccrv combined must not contain compose.
func TestGenerateContractTemplate_CombinedContainsBothSections(t *testing.T) {
	// "contract" and "" must produce the same output
	byContract, err1 := GenerateContractTemplate(TypeContract, OsCcrt)
	byEmpty, err2 := GenerateContractTemplate("", OsCcrt)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, byContract, byEmpty)
	assert.Contains(t, byContract, "workload:")
	assert.Contains(t, byContract, "env:")

	// ccrv combined has play but not compose
	ccrvCombined, err := GenerateContractTemplate(TypeContract, OsCcrv)
	assert.NoError(t, err)
	assert.Contains(t, ccrvCombined, "workload:")
	assert.Contains(t, ccrvCombined, "env:")
	assert.Contains(t, ccrvCombined, "play:")
	assert.NotContains(t, ccrvCombined, "compose:")
}

// TestGenerateContractTemplate_InvalidType verifies an unknown type returns an error.
func TestGenerateContractTemplate_InvalidType(t *testing.T) {
	result, err := GenerateContractTemplate("bogus", "")
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "unsupported template type")
}

// ---------------------------------------------------------------------------
// Output
// ---------------------------------------------------------------------------

// TestOutput_ToFile verifies the template is written to a file.
func TestOutput_ToFile(t *testing.T) {
	data := "workload: |\n  type: workload\n"
	os.Remove(testOutputPath)

	err := Output(data, testOutputPath)
	assert.NoError(t, err)

	content, _ := os.ReadFile(testOutputPath)
	assert.Equal(t, data, string(content))
	os.Remove(testOutputPath)
}

// TestOutput_ToStdout verifies no error when printing to stdout (no file path).
func TestOutput_ToStdout(t *testing.T) {
	assert.NoError(t, Output("type: env\n", ""))
}

// TestOutput_InvalidPath verifies an error is returned for an unwriteable path.
func TestOutput_InvalidPath(t *testing.T) {
	assert.Error(t, Output("data", testInvalidPath))
}
