// Copyright (c) 2026 IBM Corp.
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
	"log"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-cli/lib/sealedSecret"
	"github.com/spf13/cobra"
)

// sealedSecretCmd represents the sealed-secret command
var sealedSecretCmd = &cobra.Command{
	Use:   sealedSecret.ParameterName,
	Short: sealedSecret.ParameterShortDescription,
	Long:  sealedSecret.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		inputData, secretType, outputPath, encryptionKeyPath, signingKeyPath, err := sealedSecret.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		sealedSecretData, decryptionKey, verificationKey, _, _, err := sealedSecret.GenerateSealedSecret(inputData, secretType, encryptionKeyPath, signingKeyPath)
		if err != nil {
			log.Fatal(err)
		}

		err = sealedSecret.Output(sealedSecretData, decryptionKey, verificationKey, outputPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(sealedSecretCmd)

	requiredFlags := map[string]bool{
		"in":   true,
		"type": true,
	}

	sealedSecretCmd.PersistentFlags().String(sealedSecret.InputFlagName, "", sealedSecret.InputFlagDescription)
	sealedSecretCmd.PersistentFlags().String(sealedSecret.TypeFlagName, "", sealedSecret.TypeFlagDescription)
	sealedSecretCmd.PersistentFlags().String(sealedSecret.OutputFlagName, "", sealedSecret.OutputFlagDescription)
	sealedSecretCmd.PersistentFlags().String(sealedSecret.EncryptionKeyFlagName, "", sealedSecret.EncryptionKeyFlagDescription)
	sealedSecretCmd.PersistentFlags().String(sealedSecret.SigningKeyFlagName, "", sealedSecret.SigningKeyFlagDescription)

	common.SetCustomHelpTemplate(sealedSecretCmd, requiredFlags)
	common.SetCustomErrorTemplate(sealedSecretCmd)
}
