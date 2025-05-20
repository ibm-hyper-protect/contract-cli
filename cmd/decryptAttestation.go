// Copyright (c) 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"

	"github.com/ibm-hyper-protect/contract-go/attestation"
	"github.com/spf13/cobra"

	"github.com/ibm-hyper-protect/contract-cli/common"
)

const (
	successMessageDecryptAttestation = "Successfully decrypted attestation records"
)

// decryptAttestationCmd represents the decryptAttestation command
var decryptAttestationCmd = &cobra.Command{
	Use:   common.DecryptAttestParamName,
	Short: common.DecryptAttestParamShortDescription,
	Long:  common.DecryptAttestParamLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		encAttestPath, privateKeyPath, decryptedAttestPath, err := validateInputDecryptedAttestation(cmd)
		if err != nil {
			log.Fatal(err)
		}

		decryptedAttestationRecords, err := decryptAttestationRecords(encAttestPath, privateKeyPath)
		if err != nil {
			log.Fatal(err)
		}

		err = printDecryptAttestation(decryptedAttestationRecords, decryptedAttestPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(decryptAttestationCmd)

	decryptAttestationCmd.PersistentFlags().String(common.FileInFlagName, common.DecryptAttestFileInDefaultPath, common.DecryptAttestFileInDescription)
	decryptAttestationCmd.PersistentFlags().String(common.PrivateKeyFlagName, "", common.PrivateKeyFlagDescription)
	decryptAttestationCmd.PersistentFlags().String(common.FileOutFlagName, "", common.DecryptAttestFlagDescription)
}

func validateInputDecryptedAttestation(cmd *cobra.Command) (string, string, string, error) {
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

func decryptAttestationRecords(encryptedAttestationRecordsPath, privateKeyPath string) (string, error) {
	if !common.CheckFileFolderExists(encryptedAttestationRecordsPath) || !common.CheckFileFolderExists(privateKeyPath) {
		log.Fatal("The path to encrypted attestation records file or private key doesn't exists")
	}

	encryptedChecksum, err := common.ReadDataFromFile(encryptedAttestationRecordsPath)
	if err != nil {
		return "", err
	}

	privateKey, err := common.ReadDataFromFile(privateKeyPath)
	if err != nil {
		return "", err
	}

	decryptedAttestationRecords, err := attestation.HpcrGetAttestationRecords(encryptedChecksum, privateKey)
	if err != nil {
		return "", err
	}

	return decryptedAttestationRecords, nil
}

func printDecryptAttestation(decryptedAttestationRecords, decryptedAttestationPath string) error {
	if decryptedAttestationPath != "" {
		err := common.WriteDataToFile(decryptedAttestationPath, decryptedAttestationRecords)
		if err != nil {
			return err
		}
		fmt.Println(successMessageDecryptAttestation)
	} else {
		fmt.Println(decryptedAttestationRecords)
	}

	return nil
}
