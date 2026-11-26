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
	"strings"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-cli/lib/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// setPluginHelp installs the help renderer used by the plugin leaf commands.
// common.SetCustomHelpTemplate only prints "Usage: <path> [flags]" and the flags, which
// hides the positional arguments, the long description and the examples that these
// commands rely on. This renderer keeps the same flag layout and adds those sections.
func setPluginHelp(cmd *cobra.Command) {
	cmd.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
		out := cmd.OutOrStdout()

		if desc := strings.TrimSpace(cmd.Long); desc != "" {
			fmt.Fprintf(out, "%s\n\n", desc)
		}

		fmt.Fprintf(out, "Usage:\n  %s\n\n", cmd.UseLine())

		if examples := strings.TrimRight(cmd.Example, "\n"); strings.TrimSpace(examples) != "" {
			fmt.Fprintf(out, "Examples:\n%s\n\n", examples)
		}

		maxLen := 0
		cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
			if len(f.Name) > maxLen {
				maxLen = len(f.Name)
			}
		})
		fmt.Fprintln(out, "Optional Flags:")
		cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
			fmt.Fprintf(out, "  --%s%s%s\n", f.Name, strings.Repeat(" ", maxLen-len(f.Name)+2), f.Usage)
		})
	})
}

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
	Long: `Install a contract-cli plugin into ~/.contract-cli/plugins/.

<ref> can be:
  github.com/org/contract-plugin-<name>           install the latest GitHub release
  github.com/org/contract-plugin-<name>@v1.2.0    install a specific GitHub release
  /path/to/contract-plugin-<name>                 install from a local binary

GitHub installs download the release asset that matches this machine
(contract-plugin-<name>-<os>-<arch>) and verify it against the SHA-256 checksum file
published in the same release (<asset>.sha256, checksums.txt or SHA256SUMS). The
install fails if the checksum file is missing or does not match. Use --insecure to
skip this verification.

Local binaries must be named contract-plugin-<name> and are copied as-is; they are
never checksum-verified.

Once installed, run the plugin with: contract-cli <name> [args]`,
	Example: `  contract-cli plugin install github.com/ibm-hyper-protect/contract-plugin-ccrv
  contract-cli plugin install github.com/ibm-hyper-protect/contract-plugin-ccrv@v1.2.0
  contract-cli plugin install ./contract-plugin-ccrv
  contract-cli plugin install github.com/ibm-hyper-protect/contract-plugin-ccrv --insecure`,
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
	Long: `List all plugins installed in ~/.contract-cli/plugins/ with their version, install method and installation path.

The version is taken from the plugin's --version output. The install method is 'github'
for plugins installed from a GitHub release (these can be auto-updated) and 'local' for
plugins installed from a local binary.`,
	Example: `  contract-cli plugin list`,
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
		fmt.Print(plugin.FormatList(infos))
	},
}

// pluginUninstallCmd handles `contract-cli plugin uninstall <name>`.
var pluginUninstallCmd = &cobra.Command{
	Use:     "uninstall <name>",
	Short:   "Uninstall a plugin",
	Long:    "Remove an installed plugin by its short name (e.g. 'ccrv', not 'contract-plugin-ccrv').",
	Example: `  contract-cli plugin uninstall ccrv`,
	Args:    cobra.ExactArgs(1),
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

	// Leaf sub-commands take positional args rather than mandatory flags, so they use
	// setPluginHelp (which shows the args, long description and examples) instead of
	// common.SetCustomHelpTemplate.
	setPluginHelp(pluginInstallCmd)
	common.SetCustomErrorTemplate(pluginInstallCmd)
	setPluginHelp(pluginListCmd)
	common.SetCustomErrorTemplate(pluginListCmd)
	setPluginHelp(pluginUninstallCmd)
	common.SetCustomErrorTemplate(pluginUninstallCmd)
	setPluginHelp(pluginUpdateCmd)
	common.SetCustomErrorTemplate(pluginUpdateCmd)
}
