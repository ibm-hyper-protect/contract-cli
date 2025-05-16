package cmd

import (
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/contract"
	"github.com/spf13/cobra"
)

const (
	successMessageEncryptString = "Successfully stored encrypted text"
)

// encryptStringCmd represents the encryptString command
var encryptStringCmd = &cobra.Command{
	Use:   common.EncryptStrParamName,
	Short: common.EncryptStrParamShortDescription,
	Long:  common.EncryptStrParamLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		inputData, inputFormat, hyperProtectVersion, encCertPath, outputPath, err := validateInputEncryptString(cmd)
		if err != nil {
			log.Fatal(err)
		}

		var encCert string
		if encCertPath != "" {
			encCert, err = common.ReadDataFromFile(encCertPath)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			encCert = ""
		}

		var encryptedString string
		if inputFormat == common.DataFormatText {
			encryptedString, _, _, err = contract.HpcrTextEncrypted(inputData, hyperProtectVersion, encCert)
			if err != nil {
				log.Fatal(err)
			}
		} else if inputFormat == common.DataFormatJson {
			encryptedString, _, _, err = contract.HpcrJsonEncrypted(inputData, hyperProtectVersion, encCert)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Invalid input format")
		}

		if outputPath != "" {
			err = common.WriteDataToFile(outputPath, encryptedString)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(successMessageEncryptString)
		} else {
			fmt.Println(encryptedString)
		}
	},
}

func init() {
	rootCmd.AddCommand(encryptStringCmd)

	encryptStringCmd.PersistentFlags().String(common.FileInFlagName, "", common.EncryptStrInputFlagDescription)
	encryptStringCmd.PersistentFlags().String(common.DataFormatFlagName, common.DataFormatText, common.EncryptStrFormatFlagDescription)
	encryptStringCmd.PersistentFlags().String(common.OsVersionFlagName, "", common.OsVersionFlagDescription)
	encryptStringCmd.PersistentFlags().String(common.CertFlagName, "", common.CertFlagDescription)
	encryptStringCmd.PersistentFlags().String(common.FileOutFlagName, "", common.EncryptStrOutputFlagDescription)
}

func validateInputEncryptString(cmd *cobra.Command) (string, string, string, string, string, error) {
	inputData, err := cmd.Flags().GetString(common.FileInFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	inputFormat, err := cmd.Flags().GetString(common.DataFormatFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	hyperProtectVersion, err := cmd.Flags().GetString(common.OsVersionFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	encCertPath, err := cmd.Flags().GetString(common.CertFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	outputPath, err := cmd.Flags().GetString(common.FileOutFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	return inputData, inputFormat, hyperProtectVersion, encCertPath, outputPath, nil
}
