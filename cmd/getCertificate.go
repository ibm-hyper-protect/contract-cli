package cmd

import (
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/certificate"
	"github.com/spf13/cobra"
)

// getCertificateCmd represents the getCertificate command
var getCertificateCmd = &cobra.Command{
	Use:   common.GetCertParamName,
	Short: common.GetCertParamShortDescription,
	Long:  common.GetCertParamLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		encryptionCertsPath, version, encryptionCertOutputPath, err := validateInputGetCertificate(cmd)
		if err != nil {
			log.Fatal(err)
		}

		encryptionCertificate, err := getEncryptionCertificate(encryptionCertsPath, version)
		if err != nil {
			log.Fatal(err)
		}

		err = printCertificate(encryptionCertificate, encryptionCertOutputPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCertificateCmd)

	getCertificateCmd.PersistentFlags().String(common.FileInFlagName, "", common.GetCertFileInFlagDescription)
	getCertificateCmd.PersistentFlags().String(common.VersionFlagName, "", common.GetCertVersionFlagDescription)
	getCertificateCmd.PersistentFlags().String(common.FileOutFlagName, "", common.GetCertFileOutFlagDescription)
}

func validateInputGetCertificate(cmd *cobra.Command) (string, string, string, error) {
	encryptionCertsPath, err := cmd.Flags().GetString(common.FileInFlagName)
	if err != nil {
		return "", "", "", err
	}

	version, err := cmd.Flags().GetString(common.VersionFlagName)
	if err != nil {
		return "", "", "", err
	}

	encryptionCertificatePath, err := cmd.Flags().GetString(common.FileOutFlagName)
	if err != nil {
		return "", "", "", err
	}

	return encryptionCertsPath, version, encryptionCertificatePath, nil
}

func getEncryptionCertificate(encryptionCertsPath, version string) (string, error) {
	if !common.CheckFileFolderExists(encryptionCertsPath) {
		return "", fmt.Errorf("the path to encryption certificates doesn't exist")
	}

	encryptionCertsJson, err := common.ReadDataFromFile(encryptionCertsPath)
	if err != nil {
		return "", err
	}

	_, outputCertificate, err := certificate.HpcrGetEncryptionCertificateFromJson(encryptionCertsJson, version)
	if err != nil {
		return "", err
	}

	return outputCertificate, nil
}

func printCertificate(cert, certPath string) error {
	if certPath != "" {
		err := common.WriteDataToFile(certPath, cert)
		if err != nil {
			return err
		}
		fmt.Println("Successfully added encryption certificate ")
	} else {
		fmt.Println(cert)
	}

	return nil
}
