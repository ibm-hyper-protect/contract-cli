package cmd

import (
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/contract"
	"github.com/spf13/cobra"
)

const (
	inputMissingMessageBase64 = "input data is missing"
	invalidInputMessageBase64 = "invalid input format"
	successMessageBase64      = "successfully generated Base64"
)

// base64Cmd represents the base64 command
var base64Cmd = &cobra.Command{
	Use:   common.Base64ParamName,
	Short: common.Base64ParamShortDescription,
	Long:  common.Base64ParamLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		inputData, formatType, outputPath, err := validateInputBase64(cmd)
		if err != nil {
			log.Fatal(err)
		}

		base64String, err := processBase64(inputData, formatType)
		if err != nil {
			log.Fatal(err)
		}

		err = printBase64(base64String, outputPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(base64Cmd)

	base64Cmd.PersistentFlags().String(common.FileInFlagName, "", common.Base64InputFlagDescription)
	base64Cmd.PersistentFlags().String(common.DataFormatFlagName, common.DataFormatText, common.Base64InputFormatFlagDescription)
	base64Cmd.PersistentFlags().String(common.FileOutFlagName, "", common.Base64OutputPathFlagDescription)
}

func validateInputBase64(cmd *cobra.Command) (string, string, string, error) {
	inputData, err := cmd.Flags().GetString(common.FileInFlagName)
	if err != nil {
		return "", "", "", err
	}

	formatType, err := cmd.Flags().GetString(common.DataFormatFlagName)
	if err != nil {
		return "", "", "", err
	}

	outputPath, err := cmd.Flags().GetString(common.FileOutFlagName)
	if err != nil {
		return "", "", "", err
	}

	return inputData, formatType, outputPath, nil
}

func processBase64(inputData, formatType string) (string, error) {
	var base64String string
	var err error

	if inputData == "" {
		return "", fmt.Errorf(inputMissingMessageBase64)
	}

	if formatType == common.DataFormatText {
		base64String, _, _, err = contract.HpcrText(inputData)
		if err != nil {
			return "", err
		}
	} else if formatType == common.DataFormatJson {
		base64String, _, _, err = contract.HpcrJson(inputData)
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf(invalidInputMessageBase64)
	}

	return base64String, nil
}

func printBase64(base64String, outputPath string) error {
	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, base64String)
		if err != nil {
			return err
		}
		fmt.Println(successMessageBase64)
	} else {
		fmt.Println(base64String)
	}

	return nil
}
