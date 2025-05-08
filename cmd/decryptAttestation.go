package cmd

import (
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-go/attestation"
	"github.com/spf13/cobra"

	"github.com/ibm-hyper-protect/contract-cli/common"
)

// decryptAttestationCmd represents the decryptAttestation command
var decryptAttestationCmd = &cobra.Command{
	Use:   common.DecryptAttestParamName,
	Short: common.DecryptAttestParamShortDescription,
	Long:  common.DecryptAttestParamLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		encAttestPath, privateKeyPath, decryptedAttestPath, err := ValidateInputDecryptedAttestation(cmd)
		if err != nil {
			log.Fatal(err)
		}

		err = DecryptAttestationRecords(encAttestPath, privateKeyPath, decryptedAttestPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(decryptAttestationCmd)

	decryptAttestationCmd.PersistentFlags().String(common.FileInFlagName, common.DecryptAttestFileInDefaultPath, common.DecryptAttestFileInDescription)
	decryptAttestationCmd.PersistentFlags().String(common.PrivateKeyFlagName, "", common.PrivateKeyFlagDescription)
	decryptAttestationCmd.PersistentFlags().String(common.FileOutFlagName, common.DecryptAttestFileOutDefaultPath, common.DecryptAttestFlagDescription)
}

func ValidateInputDecryptedAttestation(cmd *cobra.Command) (string, string, string, error) {
	encAttestPath, err := cmd.Flags().GetString(common.FileInFlagName)
	if err != nil {
		return "", "", "", err
	}
	privateKeyPath, err := cmd.Flags().GetString(common.PrivateKeyFlagName)
	if err != nil {
		return "", "", "", err
	}
	decryptedAttestPath, err := cmd.Flags().GetString(common.FileOutFlagName)
	if err != nil {
		return "", "", "", err
	}

	return encAttestPath, privateKeyPath, decryptedAttestPath, nil
}

func DecryptAttestationRecords(encryptedAttestationRecordsPath, privateKeyPath, decryptedAttestationPath string) error {
	if !common.CheckFileFolderExists(encryptedAttestationRecordsPath) || !common.CheckFileFolderExists(privateKeyPath) {
		log.Fatal("The path to encrypted attestation records file or private key doesn't exists")
	}

	encryptedChecksum, err := common.ReadDataFromFile(encryptedAttestationRecordsPath)
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := common.ReadDataFromFile(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	decryptedAttestationRecords, err := attestation.HpcrGetAttestationRecords(encryptedChecksum, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = common.WriteDataToFile(decryptedAttestationPath, decryptedAttestationRecords)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully decrypted attestation records")
	return nil
}
