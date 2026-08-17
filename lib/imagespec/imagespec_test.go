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

package imagespec

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateInput_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, "quay.io/sclorg/postgresql-15-c9s:latest", "")
	cmd.Flags().String(OutputFlagName, "out.yaml", "")
	cmd.Flags().String(UserFlagName, "myuser", "")
	cmd.Flags().String(PassFlagName, "mypass", "")
	cmd.Flags().String(ContainerNameFlagName, "postgres", "")

	input, err := ValidateInput(cmd)

	require.NoError(t, err)
	assert.Equal(t, "quay.io/sclorg/postgresql-15-c9s:latest", input.ImageRef)
	assert.Equal(t, "out.yaml", input.OutputPath)
	assert.Equal(t, "myuser", input.Username)
	assert.Equal(t, "mypass", input.Password)
	assert.Equal(t, "postgres", input.ContainerName)
}

func TestValidateInput_OptionalFieldsOmitted(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String(InputFlagName, "quay.io/fedora/fedora:38", "")
	cmd.Flags().String(OutputFlagName, "", "")
	cmd.Flags().String(UserFlagName, "", "")
	cmd.Flags().String(PassFlagName, "", "")
	cmd.Flags().String(ContainerNameFlagName, "", "")

	input, err := ValidateInput(cmd)

	require.NoError(t, err)
	assert.Equal(t, "quay.io/fedora/fedora:38", input.ImageRef)
	assert.Equal(t, "", input.OutputPath)
	assert.Equal(t, "", input.Username)
	assert.Equal(t, "", input.Password)
	assert.Equal(t, "", input.ContainerName)
}

func TestValidateInput_FlagsNotRegistered(t *testing.T) {
	cmd := &cobra.Command{}
	_, err := ValidateInput(cmd)
	assert.Error(t, err)
}

func TestProcess_EmptyRef(t *testing.T) {
	result, err := Process("", "", "", "")
	assert.ErrorContains(t, err, "failed to generate image spec")
	assert.Empty(t, result.YAML)
	assert.Empty(t, result.ImageUser)
}

func TestProcess_InvalidRef(t *testing.T) {
	result, err := Process(":::bad:::", "", "", "")
	assert.ErrorContains(t, err, "failed to generate image spec")
	assert.Empty(t, result.YAML)
	assert.Empty(t, result.ImageUser)
}
