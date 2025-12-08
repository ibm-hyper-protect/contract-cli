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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "contract-cli",
		Short: "A CLI tool for managing IBM Hyper Protect contracts and configurations",
		Long: `Contract CLI automates contract generation and management for IBM Hyper Protect services.
Supports:
  - Hyper Protect Virtual Servers (HPVS) for VPC
  - Hyper Protect Container Runtime (HPCR) for RHVS
  - Hyper Protect Confidential Container (HPCC) Peer Pods

Documentation: https://github.com/ibm-hyper-protect/contract-cli/blob/main/docs/README.md`,
		Version: cliVersion,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s version %s %s %s\nRelease: %s\n", "contract-cli", cliVersion, cliOsName, cliOsArch, cliBuildDate)
		},
	}

	cliVersion   = "dev"
	cliOsName    = "unknown"
	cliOsArch    = "unknown"
	cliBuildDate = "unknown"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Disable completion command
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// SetVersionInfo - function to set CLI details from build parameters
func SetVersionInfo(version, osName, osArch, buildDate string) {
	cliVersion = version
	cliOsName = osName
	cliOsArch = osArch
	cliBuildDate = buildDate

	rootCmd.Version = fmt.Sprintf("%s %s %s build %s", version, osName, osArch, buildDate)
}
