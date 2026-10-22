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
	"strings"

	"github.com/ibm-hyper-protect/contract-cli/lib/plugin"
	"github.com/spf13/cobra"
)

const (
	builtinGroupID       = "builtin"
	pluginGroupID        = "plugins"
	toolName             = "contract-cli"
	toolShortDescription = "A CLI tool for managing IBM Confidential Computing contracts and configurations"
	toolLongDescription  = `Contract CLI automates contract generation and management for IBM Confidential Computing services.

Supports:
  - IBM Confidential Computing Container Runtime
  - IBM Confidential Computing Container Runtime for Red Hat Virtualization Solutions
  - IBM Confidential Computing Containers for Red Hat OpenShift Container Platform

Documentation: https://ibm-hyper-protect.github.io/contract-cli/`
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:     toolName,
		Short:   toolShortDescription,
		Long:    toolLongDescription,
		Version: cliVersion,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s version %s %s %s\nRelease: %s\n", toolName, cliVersion, cliOsName, cliOsArch, cliBuildDate)
		},
	}

	cliVersion   = "dev"
	cliOsName    = "unknown"
	cliOsArch    = "unknown"
	cliBuildDate = "unknown"
)

// registerPluginCommands adds a stub command for every installed plugin so that
// they are listed under "Plugin Commands" in `contract-cli --help`.
// Dispatch itself is still performed by plugin.Exec; the stub passes all arguments
// (including --help) straight through to the plugin binary.
// Plugins whose name collides with a built-in command are skipped (built-ins win).
func registerPluginCommands(pluginDir string) {
	names, err := plugin.ListPlugins(pluginDir)
	if err != nil {
		return
	}

	builtins := map[string]bool{}
	for _, c := range rootCmd.Commands() {
		builtins[c.Name()] = true
	}

	groupAdded := false
	for _, bin := range names {
		name := strings.TrimSuffix(strings.TrimPrefix(bin, plugin.PluginPrefix), ".exe")
		if name == "" || builtins[name] {
			continue
		}
		path, found := plugin.FindPlugin(name, pluginDir)
		if !found {
			continue
		}
		if !groupAdded {
			// Cobra relabels ungrouped commands "Additional Commands" once any group
			// exists, so put the built-ins in an explicit group listed first.
			rootCmd.AddGroup(
				&cobra.Group{ID: builtinGroupID, Title: "Available Commands:"},
				&cobra.Group{ID: pluginGroupID, Title: "Plugin Commands:"},
			)
			for _, c := range rootCmd.Commands() {
				c.GroupID = builtinGroupID
			}
			rootCmd.SetHelpCommandGroupID(builtinGroupID)
			groupAdded = true
		}
		rootCmd.AddCommand(&cobra.Command{
			Use:                name,
			Short:              fmt.Sprintf("Plugin command (run '%s %s --help')", toolName, name),
			GroupID:            pluginGroupID,
			DisableFlagParsing: true,
			RunE: func(cmd *cobra.Command, args []string) error {
				return plugin.Exec(path, args)
			},
		})
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Disable completion command
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true

	// List installed plugins in the help output.
	if pluginDir, err := plugin.PluginDir(); err == nil {
		registerPluginCommands(pluginDir)
	}

	// Install the plugin dispatch hook before Cobra runs.
	// If the first argument is not a known built-in command, delegate to a plugin binary.
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return nil // built-in commands execute normally
	}

	// Override the Args validator on rootCmd so unknown sub-commands are intercepted
	// before Cobra emits an "unknown command" error.
	rootCmd.Args = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return nil
		}
		// Check whether a plugin binary exists for the first argument.
		pluginDir, err := plugin.PluginDir()
		if err != nil {
			return err
		}
		if path, found := plugin.FindPlugin(args[0], pluginDir); found {
			// Delegate to the plugin — this replaces (Unix) or exits (Windows) the process.
			if err := plugin.Exec(path, args[1:]); err != nil {
				return err
			}
		}
		// No plugin found — fall through to Cobra's default unknown-command error.
		return fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
	}

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
