package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "contract-cli",
	Short: "Contract CLI: A CLI utility tool for Hyper Protect",
	Long: `This tool helps to provision HPVS/RHVS on IBM Cloud and On Prem 

Refer github.com/ibm-hyper-protect/contract-cli for details on features`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
