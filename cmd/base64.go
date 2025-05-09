package cmd

import (
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/contract"
	"github.com/spf13/cobra"
)

// base64Cmd represents the base64 command
var base64Cmd = &cobra.Command{
	Use:   common.Base64ParamName,
	Short: common.Base64ParamShortDescription,
	Long:  common.Base64ParamLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		inputData, err := cmd.Flags().GetString(common.FileInFlagName)
		if err != nil {
			log.Fatal(err)
		}

		if inputData == "" {
			log.Fatal("Input data is missing")
		}

		outputPath, err := cmd.Flags().GetString(common.FileOutFlagName)
		if err != nil {
			log.Fatal(err)
		}

		base64String, _, _, err := contract.HpcrText(inputData)
		if err != nil {
			log.Fatal(err)
		}

		if outputPath != "" {
			err := common.WriteDataToFile(outputPath, base64String)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Successfully generated Base64")
		} else {
			fmt.Println(base64String)
		}
	},
}

func init() {
	rootCmd.AddCommand(base64Cmd)

	base64Cmd.PersistentFlags().String(common.FileInFlagName, "", common.Base64InputDescription)
	base64Cmd.PersistentFlags().String(common.FileOutFlagName, "", common.Base64OutputPathDescription)
}
