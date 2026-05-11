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

package listCertVersions

import (
	"fmt"
	"strings"

	"github.com/ibm-hyper-protect/contract-go/v2/certificate"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "list-encryptioncert-versions"
	ParameterShortDescription = "List available encryption certificate versions"
	ParameterLongDescription  = `List all available embedded encryption certificate versions for IBM Confidential Computing platforms.

Shows certificate versions organized by platform (ccrt, ccrv, ccco, hpvs). Use this to discover
which certificate versions are available before using the --ver flag with encrypt commands.`
	OsVersionFlagName        = "os"
	OsVersionFlagDescription = "Filter by platform (ccrt, ccrv, ccco, or hpvs for legacy). Shows all platforms if not specified"
	OutputFlagName           = "out"
	OutputFlagDescription    = "Path to save output (optional, prints to stdout if not specified)"
	FormatFlagName           = "format"
	FormatFlagDescription    = "Output format: 'json', or 'yaml' (defaults to 'json' if empty)"
)

// ValidateInput validates the input flags for list-cert-versions command
func ValidateInput(cmd *cobra.Command) (string, string, string, error) {
	osVersion, err := cmd.Flags().GetString(OsVersionFlagName)
	if err != nil {
		return "", "", "", err
	}

	outputPath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", "", err
	}

	format, err := cmd.Flags().GetString(FormatFlagName)
	if err != nil {
		return "", "", "", err
	}

	// Validate format if provided
	if format != "" && format != "json" && format != "yaml" {
		return "", "", "", fmt.Errorf("invalid format '%s'. Valid formats: json, yaml", format)
	}

	return osVersion, outputPath, format, nil
}

// Process generates the certificate version listing
func Process(osVersion, format string) (string, error) {
	// Normalize OS version to lowercase if provided
	if osVersion != "" {
		osVersion = strings.ToLower(osVersion)
	}

	// Default to JSON if format is empty
	if format == "" {
		format = "json"
	}

	// For JSON and YAML formats, return the raw output from the library
	if format == "json" || format == "yaml" {
		result, err := certificate.HpcrListAvailableEncCertVersions(osVersion, format)
		if err != nil {
			if osVersion != "" {
				return "", fmt.Errorf("no certificates found for platform '%s'. Valid platforms: ccrt, ccrv, ccco, hpvs", osVersion)
			}
			return "", fmt.Errorf("failed to get available certificates - %v", err)
		}
		return result, nil
	}

	return "", fmt.Errorf("invalid format '%s'. Valid formats: json, yaml", format)
}

// Output returns the result string (no I/O, just returns the formatted string)
func Output(result string) string {
	return result
}
