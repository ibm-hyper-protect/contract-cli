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
	"fmt"

	"github.com/ibm-hyper-protect/contract-cli/common"
	goImageSpec "github.com/ibm-hyper-protect/contract-go/v2/imagespec"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "image-spec"
	ParameterShortDescription = "Fetch OCI image metadata and generate a Kubernetes pod YAML template"
	ParameterLongDescription  = `Fetch OCI image metadata and generate a Kubernetes pod YAML template
with the correct env, entrypoint, user, and port overrides for use with
the registryMapping feature of confidential-containers workload contracts.`

	InputFlagName  = "in"
	OutputFlagName = "out"
	UserFlagName   = "username"
	PassFlagName   = "password"

	ImageRefDescription      = "Fully-qualified image reference to inspect (e.g. quay.io/fedora/fedora:38)"
	OutputDescription        = "Path to write the generated pod YAML template (prints to stdout when omitted)"
	UsernameDescription      = "Registry username for private image access (optional)"
	PasswordDescription      = "Registry password / API key for private image access (optional)"
	ContainerNameFlagName    = "container-name"
	ContainerNameDescription = "Container name to use in the generated pod spec (default: derived from image name)"
)

// ImageSpecInput holds the validated flag values for the image-spec command.
type ImageSpecInput struct {
	ImageRef      string
	OutputPath    string
	Username      string
	Password      string
	ContainerName string
}

// ValidateInput parses and validates all flags for the image-spec command.
func ValidateInput(cmd *cobra.Command) (ImageSpecInput, error) {
	var input ImageSpecInput
	var err error

	input.ImageRef, err = cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return ImageSpecInput{}, err
	}
	if input.ImageRef == "" {
		err = fmt.Errorf("Error: required flag '--in' is missing")
		common.SetMandatoryFlagError(cmd, err)
		return ImageSpecInput{}, err
	}

	input.OutputPath, err = cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return ImageSpecInput{}, err
	}

	input.Username, err = cmd.Flags().GetString(UserFlagName)
	if err != nil {
		return ImageSpecInput{}, err
	}

	input.Password, err = cmd.Flags().GetString(PassFlagName)
	if err != nil {
		return ImageSpecInput{}, err
	}

	input.ContainerName, err = cmd.Flags().GetString(ContainerNameFlagName)
	if err != nil {
		return ImageSpecInput{}, err
	}

	return input, nil
}

// Process generates the pod YAML template for the given image reference.
// containerName defaults to the image name when empty.
func Process(imageRef, containerName, username, password string) (string, error) {
	var auth *goImageSpec.AuthConfig
	if username != "" || password != "" {
		auth = &goImageSpec.AuthConfig{
			Username: username,
			Password: password,
		}
	}

	yaml, _, _, err := goImageSpec.GenerateImageSpec(imageRef, containerName, auth)
	if err != nil {
		return "", fmt.Errorf("failed to generate image spec for %q: %w", imageRef, err)
	}

	return yaml, nil
}
