package cmd

import (
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/contract"
	"github.com/spf13/cobra"
)

// base64TgzCmd represents the base64Tgz command
var base64TgzCmd = &cobra.Command{
	Use:   common.Base64TgzParamName,
	Short: common.Base64TgzParamShortDescription,
	Long:  common.Base64TgzParamLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		inputData, outputFormat, hyperProtectVersion, encCert, outputPath, err := validateInputBase64Tgz(cmd)
		if err != nil {
			log.Fatal(err)
		}

		base64TgzData, err := processBase64Tgz(inputData, outputFormat, hyperProtectVersion, encCert)
		if err != nil {
			log.Fatal(err)
		}

		err = printBase64Tgz(base64TgzData, outputPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(base64TgzCmd)

	base64TgzCmd.PersistentFlags().String(common.FileInFlagName, "", common.Base64TgzInputFlagDescription)
	base64TgzCmd.PersistentFlags().String(common.Base64TgzOutputFormatFlagName, common.Base64TgzOutputFormatDefault, common.Base64TgzOutputFormatFlagDescription)
	base64TgzCmd.PersistentFlags().String(common.OsVersionFlagName, "", common.OsVersionFlagDescription)
	base64TgzCmd.PersistentFlags().String(common.CertFlagName, "", common.CertFlagDescription)
	base64TgzCmd.PersistentFlags().String(common.FileOutFlagName, "", common.Base64TgzOutputPathDescription)
}

func validateInputBase64Tgz(cmd *cobra.Command) (string, string, string, string, string, error) {
	inputData, err := cmd.Flags().GetString(common.FileInFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	outputFormat, err := cmd.Flags().GetString(common.Base64TgzOutputFormatFlagName)
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

	return inputData, outputFormat, hyperProtectVersion, encCertPath, outputPath, nil
}

func processBase64Tgz(inputData, outputFormat, hyperProtectVersion, encCertPath string) (string, error) {
	if outputFormat == common.Base64TgzOutputFormatUnencrypted {
		if !common.CheckFileFolderExists(inputData) {
			return "", fmt.Errorf("the path to docker-compose.yaml or pods.yaml is not accessible")
		}

		base64Data, _, _, err := contract.HpcrTgz(inputData)
		if err != nil {
			return "", fmt.Errorf("failed to generate base64 tar - %v", err)
		}

		return base64Data, nil
	} else if outputFormat == common.Base64TgzOutputFormatencrypted {
		encCert, err := common.GetEncryptionCertificate(encCertPath)
		if err != nil {
			return "", err
		}

		encryptedBase64Data, _, _, err := contract.HpcrTgzEncrypted(inputData, hyperProtectVersion, encCert)
		if err != nil {
			return "", fmt.Errorf("failed to generate encrypted base64 tar - %v", err)
		}

		return encryptedBase64Data, nil
	} else {
		return "", fmt.Errorf("invalid output format (supported: plain / encrypt)")
	}
}

func printBase64Tgz(tarTgzData, outputPath string) error {
	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, tarTgzData)
		if err != nil {
			return err
		}
		fmt.Println("Successfully stored tar tgz data")
	} else {
		fmt.Println(tarTgzData)
	}

	return nil
}
