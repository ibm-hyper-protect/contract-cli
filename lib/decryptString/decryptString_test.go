// Copyright (c) 2026 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0

package decryptString

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestProcessMissingPrivKey - missing private key path returns error
func TestProcessMissingPrivKey(t *testing.T) {
	_, err := Process("hyper-protect-basic.abc.def", "", "")
	assert.Error(t, err)
}

// TestProcessMissingInput - empty input returns error
func TestProcessMissingInput(t *testing.T) {
	_, err := Process("", "/some/key.pem", "")
	assert.Error(t, err)
}

// TestProcessWrongEncryptionFormat - wrong encryption format returns error
func TestProcessWrongEncryptionFormat(t *testing.T) {
	_, err := Process("invalid-format.wrongpassword.wrongdata", "../../samples/decrypt/private.key", "")
	assert.Error(t, err)
}

// TestOutputToStdout - output with no file path succeeds
func TestOutputToStdout(t *testing.T) {
	err := Output("", "decrypted text")
	assert.NoError(t, err)
}

// TestOutputToFile - output to file succeeds
func TestOutputToFile(t *testing.T) {
	tmpFile := t.TempDir() + "/out.txt"
	err := Output(tmpFile, "decrypted text")
	assert.NoError(t, err)
}
