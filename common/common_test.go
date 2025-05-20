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

	sampleCertPath = "../samples/contract-expiry/personal_ca.crt"
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
	if err != nil {
		t.Errorf("failed to write data to file - %v", err)
	}
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

func TestGetPrivateKey(t *testing.T) {
	result, err := GetPrivateKey("")
	if err != nil {
		t.Errorf("failed to get private key - %v", err)
	}

	assert.NotEmpty(t, result)
}

func TestGeneratePrivateKey(t *testing.T) {
	result, err := generatePrivateKey()
	if err != nil {
		t.Errorf("failed to generate private key - %v", err)
	}

	assert.NotEmpty(t, result)
}

func TestGetDataFromFile(t *testing.T) {
	result, err := GetDataFromFile(sampleCertPath)
	if err != nil {
		t.Errorf("failed to get data from file - %v", err)
	}

	assert.NotEmpty(t, result)
}
