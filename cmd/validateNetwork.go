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

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-cli/lib/validateNetwork"
	"github.com/ibm-hyper-protect/contract-go/v2/network"
	"github.com/spf13/cobra"
)

// validateNetworkConfigCmd represents the validate-networkConfig command
var validateNetworkConfigCmd = &cobra.Command{
	Use:   validateNetwork.ParameterName,
	Short: validateNetwork.ParameterShortDescription,
	Long:  validateNetwork.ParameterLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		networkConfigPath, err := validateNetwork.ValidateInput(cmd)
		if err != nil {
			log.Fatal(err)
		}

		if !common.CheckFileFolderExists(networkConfigPath) {
			log.Fatal("The path to network-config doesn't exist")
		}

		networkConfigData, err := common.ReadDataFromFile(networkConfigPath)
		if err != nil {
			log.Fatal(err)
		}

		err = network.HpcrVerifyNetworkConfig(networkConfigData)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("network-config validated successfully")
	},
}

// init - cobra init function
func init() {
	rootCmd.AddCommand(validateNetworkConfigCmd)

	requiredFlags := map[string]bool{
		"in": true,
	}
	validateNetworkConfigCmd.PersistentFlags().String(validateNetwork.InputFlagName, "", validateNetwork.InputFlagDescription)
	common.SetCustomHelpTemplate(validateNetworkConfigCmd, requiredFlags)
}
