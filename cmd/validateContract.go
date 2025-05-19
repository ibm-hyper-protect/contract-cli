package cmd

import (
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/contract"
	"github.com/spf13/cobra"
)

// validateContractCmd represents the validateContract command
var validateContractCmd = &cobra.Command{
	Use:   common.ValidateContractParamName,
	Short: common.ValidateContractParamShortDescription,
	Long:  common.ValidateContractParamLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		contractPath, version, err := validateInputContract(cmd)
		if err != nil {
			log.Fatal(err)
		}

		contractData, err := common.ReadDataFromFile(contractPath)
		if err != nil {
			log.Fatal(err)
		}

		if !common.CheckFileFolderExists(contractPath) {
			log.Fatal("The path to contract doesn't exist")
		}

		err = contract.HpcrVerifyContract(contractData, version)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("True")
	},
}

func init() {
	rootCmd.AddCommand(validateContractCmd)

	validateContractCmd.PersistentFlags().String(common.FileInFlagName, "", common.ValidateContractInputFlagDescription)
	validateContractCmd.PersistentFlags().String(common.OsVersionFlagName, "", common.OsVersionFlagDescription)
}

func validateInputContract(cmd *cobra.Command) (string, string, error) {
	contract, err := cmd.Flags().GetString(common.FileInFlagName)
	if err != nil {
		return "", "", err
	}

	version, err := cmd.Flags().GetString(common.OsVersionFlagName)
	if err != nil {
		return "", "", err
	}

	return contract, version, nil
}
