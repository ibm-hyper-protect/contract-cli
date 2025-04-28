package common

import (
	"fmt"
	"io"
	"os"
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
