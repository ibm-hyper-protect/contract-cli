package cmd

import (
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/certificate"
	"github.com/spf13/cobra"
)

var (
	versions []string

	downloadCertificatesCmd = &cobra.Command{
		Use:   common.DownloadCertParamName,
		Short: common.DownloadCertParamShortDescription,
		Long:  common.DownloadCertParamLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			formatType, certificatePath, err := ValidateInputDownloadCertificates(cmd)
			if err != nil {
				log.Fatal(err)
			}

			certificates, err := certificate.HpcrDownloadEncryptionCertificates(versions, formatType, "")
			if err != nil {
				log.Fatal(err)
			}

			if certificatePath != "" {
				err := common.WriteDataToFile(certificatePath, certificates)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Successfully stored certificates")
			} else {
				fmt.Println(certificates)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(downloadCertificatesCmd)

	downloadCertificatesCmd.PersistentFlags().StringSliceVarP(&versions, common.VersionFlagName, "", []string{}, common.EncryptionCertVersionFlagDescription)
	downloadCertificatesCmd.PersistentFlags().String(common.DataFormatFlagName, common.DataFormatDefault, common.DataFormatFlagDescription)
	downloadCertificatesCmd.PersistentFlags().String(common.FileOutFlagName, "", common.EncryptionCertsFlagDescription)
}

func ValidateInputDownloadCertificates(cmd *cobra.Command) (string, string, error) {
	formatType, err := cmd.Flags().GetString(common.DataFormatFlagName)
	if err != nil {
		return "", "", err
	}

	certificatePath, err := cmd.Flags().GetString(common.FileOutFlagName)
	if err != nil {
		return "", "", err
	}

	return formatType, certificatePath, nil
}
