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
	cryptoLib "github.com/ibm-hyper-protect/contract-cli/lib/crypto"
	"github.com/spf13/cobra"
)

// cryptoCmd represents the openssl-keycert command
var cryptoCmd = &cobra.Command{
	Use:   cryptoLib.ParameterName,
	Short: cryptoLib.ParameterShortDescription,
	Long:  cryptoLib.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		artifactType, password, commonName, sans,
			privPath, pubPath, caPath, certPath, certKeyPath,
			keySize, validDays, err := cryptoLib.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		err = cryptoLib.Generate(artifactType, password, commonName, sans,
			privPath, pubPath, caPath, certPath, certKeyPath,
			keySize, validDays)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init registers cryptoCmd with the root command and declares all flags.
func init() {
	rootCmd.AddCommand(cryptoCmd)

	requiredFlags := map[string]bool{
		cryptoLib.TypeFlagName: true,
	}

	cryptoCmd.PersistentFlags().String(cryptoLib.TypeFlagName, "", cryptoLib.TypeFlagDescription)
	cryptoCmd.PersistentFlags().String(cryptoLib.OutFlagName, "", cryptoLib.OutFlagDescription)
	cryptoCmd.PersistentFlags().Int(cryptoLib.SizeFlagName, 0, cryptoLib.SizeFlagDescription)
	cryptoCmd.PersistentFlags().String(cryptoLib.PasswordFlagName, "", cryptoLib.PasswordFlagDescription)
	cryptoCmd.PersistentFlags().String(cryptoLib.SANFlagName, "", cryptoLib.SANFlagDescription)
	cryptoCmd.PersistentFlags().Int(cryptoLib.DaysFlagName, 0, cryptoLib.DaysFlagDescription)
	cryptoCmd.PersistentFlags().String(cryptoLib.CNFlagName, "", cryptoLib.CNFlagDescription)

	common.SetCustomHelpTemplate(cryptoCmd, requiredFlags)
	common.SetCustomErrorTemplate(cryptoCmd)
}
