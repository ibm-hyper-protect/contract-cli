package cmd

import (
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/contract"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   common.EncryptParamName,
	Short: common.EncryptParamShortDescription,
	Long:  common.EncryptParamLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		inputDataPath, osVersion, certPath, privateKeyPath, outputPath, err := validateInputEncrypt(cmd)
		if err != nil {
			log.Fatal(err)
		}

		signedEncryptContract, err := generateSignedEncryptContract(inputDataPath, osVersion, certPath, privateKeyPath)
		if err != nil {
			log.Fatal(err)
		}

		err = printSignedEncryptContract(signedEncryptContract, outputPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.PersistentFlags().String(common.FileInFlagName, "", common.EncryptInputFlagDescription)
	encryptCmd.PersistentFlags().String(common.OsVersionFlagName, "", common.OsVersionFlagDescription)
	encryptCmd.PersistentFlags().String(common.CertFlagName, "", common.CertFlagDescription)
	encryptCmd.PersistentFlags().String(common.PrivateKeyFlagName, "", common.PrivateKeyFlagDescription)
	encryptCmd.PersistentFlags().String(common.FileOutFlagName, "", common.EncryptOutputFlagDescription)
}

func validateInputEncrypt(cmd *cobra.Command) (string, string, string, string, string, error) {
	inputData, err := cmd.Flags().GetString(common.FileInFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	osVersion, err := cmd.Flags().GetString(common.OsVersionFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	certPath, err := cmd.Flags().GetString(common.CertFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	privateKeyPath, err := cmd.Flags().GetString(common.PrivateKeyFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	outputPath, err := cmd.Flags().GetString(common.FileOutFlagName)
	if err != nil {
		return "", "", "", "", "", err
	}

	return inputData, osVersion, certPath, privateKeyPath, outputPath, nil
}

func generateSignedEncryptContract(inputDataPath, osVersion, certPath, privateKeyPath string) (string, error) {
	if !common.CheckFileFolderExists(inputDataPath) {
		log.Fatal("The contract path doesn't exist")
	}

	inputData, err := common.ReadDataFromFile(inputDataPath)
	if err != nil {
		return "", err
	}

	var privateKey string
	if privateKeyPath == "" {
		privateKey, err = generatePrivateKey()
		if err != nil {
			return "", err
		}
	} else {
		if common.CheckFileFolderExists(privateKeyPath) {
			privateKey, err = common.ReadDataFromFile(privateKeyPath)
			if err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("private key path doesn't exist")
		}
	}

	signedEncryptContract, _, _, err := contract.HpcrContractSignedEncrypted(inputData, osVersion, certPath, privateKey)
	if err != nil {
		return "", err
	}

	return signedEncryptContract, nil
}

func generatePrivateKey() (string, error) {
	err := common.OpensslCheck()
	if err != nil {
		return "", fmt.Errorf("openssl not found - %v", err)
	}

	privateKey, err := common.ExecCommand("openssl", "", "genrsa", "4096")
	if err != nil {
		return "", fmt.Errorf("failed to generate private key - %v", err)
	}

	return privateKey, nil
}

func printSignedEncryptContract(signedEncryptContract, outputPath string) error {
	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, signedEncryptContract)
		if err != nil {
			return err
		}
		fmt.Println("Successfully generated signed and encrypted contract")
	} else {
		fmt.Println(signedEncryptContract)
	}

	return nil
}
