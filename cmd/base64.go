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

		formatType, err := cmd.Flags().GetString(common.DataFormatFlagName)
		if err != nil {
			log.Fatal(err)
		}

		outputPath, err := cmd.Flags().GetString(common.FileOutFlagName)
		if err != nil {
			log.Fatal(err)
		}

		if inputData == "" {
			log.Fatal("Input data is missing")
		}

		var base64String string

		if formatType == common.DataFormatText {
			base64String, _, _, err = contract.HpcrText(inputData)
			if err != nil {
				log.Fatal(err)
			}
		} else if formatType == common.DataFormatJson {
			base64String, _, _, err = contract.HpcrJson(inputData)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Invalid input format")
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

	base64Cmd.PersistentFlags().String(common.FileInFlagName, "", common.Base64InputFlagDescription)
	base64Cmd.PersistentFlags().String(common.DataFormatFlagName, common.DataFormatText, common.Base64InputFormatFlagDescription)
	base64Cmd.PersistentFlags().String(common.FileOutFlagName, "", common.Base64OutputPathFlagDescription)
}
