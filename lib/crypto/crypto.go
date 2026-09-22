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
	"fmt"
	"os"
	"strings"

	"github.com/ibm-hyper-protect/contract-cli/common"
	"github.com/ibm-hyper-protect/contract-go/v2/crypto"
	"github.com/spf13/cobra"
)

const (
	ParameterName             = "crypto"
	ParameterShortDescription = "Generate OpenSSL RSA key pair or certificate bundle"
	ParameterLongDescription  = `Generate OpenSSL RSA key pair or CA-signed certificate bundle.

Use --type key to generate an RSA private/public key pair.
	 Always outputs a plain RSA public key PEM alongside the private key.
	 --days is not supported for --type key (key pairs do not carry an expiry).

Use --type cert to generate a self-signed CA certificate and a client
certificate signed by that CA, along with the client's private key.
	 --days sets the certificate validity period (default: 365).

The --out flag accepts comma-separated filenames for each output file.
For --type key:  --out <private>,<public>   (2 names max)
For --type cert: --out <ca>,<client-cert>,<client-key>  (3 names max)
If fewer names than files are supplied, missing slots use fixed defaults.
When --out is omitted entirely, all PEM artifacts are printed to stdout.`

	TypeFlagName        = "type"
	TypeFlagDescription = "Generation mode: 'key' for RSA key pair, 'cert' for CA + client certificate bundle (required)"

	OutFlagName        = "out"
	OutFlagDescription = "Comma-separated output filenames. For --type key: <private>,<public>. For --type cert: <ca>,<client-cert>,<client-key>. Missing slots use defaults. Omit to print to stdout."

	SizeFlagName        = "size"
	SizeFlagDescription = "RSA key size in bits: 2048, 3072, or 4096 (default: 4096; omit flag to use default)"

	PasswordFlagName        = "password"
	PasswordFlagDescription = "Passphrase to encrypt the private key with AES-256 (optional, key is unencrypted when not specified)"

	SANFlagName        = "san"
	SANFlagDescription = "Comma-separated Subject Alternative Names, e.g. 'example.com,www.example.com,192.168.1.1' (default: 'example.com', used only with --type cert)"

	DaysFlagName        = "days"
	DaysFlagDescription = "Certificate validity period in days (must be > 0; default: 365). Only valid for --type cert."

	CNFlagName        = "cn"
	CNFlagDescription = "X.509 subject Common Name for the client certificate (default: first SAN entry, used only with --type cert)"

	// Default values applied by the CLI before calling contract-go.
	DefaultKeySize   = 4096
	DefaultSAN       = "example.com"
	DefaultValidDays = 365

	// Default output filenames used when --out provides fewer names than files.
	DefaultKeyPrivateName = "private.pem"
	DefaultKeyPublicName  = "public.pem"
	DefaultCertCAName     = "cert.pem"
	DefaultCertClientName = "cert.pem"
	DefaultCertKeyName    = "privatekey.pem"

	// Maximum number of comma-separated --out values per type.
	maxOutValuesKey  = 2
	maxOutValuesCert = 3
)

// resolveKeyPaths derives the private and public key output paths from the
// raw comma-separated --out value.
//
//   - ""                  → both empty (stdout mode)
//   - "priv.pem"          → privPath="priv.pem", pubPath=DefaultKeyPublicName + log
//   - "priv.pem,pub.pem"  → privPath="priv.pem", pubPath="pub.pem"
//   - 3+ tokens           → first 2 used, warning logged
func resolveKeyPaths(outRaw string) (privPath, pubPath string) {
	if outRaw == "" {
		return "", ""
	}

	tokens := strings.Split(outRaw, ",")

	if len(tokens) > maxOutValuesKey {
		fmt.Printf("Warning: --type key produces %d output files; ignoring extra values\n", maxOutValuesKey)
		tokens = tokens[:maxOutValuesKey]
	}

	privPath = tokens[0]

	if len(tokens) >= 2 {
		pubPath = tokens[1]
	} else {
		pubPath = DefaultKeyPublicName
		fmt.Printf("Public key filename not specified, using default: %s\n", pubPath)
	}

	return privPath, pubPath
}

// resolveCertPaths derives the CA cert, client cert, and client key output
// paths from the raw comma-separated --out value.
//
//   - ""                        → all empty (stdout mode)
//   - "ca.crt"                  → caPath="ca.crt", certPath=DefaultCertClientName + log, keyPath=DefaultCertKeyName + log
//   - "ca.crt,client.crt"       → caPath="ca.crt", certPath="client.crt", keyPath=DefaultCertKeyName + log
//   - "ca.crt,client.crt,k.pem" → caPath="ca.crt", certPath="client.crt", keyPath="k.pem"
//   - 4+ tokens                 → first 3 used, warning logged
func resolveCertPaths(outRaw string) (caPath, certPath, keyPath string) {
	if outRaw == "" {
		return "", "", ""
	}

	tokens := strings.Split(outRaw, ",")

	if len(tokens) > maxOutValuesCert {
		fmt.Printf("Warning: --type cert produces %d output files; ignoring extra values\n", maxOutValuesCert)
		tokens = tokens[:maxOutValuesCert]
	}

	caPath = tokens[0]

	if len(tokens) >= 2 {
		certPath = tokens[1]
	} else {
		certPath = DefaultCertClientName
		fmt.Printf("Client certificate filename not specified, using default: %s\n", certPath)
	}

	if len(tokens) >= 3 {
		keyPath = tokens[2]
	} else {
		keyPath = DefaultCertKeyName
		fmt.Printf("Client private key filename not specified, using default: %s\n", keyPath)
	}

	return caPath, certPath, keyPath
}

// ValidateInput reads and validates all flags from the cobra command.
// It applies defaults for optional flags, splits the --out comma list into
// individual named paths, and returns them ready for the Generate call.
//
// Parameters:
//   - cmd: cobra.Command populated by the CLI flag parser
//
// Returns:
//   - artifactType:  "key" or "cert"
//   - password:      optional AES-256 passphrase (empty string = no encryption)
//   - commonName:    X.509 CN (cert type only; defaults to first SAN token)
//   - sans:          comma-separated SAN list (cert type only; default: "example.com")
//   - privPath:      private key output path  (key type); empty = stdout
//   - pubPath:       public key output path   (key type); empty = stdout
//   - caPath:        CA cert output path      (cert type); empty = stdout
//   - certPath:      client cert output path  (cert type); empty = stdout
//   - certKeyPath:   client key output path   (cert type); empty = stdout
//   - keySize:       RSA key size in bits (default: 4096)
//   - validDays:     certificate validity in days (cert type only; default: 365)
//   - error:         flag parsing error
func ValidateInput(cmd *cobra.Command) (
	artifactType, password, commonName, sans string,
	privPath, pubPath, caPath, certPath, certKeyPath string,
	keySize, validDays int,
	err error,
) {
	artifactType, err = cmd.Flags().GetString(TypeFlagName)
	if err != nil {
		return
	}
	if artifactType == "" {
		common.SetMandatoryFlagError(cmd, fmt.Errorf("Error: required flag '--type' is missing"))
	}
	if artifactType != "key" && artifactType != "cert" {
		common.SetMandatoryFlagError(cmd, fmt.Errorf("Error: invalid value for '--type'. Must be 'key' or 'cert'"))
	}

	outRaw, err := cmd.Flags().GetString(OutFlagName)
	if err != nil {
		return
	}

	switch artifactType {
	case "key":
		privPath, pubPath = resolveKeyPaths(outRaw)
	case "cert":
		caPath, certPath, certKeyPath = resolveCertPaths(outRaw)
	}

	keySize, err = cmd.Flags().GetInt(SizeFlagName)
	if err != nil {
		return
	}
	if !cmd.Flags().Changed(SizeFlagName) {
		keySize = DefaultKeySize
	}

	password, err = cmd.Flags().GetString(PasswordFlagName)
	if err != nil {
		return
	}

	sans, err = cmd.Flags().GetString(SANFlagName)
	if err != nil {
		return
	}
	if sans == "" {
		sans = DefaultSAN
	}

	validDays, err = cmd.Flags().GetInt(DaysFlagName)
	if err != nil {
		return
	}
	if cmd.Flags().Changed(DaysFlagName) {
		// --days is not allowed for --type key
		if artifactType == "key" {
			err = fmt.Errorf("Error: --days is not supported for --type key: key pairs do not carry an expiry; use --type cert to generate a certificate with a validity period")
			return
		}
		// For --type cert, --days must be a positive integer
		if validDays <= 0 {
			err = fmt.Errorf("Error: invalid value for '--days': must be a positive integer, got %d", validDays)
			return
		}
	} else if artifactType == "cert" {
		// Apply the 365-day default when --days is not provided for cert
		validDays = DefaultValidDays
	}

	commonName, err = cmd.Flags().GetString(CNFlagName)
	if err != nil {
		return
	}
	if commonName == "" {
		// Default CN to the first SAN token.
		commonName = strings.SplitN(sans, ",", 2)[0]
	}
	return
}

// Generate calls GenerateOpenSSLArtifacts with the pre-validated parameters
// and writes or prints the resulting PEM artifacts.
//
// For --type key:
//   - privPath / pubPath empty → all output printed to stdout
//   - privPath / pubPath set   → written to those exact paths
//
// For --type cert:
//   - caPath / certPath / certKeyPath empty → all output printed to stdout
//   - set                                   → written to those exact paths
func Generate(
	artifactType, password, commonName, sans string,
	privPath, pubPath, caPath, certPath, certKeyPath string,
	keySize, validDays int,
) error {
	privKey, pubKey, caCert, clientCert, clientKey, _, _, _, _, _, err := crypto.GenerateOpenSSLArtifacts(
		artifactType, password, commonName, sans, keySize, validDays,
	)
	if err != nil {
		return fmt.Errorf("failed to generate OpenSSL artifacts: %w", err)
	}

	switch artifactType {
	case "key":
		if privPath == "" {
			return printToStdout(artifactType, privKey, pubKey, "", "", "")
		}
		return writeKeyFiles(privPath, pubPath, privKey, pubKey)

	case "cert":
		if caPath == "" {
			return printToStdout(artifactType, "", "", caCert, clientCert, clientKey)
		}
		return writeCertFiles(caPath, certPath, certKeyPath, caCert, clientCert, clientKey)
	}

	return nil
}

// printToStdout writes labelled PEM blocks to stdout.
func printToStdout(artifactType, privKey, pubKey, caCert, clientCert, clientKey string) error {
	switch artifactType {
	case "key":
		fmt.Printf("Private Key:\n%s\nPublic Key:\n%s\n", privKey, pubKey)
	case "cert":
		fmt.Printf("CA Certificate:\n%s\nClient Certificate:\n%s\nClient Private Key:\n%s\n",
			caCert, clientCert, clientKey)
	}
	return nil
}

// writeKeyFiles writes the private and public key PEM files.
func writeKeyFiles(privPath, pubPath, privKey, pubKey string) error {
	if err := writeFile(privPath, privKey, 0600); err != nil {
		return err
	}
	if err := writeFile(pubPath, pubKey, 0644); err != nil {
		return err
	}
	fmt.Printf("Private key written to: %s\nPublic key written to:  %s\n", privPath, pubPath)
	return nil
}

// writeCertFiles writes the CA cert, client cert, and client private key PEM files.
func writeCertFiles(caPath, certPath, keyPath, caCert, clientCert, clientKey string) error {
	if err := writeFile(caPath, caCert, 0644); err != nil {
		return err
	}
	if err := writeFile(certPath, clientCert, 0644); err != nil {
		return err
	}
	if err := writeFile(keyPath, clientKey, 0600); err != nil {
		return err
	}
	fmt.Printf("CA certificate written to:     %s\nClient certificate written to: %s\nClient private key written to: %s\n",
		caPath, certPath, keyPath)
	return nil
}

// writeFile writes data to path with the given permission bits.
func writeFile(path, data string, perm os.FileMode) error {
	if err := os.WriteFile(path, []byte(data), perm); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}
