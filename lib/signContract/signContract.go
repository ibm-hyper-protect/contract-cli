package signContract

import (
	"fmt"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/v2/contract"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "sign-contract"
	ParameterShortDescription = "Sign an encrypted contract"
	ParameterLongDescription  = `Sign an contract with encrypted workload and env`
	InputFlagName             = "in"
	InputFlagDescription      = "Path to encrypted contract"
	PrivateKeyFlagName        = "priv"
	PrivateKeyFlagDescription = "Path to private key file for signing"
	PasswordFlagName          = "password"
	PasswordFlagDescription   = "Password for encrypted private key"
	OutputFlagName            = "out"
	OutputFlagDescription     = "Path to save encrypted output"
)

// ValidateInput - function to validate inputs of sign-contract
func ValidateInput(cmd *cobra.Command) (string, string, string, string, error) {
	inputData, err := cmd.Flags().GetString(InputFlagName)
	if err != nil {
		return "", "", "", "", err
	}

	privateKeyPath, err := cmd.Flags().GetString(PrivateKeyFlagName)
	if err != nil {
		return "", "", "", "", err
	}

	if inputData == "" || privateKeyPath == "" {
		err := fmt.Errorf("Error: required flag '--in' or '--priv' is missing")
		common.SetMandatoryFlagError(cmd, err)
	}

	// Validate stdin input conflicts
	common.ValidateStdinInput(cmd, inputData)

	outputPath, err := cmd.Flags().GetString(OutputFlagName)
	if err != nil {
		return "", "", "", "", err
	}

	password, err := cmd.Flags().GetString(PasswordFlagName)
	if err != nil {
		return "", "", "", "", err
	}

	return inputData, privateKeyPath, outputPath, password, nil
}

func GenerateSignContract(inputDataPath, privateKeyPath, password string) (string, error) {
	var inputData string
	var err error

	if inputDataPath == "-" {
		inputData, err = common.ReadDataFromStdin()
		if err != nil {
			return "", fmt.Errorf("unable to read input from standard input: %w", err)
		}
	} else {
		if !common.CheckFileFolderExists(inputDataPath) {
			return "", fmt.Errorf("the contract path doesn't exist")
		}

		inputData, err = common.ReadDataFromFile(inputDataPath)
		if err != nil {
			return "", err
		}
	}

	privateKey, err := common.GetPrivateKey(privateKeyPath)
	if err != nil {
		return "", err
	}

	signedContract, _, _, err := contract.HpcrContractSign(inputData, privateKey, password)
	if err != nil {
		return "", err
	}

	return signedContract, nil
}

// Output - function to print signed contract or redirect it to a file
func Output(signedContract, outputPath string) error {
	if outputPath != "" {
		err := common.WriteDataToFile(outputPath, signedContract)
		if err != nil {
			return err
		}
		fmt.Println("Successfully generated signed contract")
	} else {
		fmt.Println(signedContract)
	}

	return nil
}
