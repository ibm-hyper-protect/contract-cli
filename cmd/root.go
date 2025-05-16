package cmd

import (
	"os"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   common.ContractCliName,
	Short: common.ContractCliShortDescription,
	Long:  common.ContractCliLongDescription,
}

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
