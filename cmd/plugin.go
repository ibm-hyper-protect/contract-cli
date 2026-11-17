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

package cmd

import (
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-cli/lib/plugin"
	"github.com/spf13/cobra"
)

// pluginCmd is the parent command that groups all plugin management sub-commands.
var pluginCmd = &cobra.Command{
	Use:   plugin.ParameterName,
	Short: plugin.ParameterShortDescription,
	Long:  plugin.ParameterLongDescription,
}

// pluginInstallCmd handles `contract-cli plugin install <ref>`.
var pluginInstallCmd = &cobra.Command{
	Use:   "install <ref>",
	Short: "Install a plugin from GitHub or a local path",
	Long: `Install a contract-cli plugin.

<ref> can be:
  github.com/org/contract-plugin-<name>           install latest release
  github.com/org/contract-plugin-<name>@v1.2.0    install specific version
  /path/to/contract-plugin-<name>                 install from local binary

GitHub releases must include a SHA-256 checksum file (<asset>.sha256, checksums.txt
or SHA256SUMS); the downloaded binary is verified against it. Use --insecure to skip
this verification.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ref, err := plugin.ValidateInstallInput(args)
		if err != nil {
			log.Fatal(err)
		}
		insecure, err := cmd.Flags().GetBool(plugin.InsecureFlag)
		if err != nil {
			log.Fatal(err)
		}
		pluginDir, err := plugin.PluginDir()
		if err != nil {
			log.Fatal(err)
		}
		if err := plugin.Install(ref, pluginDir, insecure); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Plugin installed successfully from %s\n", ref)
	},
}

// pluginListCmd handles `contract-cli plugin list`.
var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed plugins",
	Long:  "List all plugins installed in ~/.contract-cli/plugins/ along with their versions.",
	Run: func(cmd *cobra.Command, args []string) {
		pluginDir, err := plugin.PluginDir()
		if err != nil {
			log.Fatal(err)
		}
		infos, err := plugin.List(pluginDir)
		if err != nil {
			log.Fatal(err)
		}
		if len(infos) == 0 {
			fmt.Println("No plugins installed.")
			fmt.Println("Install a plugin with: contract-cli plugin install github.com/org/contract-plugin-<name>")
			return
		}
		fmt.Printf("%-25s %-15s %s\n", "NAME", "VERSION", "PATH")
		fmt.Printf("%-25s %-15s %s\n", "----", "-------", "----")
		for _, info := range infos {
			fmt.Printf("%-25s %-15s %s\n", info.Name, info.Version, info.Path)
		}
	},
}

// pluginUninstallCmd handles `contract-cli plugin uninstall <name>`.
var pluginUninstallCmd = &cobra.Command{
	Use:   "uninstall <name>",
	Short: "Uninstall a plugin",
	Long:  "Remove an installed plugin by its short name (e.g. 'ccrv', not 'contract-plugin-ccrv').",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, err := plugin.ValidateInstallInput(args)
		if err != nil {
			log.Fatal(err)
		}
		pluginDir, err := plugin.PluginDir()
		if err != nil {
			log.Fatal(err)
		}
		if err := plugin.Uninstall(name, pluginDir); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Plugin %q uninstalled successfully.\n", name)
	},
}

// pluginUpdateCmd handles `contract-cli plugin update [<name[@version]>]`.
var pluginUpdateCmd = &cobra.Command{
	Use:   "update [<name[@version]>]",
	Short: "Update one or all plugins",
	Long: `Update an installed plugin to its latest release, or to a specific version.

  contract-cli plugin update                  update all GitHub-sourced plugins
  contract-cli plugin update ccrv             update ccrv to latest
  contract-cli plugin update ccrv@v2.0.0      update ccrv to v2.0.0

Plugins installed from a local path cannot be auto-updated and will be skipped.

Downloads are verified against the release's SHA-256 checksum file. Use --insecure
to skip this verification.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		insecure, err := cmd.Flags().GetBool(plugin.InsecureFlag)
		if err != nil {
			log.Fatal(err)
		}
		pluginDir, err := plugin.PluginDir()
		if err != nil {
			log.Fatal(err)
		}
		ref := ""
		if len(args) == 1 {
			ref = args[0]
		}
		if err := plugin.Update(ref, pluginDir, insecure); err != nil {
			log.Fatal(err)
		}
		if ref == "" {
			fmt.Println("All plugins updated successfully.")
		} else {
			fmt.Printf("Plugin %q updated successfully.\n", ref)
		}
	},
}

func init() {
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.AddCommand(pluginInstallCmd)
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginUninstallCmd)
	pluginCmd.AddCommand(pluginUpdateCmd)

	pluginInstallCmd.Flags().Bool(plugin.InsecureFlag, false, plugin.InsecureFlagDescription)
	pluginUpdateCmd.Flags().Bool(plugin.InsecureFlag, false, plugin.InsecureFlagDescription)

	// pluginCmd is a sub-command group with no flags of its own — do NOT apply
	// SetCustomHelpTemplate here, as that would replace Cobra's default renderer
	// and suppress the sub-command listing entirely.
	common.SetCustomErrorTemplate(pluginCmd)

	// Leaf sub-commands use positional args, not named flags.
	// Pass an empty requiredFlags map — the template shows no Mandatory Flags section.
	common.SetCustomHelpTemplate(pluginInstallCmd, map[string]bool{})
	common.SetCustomErrorTemplate(pluginInstallCmd)
	common.SetCustomHelpTemplate(pluginListCmd, map[string]bool{})
	common.SetCustomErrorTemplate(pluginListCmd)
	common.SetCustomHelpTemplate(pluginUninstallCmd, map[string]bool{})
	common.SetCustomErrorTemplate(pluginUninstallCmd)
	common.SetCustomHelpTemplate(pluginUpdateCmd, map[string]bool{})
	common.SetCustomErrorTemplate(pluginUpdateCmd)
}
