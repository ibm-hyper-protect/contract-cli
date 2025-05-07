package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/image"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type ImageDetails struct {
	Id       string `json:"id"       yaml:"id"`
	Name     string `json:"name"     yaml:"name"`
	Checksum string `json:"checksum" yaml:"checksum"`
	Version  string `json:"version"  yaml:"version"`
}

const (
	FormatTypeJson = "json"
	FormatTypeYaml = "yaml"
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   common.ImageParamName,
	Short: common.ImageParamShortDescription,
	Long:  common.ImageParamLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		imageListJsonPath, err := cmd.Flags().GetString(common.FileInFlagName)
		if err != nil {
			log.Fatal()
		}
		versionName, err := cmd.Flags().GetString(common.VersionFlagName)
		if err != nil {
			log.Fatal()
		}
		formatType, err := cmd.Flags().GetString(common.DataFormatFlagName)
		if err != nil {
			log.Fatal(err)
		}
		hpcrImagePath, err := cmd.Flags().GetString(common.FileOutFlagName)
		if err != nil {
			log.Fatal(err)
		}

		imageDetail, err := GetImageDetails(imageListJsonPath, versionName)
		if err != nil {
			log.Fatal(err)
		}

		result, err := PrintData(imageDetail, formatType)
		if err != nil {
			log.Fatal(err)
		}

		err = common.WriteDataToFile(hpcrImagePath, result)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Successfully stored HPCR image details")
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)

	imageCmd.PersistentFlags().String(common.FileInFlagName, "", common.IbmCloudJsonInputDescription)
	imageCmd.PersistentFlags().String(common.VersionFlagName, "", common.HpcrVersionFlagDescription)
	imageCmd.PersistentFlags().String(common.DataFormatFlagName, common.DataFormatDefault, common.DataFormatFlagDescription)
	imageCmd.PersistentFlags().String(common.FileOutFlagName, "", common.HpcrImageFlagDescription)
}

func GetImageDetails(imageDetailsJsonPath, versionName string) (ImageDetails, error) {
	if !common.CheckFileFolderExists(imageDetailsJsonPath) {
		log.Fatal("The Image details path doesn't file")
	}

	imageDataJson, err := common.ReadDataFromFile(imageDetailsJsonPath)
	if err != nil {
		log.Fatal(err)
	}

	imageId, imageName, imageChecksum, ImageVersion, err := image.HpcrSelectImage(imageDataJson, versionName)
	if err != nil {
		log.Fatal(err)
	}

	return ImageDetails{imageId, imageName, imageChecksum, ImageVersion}, nil
}

func PrintData(imageDetail ImageDetails, format string) (string, error) {
	if format == FormatTypeJson {
		imageJson, err := json.MarshalIndent(imageDetail, "", "  ")
		if err != nil {
			return "", err
		}
		return string(imageJson), nil
	} else if format == FormatTypeYaml {
		imageYaml, err := yaml.Marshal(imageDetail)
		if err != nil {
			return "", err
		}
		return string(imageYaml), nil
	} else {
		return "", fmt.Errorf("invalid output format")
	}
}
