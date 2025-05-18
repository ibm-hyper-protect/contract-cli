package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	sampleFile = "../samples/attestation/se-checksums.txt.enc"

	simpleSampleTextPath = "../samples/simple_file.txt"
	simpleSampleText     = "Testing"

	simpleSampleWritePath = "../build/simple_file.txt"
)

func TestCheckFileFolderExists(t *testing.T) {
	result := CheckFileFolderExists(sampleFile)

	assert.True(t, result)
}

func TestReadDataFromFile(t *testing.T) {
	content, err := ReadDataFromFile(simpleSampleTextPath)
	if err != nil {
		t.Errorf("failed to read text from file - %v", err)
	}

	assert.Equal(t, content, simpleSampleText)
}

func TestWriteDataToFile(t *testing.T) {
	err := WriteDataToFile(simpleSampleWritePath, simpleSampleText)
	assert.NoError(t, err)
}

func TestExecCommand(t *testing.T) {
	_, err := ExecCommand("openssl", "", "version")
	if err != nil {
		t.Errorf("failed to execute command - %v", err)
	}
}

func TestOpensslCheck(t *testing.T) {
	err := OpensslCheck()
	if err != nil {
		t.Errorf("openssl check failed - %v", err)
	}
}
