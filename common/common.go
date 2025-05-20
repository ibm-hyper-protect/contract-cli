package common

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
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
