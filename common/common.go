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

package common

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CheckFileFolderExists - function to check if file or folder exists
func CheckFileFolderExists(folderFilePath string) bool {
	_, err := os.Stat(folderFilePath)
	return !os.IsNotExist(err)
}

// ReadDataFromFile - function to read data from file
func ReadDataFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// WriteDataToFile - function to write data to file (create file if doesn't exists)
func WriteDataToFile(filePath, data string) error {
	DataFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file - %v", err)
	}
	defer DataFile.Close()

	_, err = DataFile.WriteString(data)
	if err != nil {
		return fmt.Errorf("failed to write data - %v", err)
	}

	return nil
}

// ExecCommand - function to run os commands
func ExecCommand(name string, stdinInput string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	// Check for standard input
	if stdinInput != "" {
		stdinPipe, err := cmd.StdinPipe()
		if err != nil {
			return "", err
		}
		defer stdinPipe.Close()

		go func() {
			defer stdinPipe.Close()
			stdinPipe.Write([]byte(stdinInput))
		}()
	}

	// Buffer to capture the output from the command.
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command.
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Return the output from the command and nil for the error.
	return out.String(), nil
}

// OpensslCheck - function to check if openssl exists
func OpensslCheck() error {
	_, err := ExecCommand("openssl", "", "version")

	if err != nil {
		return err
	}

	return nil
}

// GetPrivateKey - function to get private key from path or generate a new private key
func GetPrivateKey(privateKeyPath string) (string, error) {
	var privateKey string
	var err error

	if privateKeyPath == "" {
		privateKey, err = generatePrivateKey()
		if err != nil {
			return "", err
		}
	} else {
		if CheckFileFolderExists(privateKeyPath) {
			privateKey, err = ReadDataFromFile(privateKeyPath)
			if err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("private key path doesn't exist")
		}
	}

	return privateKey, nil
}

// generatePrivateKey - function to generate private key
func generatePrivateKey() (string, error) {
	err := OpensslCheck()
	if err != nil {
		return "", fmt.Errorf("openssl not found - %v", err)
	}

	privateKey, err := ExecCommand("openssl", "", "genrsa", "4096")
	if err != nil {
		return "", fmt.Errorf("failed to generate private key - %v", err)
	}

	return privateKey, nil
}

// GetDataFromFile - function to get data from file
func GetDataFromFile(certPath string) (string, error) {
	var encCert string
	var err error

	if certPath != "" {
		encCert, err = ReadDataFromFile(certPath)
		if err != nil {
			return "", err
		}
	} else {
		encCert = ""
	}

	return encCert, nil
}

// SetCustomHelpTemplate - funtion to format contract-cli --help command
func SetCustomHelpTemplate(cmd *cobra.Command, requiredFlags map[string]bool) {
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		out := cmd.OutOrStdout()

		fmt.Fprintln(out, "Usage:")
		fmt.Fprintf(out, "  %s [flags]\n\n", cmd.CommandPath())

		maxLen := 0
		cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
			if len(f.Name) > maxLen {
				maxLen = len(f.Name)
			}
		})

		printFlags := func(mandatory bool) {
			cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
				if requiredFlags[f.Name] != mandatory {
					return
				}
				padding := maxLen - len(f.Name) + 2
				fmt.Fprintf(out, "  --%s%s%s\n", f.Name, strings.Repeat(" ", padding), f.Usage)
			})
		}

		fmt.Fprintln(out, "Mandatory Flags:")
		printFlags(true)

		fmt.Fprintln(out, "\nOptional Flags:")
		printFlags(false)
	})
}

// SetCustomErrorTemplate function customizes Cobra's default flag error handling to
// print a spaced error message followed by command usage.
func SetCustomErrorTemplate(cmd *cobra.Command) {
	cmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		SetMandatoryFlagError(cmd, err)
		return nil
	})

	cmd.SilenceUsage = true
}

// SetMandatoryFlagError function print error if required flag is missing
func SetMandatoryFlagError(cmd *cobra.Command, err error) {
	out := cmd.ErrOrStderr()
	fmt.Fprintln(out, err)
	fmt.Fprintln(out)
	cmd.Usage()
	os.Exit(1)
}
