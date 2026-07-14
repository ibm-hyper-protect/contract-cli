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
	for _, osVal := range []string{OsHpvs, OsCcrt, OsCcrv, OsCccoPeerpod, OsCccoBmtl} {
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

// TestValidateInput_LegacyCccoRejected verifies that the old "ccco" value is no longer accepted.
func TestValidateInput_LegacyCccoRejected(t *testing.T) {
	_, _, _, err := ValidateInput(newCmd(TypeContract, "ccco", ""))
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
// returned for each OS.
func TestGenerateContractTemplate_WorkloadPerOs(t *testing.T) {
	cases := []struct {
		os                  string
		hasCompose          bool
		hasConfidentialCont bool
	}{
		{"", true, false},
		{OsHpvs, true, false},
		{OsCcrt, true, false},
		{OsCcrv, false, false},
		{OsCccoPeerpod, false, true},
		{OsCccoBmtl, false, true},
	}
	for _, tc := range cases {
		result, err := GenerateContractTemplate(TypeWorkload, tc.os)
		assert.NoError(t, err, "os=%s", tc.os)
		assert.Contains(t, result, "type: workload", "os=%s", tc.os)
		if tc.hasCompose {
			assert.Contains(t, result, "compose:", "os=%s should have compose", tc.os)
		} else {
			assert.NotContains(t, result, "compose:", "os=%s should not have compose", tc.os)
		}
		if tc.hasConfidentialCont {
			assert.Contains(t, result, "confidential-containers:", "os=%s should have confidential-containers", tc.os)
		} else {
			assert.NotContains(t, result, "confidential-containers:", "os=%s should not have confidential-containers", tc.os)
		}
	}
}

// TestGenerateContractTemplate_EnvPerOs verifies correct env template is returned per OS.
// hpvs/ccrt/ccrv/"" use the standard env; ccco-peerpod and ccco-bmtl have distinct templates.
func TestGenerateContractTemplate_EnvPerOs(t *testing.T) {
	// Standard env (no host-attestation, no volumes)
	for _, osVal := range []string{"", OsHpvs, OsCcrt, OsCcrv} {
		result, err := GenerateContractTemplate(TypeEnv, osVal)
		assert.NoError(t, err, "os=%s", osVal)
		assert.Contains(t, result, "type: env", "os=%s", osVal)
	}

	// ccco-peerpod: has logRouter, no volumes, no host-attestation
	peerpodEnv, err := GenerateContractTemplate(TypeEnv, OsCccoPeerpod)
	assert.NoError(t, err)
	assert.Contains(t, peerpodEnv, "type: env")
	assert.Contains(t, peerpodEnv, "logRouter:")
	assert.NotContains(t, peerpodEnv, "host-attestation:")
	assert.NotContains(t, peerpodEnv, "volumes:")

	// ccco-bmtl: has logRouter, volumes, and host-attestation
	bmtlEnv, err := GenerateContractTemplate(TypeEnv, OsCccoBmtl)
	assert.NoError(t, err)
	assert.Contains(t, bmtlEnv, "type: env")
	assert.Contains(t, bmtlEnv, "logRouter:")
	assert.Contains(t, bmtlEnv, "volumes:")
	assert.Contains(t, bmtlEnv, "host-attestation:")

	// ccco-peerpod and ccco-bmtl env templates must differ from each other
	assert.NotEqual(t, peerpodEnv, bmtlEnv)
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
	assert.NotContains(t, ccrvCombined, "compose:")

	// ccco-peerpod combined has confidential-containers in workload section
	peerpodCombined, err := GenerateContractTemplate(TypeContract, OsCccoPeerpod)
	assert.NoError(t, err)
	assert.Contains(t, peerpodCombined, "workload:")
	assert.Contains(t, peerpodCombined, "env:")
	assert.Contains(t, peerpodCombined, "confidential-containers:")
	assert.NotContains(t, peerpodCombined, "host-attestation:")

	// ccco-bmtl combined has confidential-containers and host-attestation
	bmtlCombined, err := GenerateContractTemplate(TypeContract, OsCccoBmtl)
	assert.NoError(t, err)
	assert.Contains(t, bmtlCombined, "workload:")
	assert.Contains(t, bmtlCombined, "env:")
	assert.Contains(t, bmtlCombined, "confidential-containers:")
	assert.Contains(t, bmtlCombined, "host-attestation:")
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
