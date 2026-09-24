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

package crypto

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureStdout redirects os.Stdout to a temp file for the duration of fn,
// then restores it and returns whatever was written.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	tmpFile, err := os.CreateTemp(t.TempDir(), "stdout-*")
	require.NoError(t, err)

	orig := os.Stdout
	os.Stdout = tmpFile
	defer func() {
		os.Stdout = orig
		tmpFile.Close()
	}()

	fn()

	_, err = tmpFile.Seek(0, 0)
	require.NoError(t, err)
	data, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)
	return string(data)
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// newCmd builds a minimal cobra.Command with all crypto flags registered,
// then parses the provided flag values into it.
func newCmd(t *testing.T, flags map[string]string) *cobra.Command {
	t.Helper()
	cmd := &cobra.Command{Use: ParameterName}
	cmd.PersistentFlags().String(TypeFlagName, "", TypeFlagDescription)
	cmd.PersistentFlags().String(OutFlagName, "", OutFlagDescription)
	cmd.PersistentFlags().Int(SizeFlagName, 0, SizeFlagDescription)
	cmd.PersistentFlags().String(PasswordFlagName, "", PasswordFlagDescription)
	cmd.PersistentFlags().String(SANFlagName, "", SANFlagDescription)
	cmd.PersistentFlags().Int(DaysFlagName, 0, DaysFlagDescription)
	cmd.PersistentFlags().String(CNFlagName, "", CNFlagDescription)

	args := make([]string, 0, len(flags)*2)
	for k, v := range flags {
		args = append(args, "--"+k, v)
	}
	require.NoError(t, cmd.ParseFlags(args))
	return cmd
}

// ---------------------------------------------------------------------------
// resolveKeyPaths unit tests
// ---------------------------------------------------------------------------

// TestResolveKeyPaths_Empty verifies stdout mode when outRaw is empty.
func TestResolveKeyPaths_Empty(t *testing.T) {
	priv, pub := resolveKeyPaths("")
	assert.Equal(t, "", priv)
	assert.Equal(t, "", pub)
}

// TestResolveKeyPaths_OneToken verifies that a single name gets the fixed
// default public key name and that a log line is printed.
func TestResolveKeyPaths_OneToken(t *testing.T) {
	var priv, pub string
	out := captureStdout(t, func() {
		priv, pub = resolveKeyPaths("mykey.pem")
	})
	assert.Equal(t, "mykey.pem", priv)
	assert.Equal(t, DefaultKeyPublicName, pub)
	assert.Contains(t, out, DefaultKeyPublicName)
}

// TestResolveKeyPaths_TwoTokens verifies both names are used as-is.
func TestResolveKeyPaths_TwoTokens(t *testing.T) {
	var priv, pub string
	captureStdout(t, func() {
		priv, pub = resolveKeyPaths("priv.pem,pub.pem")
	})
	assert.Equal(t, "priv.pem", priv)
	assert.Equal(t, "pub.pem", pub)
}

// TestResolveKeyPaths_ExcessTokens verifies that 3+ tokens are truncated to 2
// and a warning is printed.
func TestResolveKeyPaths_ExcessTokens(t *testing.T) {
	var priv, pub string
	out := captureStdout(t, func() {
		priv, pub = resolveKeyPaths("a.pem,b.pem,c.pem")
	})
	assert.Equal(t, "a.pem", priv)
	assert.Equal(t, "b.pem", pub)
	assert.Contains(t, out, "Warning")
	assert.Contains(t, out, "ignoring extra values")
}

// ---------------------------------------------------------------------------
// resolveCertPaths unit tests
// ---------------------------------------------------------------------------

// TestResolveCertPaths_Empty verifies stdout mode when outRaw is empty.
func TestResolveCertPaths_Empty(t *testing.T) {
	ca, cert, key := resolveCertPaths("")
	assert.Equal(t, "", ca)
	assert.Equal(t, "", cert)
	assert.Equal(t, "", key)
}

// TestResolveCertPaths_OneToken verifies that missing slots 2 and 3 receive
// their fixed defaults and that log lines are printed for both.
func TestResolveCertPaths_OneToken(t *testing.T) {
	var ca, cert, key string
	out := captureStdout(t, func() {
		ca, cert, key = resolveCertPaths("ca.crt")
	})
	assert.Equal(t, "ca.crt", ca)
	assert.Equal(t, DefaultCertClientName, cert)
	assert.Equal(t, DefaultCertKeyName, key)
	assert.Contains(t, out, DefaultCertClientName)
	assert.Contains(t, out, DefaultCertKeyName)
}

// TestResolveCertPaths_TwoTokens verifies that missing slot 3 gets its default.
func TestResolveCertPaths_TwoTokens(t *testing.T) {
	var ca, cert, key string
	out := captureStdout(t, func() {
		ca, cert, key = resolveCertPaths("ca.crt,client.crt")
	})
	assert.Equal(t, "ca.crt", ca)
	assert.Equal(t, "client.crt", cert)
	assert.Equal(t, DefaultCertKeyName, key)
	assert.Contains(t, out, DefaultCertKeyName)
}

// TestResolveCertPaths_ThreeTokens verifies all three names are used as-is.
func TestResolveCertPaths_ThreeTokens(t *testing.T) {
	var ca, cert, key string
	captureStdout(t, func() {
		ca, cert, key = resolveCertPaths("ca.crt,client.crt,client.pem")
	})
	assert.Equal(t, "ca.crt", ca)
	assert.Equal(t, "client.crt", cert)
	assert.Equal(t, "client.pem", key)
}

// TestResolveCertPaths_ExcessTokens verifies that 4+ tokens are truncated to 3
// and a warning is printed.
func TestResolveCertPaths_ExcessTokens(t *testing.T) {
	var ca, cert, key string
	out := captureStdout(t, func() {
		ca, cert, key = resolveCertPaths("a.crt,b.crt,c.pem,d.pem")
	})
	assert.Equal(t, "a.crt", ca)
	assert.Equal(t, "b.crt", cert)
	assert.Equal(t, "c.pem", key)
	assert.Contains(t, out, "Warning")
	assert.Contains(t, out, "ignoring extra values")
}

// ---------------------------------------------------------------------------
// ValidateInput
// ---------------------------------------------------------------------------

// TestValidateInput_KeyDefaults verifies that omitted optional flags receive
// correct defaults for --type key, and that paths are empty (stdout mode).
func TestValidateInput_KeyDefaults(t *testing.T) {
	cmd := newCmd(t, map[string]string{
		TypeFlagName: "key",
	})
	artifactType, password, _, _, privPath, pubPath, _, _, _, keySize, _, err := ValidateInput(cmd)
	require.NoError(t, err)

	assert.Equal(t, "key", artifactType)
	assert.Equal(t, DefaultKeySize, keySize)
	assert.Equal(t, "", password)
	assert.Equal(t, "", privPath)
	assert.Equal(t, "", pubPath)
}

// TestValidateInput_CertDefaults verifies that omitted optional flags receive
// correct defaults for --type cert, including CN derived from the first SAN.
func TestValidateInput_CertDefaults(t *testing.T) {
	cmd := newCmd(t, map[string]string{
		TypeFlagName: "cert",
	})
	artifactType, _, commonName, sans, _, _, caPath, certPath, certKeyPath, keySize, validDays, err := ValidateInput(cmd)
	require.NoError(t, err)

	assert.Equal(t, "cert", artifactType)
	assert.Equal(t, DefaultKeySize, keySize)
	assert.Equal(t, DefaultSAN, sans)
	assert.Equal(t, DefaultValidDays, validDays)
	assert.Equal(t, strings.SplitN(DefaultSAN, ",", 2)[0], commonName)
	// all paths empty = stdout mode
	assert.Equal(t, "", caPath)
	assert.Equal(t, "", certPath)
	assert.Equal(t, "", certKeyPath)
}

// TestValidateInput_KeyOutOneToken verifies that a single --out token for
// --type key resolves privPath to that token and pubPath to the default.
// The private key path is rooted inside t.TempDir() to prevent any accidental
// file creation in the package source directory.
func TestValidateInput_KeyOutOneToken(t *testing.T) {
	tmpDir := t.TempDir()
	privName := filepath.Join(tmpDir, "mykey.pem")
	cmd := newCmd(t, map[string]string{
		TypeFlagName: "key",
		OutFlagName:  privName,
	})
	var privPath, pubPath string
	captureStdout(t, func() {
		_, _, _, _, privPath, pubPath, _, _, _, _, _, _ = ValidateInput(cmd)
	})
	assert.Equal(t, privName, privPath)
	assert.Equal(t, DefaultKeyPublicName, pubPath)
}

// TestValidateInput_KeyOutTwoTokens verifies that two --out tokens for
// --type key resolve to exactly those paths.
func TestValidateInput_KeyOutTwoTokens(t *testing.T) {
	cmd := newCmd(t, map[string]string{
		TypeFlagName: "key",
		OutFlagName:  "priv.pem,pub.pem",
	})
	_, _, _, _, privPath, pubPath, _, _, _, _, _, err := ValidateInput(cmd)
	require.NoError(t, err)
	assert.Equal(t, "priv.pem", privPath)
	assert.Equal(t, "pub.pem", pubPath)
}

// TestValidateInput_CertOutOneToken verifies that a single --out token for
// --type cert resolves caPath to that token and the remaining to defaults.
// The CA cert path is rooted inside t.TempDir() to prevent any accidental
// file creation in the package source directory.
func TestValidateInput_CertOutOneToken(t *testing.T) {
	tmpDir := t.TempDir()
	caName := filepath.Join(tmpDir, "ca.crt")
	cmd := newCmd(t, map[string]string{
		TypeFlagName: "cert",
		OutFlagName:  caName,
	})
	var caPath, certPath, certKeyPath string
	captureStdout(t, func() {
		_, _, _, _, _, _, caPath, certPath, certKeyPath, _, _, _ = ValidateInput(cmd)
	})
	assert.Equal(t, caName, caPath)
	assert.Equal(t, DefaultCertClientName, certPath)
	assert.Equal(t, DefaultCertKeyName, certKeyPath)
}

// TestValidateInput_CertOutThreeTokens verifies that three --out tokens for
// --type cert resolve to exactly those paths.
func TestValidateInput_CertOutThreeTokens(t *testing.T) {
	cmd := newCmd(t, map[string]string{
		TypeFlagName: "cert",
		OutFlagName:  "ca.crt,client.crt,client.pem",
	})
	_, _, _, _, _, _, caPath, certPath, certKeyPath, _, _, err := ValidateInput(cmd)
	require.NoError(t, err)
	assert.Equal(t, "ca.crt", caPath)
	assert.Equal(t, "client.crt", certPath)
	assert.Equal(t, "client.pem", certKeyPath)
}

// TestValidateInput_AllFlagsExplicit verifies that explicitly provided flags
// override all defaults (cert type, 3 out tokens).
func TestValidateInput_AllFlagsExplicit(t *testing.T) {
	cmd := newCmd(t, map[string]string{
		TypeFlagName:     "cert",
		OutFlagName:      "ca.crt,client.crt,key.pem",
		SizeFlagName:     "2048",
		PasswordFlagName: "s3cr3t",
		SANFlagName:      "foo.com,bar.com",
		DaysFlagName:     "90",
		CNFlagName:       "foo.com",
	})
	artifactType, password, commonName, sans, _, _, caPath, certPath, certKeyPath, keySize, validDays, err := ValidateInput(cmd)
	require.NoError(t, err)

	assert.Equal(t, "cert", artifactType)
	assert.Equal(t, 2048, keySize)
	assert.Equal(t, "s3cr3t", password)
	assert.Equal(t, "foo.com,bar.com", sans)
	assert.Equal(t, 90, validDays)
	assert.Equal(t, "foo.com", commonName)
	assert.Equal(t, "ca.crt", caPath)
	assert.Equal(t, "client.crt", certPath)
	assert.Equal(t, "key.pem", certKeyPath)
}

// TestValidateInput_CNDefaultsToFirstSAN verifies CN defaults to first SAN token.
func TestValidateInput_CNDefaultsToFirstSAN(t *testing.T) {
	cmd := newCmd(t, map[string]string{
		TypeFlagName: "cert",
		SANFlagName:  "first.example.com,second.example.com",
	})
	_, _, commonName, _, _, _, _, _, _, _, _, err := ValidateInput(cmd)
	require.NoError(t, err)
	assert.Equal(t, "first.example.com", commonName)
}

// TestValidateInput_KeyWithPassword verifies password is forwarded unchanged.
func TestValidateInput_KeyWithPassword(t *testing.T) {
	cmd := newCmd(t, map[string]string{
		TypeFlagName:     "key",
		PasswordFlagName: "my-passphrase",
	})
	_, password, _, _, _, _, _, _, _, _, _, err := ValidateInput(cmd)
	require.NoError(t, err)
	assert.Equal(t, "my-passphrase", password)
}

// ---------------------------------------------------------------------------
// Generate — stdout path (all paths empty)
// ---------------------------------------------------------------------------

// TestGenerate_KeyToStdout verifies that Generate writes valid PEM to stdout
// when all paths are empty.
func TestGenerate_KeyToStdout(t *testing.T) {
	var err error
	out := captureStdout(t, func() {
		err = Generate("key", "", "", "", "", "", "", "", "", 2048, 0)
	})
	require.NoError(t, err)
	assert.Contains(t, out, "PRIVATE KEY")
	assert.Contains(t, out, "PUBLIC KEY")
}

// TestGenerate_CertToStdout verifies that Generate writes valid PEM to stdout
// for --type cert when all paths are empty.
func TestGenerate_CertToStdout(t *testing.T) {
	var err error
	out := captureStdout(t, func() {
		err = Generate("cert", "", "example.com", "example.com", "", "", "", "", "", 2048, 365)
	})
	require.NoError(t, err)
	assert.Contains(t, out, "BEGIN CERTIFICATE")
	assert.Contains(t, out, "PRIVATE KEY")
}

// ---------------------------------------------------------------------------
// Generate — file path (explicit names)
// ---------------------------------------------------------------------------

// TestGenerate_KeyToFiles_TwoNames verifies that --type key writes to exactly
// the two provided filenames.
func TestGenerate_KeyToFiles_TwoNames(t *testing.T) {
	tmpDir := t.TempDir()
	privPath := filepath.Join(tmpDir, "mykey.pem")
	pubPath := filepath.Join(tmpDir, "mykey_pub.pem")

	var err error
	captureStdout(t, func() {
		err = Generate("key", "", "", "", privPath, pubPath, "", "", "", 2048, 0)
	})
	require.NoError(t, err)

	assert.FileExists(t, privPath)
	assert.FileExists(t, pubPath)

	privData, err := os.ReadFile(privPath)
	require.NoError(t, err)
	assert.Contains(t, string(privData), "PRIVATE KEY")

	pubData, err := os.ReadFile(pubPath)
	require.NoError(t, err)
	assert.Contains(t, string(pubData), "PUBLIC KEY")

	info, err := os.Stat(privPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm(),
		"private key must have 0600 permissions")
}

// TestGenerate_KeyToFiles_DefaultPublic verifies the default public key name
// is used when only the private key path is provided via resolveKeyPaths.
func TestGenerate_KeyToFiles_DefaultPublic(t *testing.T) {
	tmpDir := t.TempDir()
	privPath := filepath.Join(tmpDir, "priv.pem")

	// Simulate single-token --out: public key falls back to DefaultKeyPublicName
	// placed in the working directory; use tmpDir as working dir via absolute path.
	pubPath := filepath.Join(tmpDir, DefaultKeyPublicName)

	var err error
	captureStdout(t, func() {
		err = Generate("key", "", "", "", privPath, pubPath, "", "", "", 2048, 0)
	})
	require.NoError(t, err)

	assert.FileExists(t, privPath)
	assert.FileExists(t, pubPath)
}

// TestGenerate_CertToFiles_ThreeNames verifies that --type cert writes to
// exactly the three provided filenames.
func TestGenerate_CertToFiles_ThreeNames(t *testing.T) {
	tmpDir := t.TempDir()
	caPath := filepath.Join(tmpDir, "ca.crt")
	certPath := filepath.Join(tmpDir, "client.crt")
	keyPath := filepath.Join(tmpDir, "client.pem")

	var err error
	captureStdout(t, func() {
		err = Generate("cert", "", "example.com", "example.com", "", "", caPath, certPath, keyPath, 2048, 365)
	})
	require.NoError(t, err)

	assert.FileExists(t, caPath)
	assert.FileExists(t, certPath)
	assert.FileExists(t, keyPath)

	caData, err := os.ReadFile(caPath)
	require.NoError(t, err)
	assert.Contains(t, string(caData), "BEGIN CERTIFICATE")

	clientCertData, err := os.ReadFile(certPath)
	require.NoError(t, err)
	assert.Contains(t, string(clientCertData), "BEGIN CERTIFICATE")

	clientKeyData, err := os.ReadFile(keyPath)
	require.NoError(t, err)
	assert.Contains(t, string(clientKeyData), "PRIVATE KEY")

	info, err := os.Stat(keyPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm(),
		"client private key must have 0600 permissions")
}

// TestGenerate_KeyWithPasswordToFiles verifies that a password-protected key
// file is written encrypted.
func TestGenerate_KeyWithPasswordToFiles(t *testing.T) {
	tmpDir := t.TempDir()
	privPath := filepath.Join(tmpDir, "enckey.pem")
	pubPath := filepath.Join(tmpDir, "enckey_pub.pem")

	var err error
	captureStdout(t, func() {
		err = Generate("key", "test-passphrase", "", "", privPath, pubPath, "", "", "", 2048, 0)
	})
	require.NoError(t, err)

	privData, err := os.ReadFile(privPath)
	require.NoError(t, err)
	assert.Contains(t, string(privData), "ENCRYPTED",
		"private key written to file must be encrypted when --password is set")
}

// ---------------------------------------------------------------------------
// ValidateInput — --days for --type key (rejected)
// ---------------------------------------------------------------------------

// TestValidateInput_KeyDaysDefault verifies that --days defaults to 0 (not
// set) when --type key is used without an explicit --days flag.
func TestValidateInput_KeyDaysDefault(t *testing.T) {
	cmd := newCmd(t, map[string]string{
		TypeFlagName: "key",
	})
	_, _, _, _, _, _, _, _, _, _, validDays, err := ValidateInput(cmd)
	require.NoError(t, err)
	assert.Equal(t, 0, validDays, "--days should be 0 for --type key when not provided")
}

// TestValidateInput_KeyDaysRejected verifies that passing --days with --type
// key returns an error, since key pairs do not carry an expiry.
func TestValidateInput_KeyDaysRejected(t *testing.T) {
	cmd := newCmd(t, map[string]string{
		TypeFlagName: "key",
		DaysFlagName: "180",
	})
	_, _, _, _, _, _, _, _, _, _, _, err := ValidateInput(cmd)
	require.Error(t, err, "--days with --type key must produce an error")
	assert.Contains(t, err.Error(), "--days is not supported for --type key")
}

// ---------------------------------------------------------------------------
// ValidateInput — negative / invalid input test cases
// ---------------------------------------------------------------------------

// TestValidateInput_TypeValidation verifies that only "key" and "cert" are
// accepted as valid --type values.
//
// Note: passing a truly invalid type through ValidateInput is not testable here
// because the invalid-type path calls common.SetMandatoryFlagError → os.Exit(1),
// which would terminate the test process. Following the project convention
// (see cmd/signContract_test.go), we test the acceptance boundary instead:
// confirm "key" and "cert" pass, and document that other values exit via the
// SetMandatoryFlagError path in cmd/crypto.go.
func TestValidateInput_TypeValidation(t *testing.T) {
	for _, validType := range []string{"key", "cert"} {
		cmd := newCmd(t, map[string]string{
			TypeFlagName: validType,
		})
		artifactType, _, _, _, _, _, _, _, _, _, _, err := ValidateInput(cmd)
		require.NoError(t, err)
		assert.Equal(t, validType, artifactType,
			"valid --type %q must be accepted without error", validType)
	}
}

// TestValidateInput_NegativeDays verifies that a negative --days value is
// rejected with a descriptive error.
func TestValidateInput_NegativeDays(t *testing.T) {
	cmd := newCmd(t, map[string]string{
		TypeFlagName: "cert",
		DaysFlagName: "-5",
	})
	_, _, _, _, _, _, _, _, _, _, _, err := ValidateInput(cmd)
	require.Error(t, err, "a negative --days value must produce an error")
	assert.Contains(t, err.Error(), "--days",
		"error message should reference the --days flag")
	assert.Contains(t, err.Error(), "-5",
		"error message should include the bad value")
}
