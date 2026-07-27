<!--
Copyright (c) 2026 IBM Corp.
All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Contract CLI Documentation

Complete command reference and usage guide for the IBM Confidential Computing Contract CLI.

## Table of Contents

- [Introduction](#introduction)
- [Installation](#installation)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
  - [base64](#base64)
  - [base64-tgz](#base64-tgz)
  - [contract-template](#contract-template)
  - [decrypt-attestation](#decrypt-attestation)
  - [download-certificate](#download-certificate)
  - [sign-contract](#sign-contract)
  - [encrypt](#encrypt)
  - [encrypt-string](#encrypt-string)
  - [get-certificate](#get-certificate)
  - [image](#image)
  - [list-encryptioncert-versions](#list-encryptioncert-versions)
  - [sealed-secret](#sealed-secret)
  - [validate-contract](#validate-contract)
  - [validate-network](#validate-network)
  - [validate-encryption-certificate](#validate-encryption-certificate)
  - [initdata](#initdata)
- [Common Workflows](#common-workflows)
- [CI/CD Integration](#cicd-integration)
- [Exit Codes](#exit-codes)
- [Troubleshooting](#troubleshooting)
- [Examples](#examples)

## Introduction

The Contract CLI automates the process of generating and managing contracts for provisioning IBM Confidential Computing services including IBM Confidential Computing Container Runtime, IBM Confidential Computing Container Runtime for Red Hat Virtualization Solutions, and IBM Confidential Computing Containers for Red Hat OpenShift Container Platform. It provides a comprehensive set of commands for:

- Generating signed and encrypted contracts
- Managing encryption certificates
- Validating contracts and network configurations
- Handling attestation records
- Working with container configurations

## Installation

Download the latest release for your platform from the [releases page](https://github.com/ibm-hyper-protect/contract-cli/releases/latest).

### Available Platforms

- **Linux**: amd64, arm64, s390x, ppc64le
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64, arm64

### Verify Installation

```bash
# Check version
contract-cli --version

# View available commands
contract-cli --help
```

## Prerequisites

### OpenSSL

OpenSSL is required for all cryptographic operations. The CLI will use the `openssl` binary from your system PATH.

**Installation:**
- **Linux**: `apt-get install openssl` or `yum install openssl`
- **macOS**: `brew install openssl`
- **Windows**: [Download OpenSSL](https://slproweb.com/products/Win32OpenSSL.html)

### Custom OpenSSL Path (Optional)

If OpenSSL is not in your system PATH, configure the `OPENSSL_BIN` environment variable:

**Linux/macOS:**
```bash
export OPENSSL_BIN=/usr/bin/openssl
```

**Windows (PowerShell):**
```powershell
$env:OPENSSL_BIN="C:\Program Files\OpenSSL-Win64\bin\openssl.exe"
```

## Quick Start

### Generate a Complete Contract

```bash
# 1. Generate RSA key pair
openssl genrsa -out private.pem 4096

# 2. Create your contract YAML
cat > contract.yaml <<EOF
env: |
  type: env
  logging:
    logRouter:
      hostname: example.logs.cloud.ibm.com
      iamApiKey: your-api-key
workload: |
  type: workload
  compose:
    archive: your-archive
EOF 

# 3. Validate the contract
contract-cli validate-contract --in contract.yaml --os hpvs

# 4. Generate signed and encrypted contract
contract-cli encrypt --in contract.yaml --priv private.pem --out encrypted-contract.yaml
```

---

## Command Reference

### base64

Encode text or JSON data to Base64 format. Useful for encoding data that needs to be included in contracts or configurations.

#### Usage

```bash
contract-cli base64 [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Input data to encode (text or JSON) (use '-' for standard input) |
| `--format` | string | No | Input data format (text or json) |
| `--out` | string | No | Path to save Base64 encoded output |
| `-h, --help` | - | No | Display help information |

#### Examples

**Basic text encoding:**
```bash
contract-cli base64 --in "Hello World" --format text
```

**JSON encoding:**
```bash
contract-cli base64 --in '{"type": "workload"}' --format json
```

**Save to file:**
```bash
contract-cli base64 --in "Hello World" --format text --out encoded.txt
```

**Using standard input (pipe input):**
```bash
echo "Hello World" | contract-cli base64 --in - --format text
```

---

### base64-tgz

Generate Base64-encoded tar.gz archive of docker-compose.yaml or pods.yaml. Creates a compressed archive of your container configuration files, encoded as Base64 for inclusion in Confidential Computing contracts. Supports both plain and encrypted output.

#### Usage

```bash
contract-cli base64-tgz [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to folder containing `docker-compose.yaml` or `pods.yaml` (use '-' for standard input) |
| `--output` | string | No | Output type: `plain` or `encrypted` (default: `plain`) |
| `--cert` | string | No | Path to encryption certificate (uses latest embedded certificate for the provided IBM Confidential Computing platform, if not specified) |
| `--os` | string | No | Target IBM Confidential Computing platform: `ccrt`, `ccrv`, `ccco`, or `hpvs` (default: `hpvs`) |
| `--ver` | string | No | Specific encryption certificate version (e.g., `26.2.0`). Uses latest version if not specified. Use `list-encryptioncert-versions` to see available versions |
| `--out` | string | No | Path to save the output |
| `-h, --help` | - | No | Display help information |

#### Examples

**Plain Base64 archive:**
```bash
contract-cli base64-tgz --in ./compose-folder
```

**Encrypted archive with latest certificate:**
```bash
contract-cli base64-tgz --in ./compose-folder --output encrypted
```

**Encrypted archive with custom certificate:**
```bash
contract-cli base64-tgz \
  --in ./compose-folder \
  --output encrypted \
  --cert encryption.crt
```

**For HPCR-RHVS:**
```bash
contract-cli base64-tgz \
  --in ./pods-folder \
  --output encrypted \
  --os ccrv
```

**For CCCO:**
```bash
contract-cli base64-tgz \
  --in ./pods-folder \
  --output encrypted \
  --os ccco
```

**Save to file:**
```bash
contract-cli base64-tgz --in ./compose-folder --out archive.txt
```

**With specific certificate version:**
```bash
contract-cli base64-tgz \
  --in ./compose-folder \
  --output encrypted \
  --os ccco \
  --ver 25.12.0
```

**Using standard input (pipe input):**
```bash
echo "pods-folder" | contract-cli base64-tgz --in -
```

---

### contract-template

Generate a contract YAML template for IBM Confidential Computing deployments. Returns a pre-filled YAML scaffold for the workload section, env section, or a combined contract containing both. Use this as a starting point when authoring a new contract.

#### Usage

```bash
contract-cli contract-template [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--type` | string | No | Template type to generate: `env`, `workload`, or `contract` (default: `contract`) |
| `--os` | string | No | Target platform (default: `hpvs`). See OS values table below. |
| `--out` | string | No | Path to save the generated template (prints to terminal if not specified) |
| `-h, --help` | - | No | Display help information |

#### OS Values

| Value | Platform | Workload template | Env template |
|-------|----------|-------------------|--------------|
| `hpvs` | IBM Hyper Protect Virtual Servers | compose + play + volumes | standard (syslog, env vars, volumes) |
| `ccrt` | IBM Confidential Computing Container Runtime | compose + play + volumes | standard |
| `ccrv` | IBM CCRT for Red Hat Virtualization | play only (no compose) | standard |
| `ccco-peerpod` | IBM CCCO Peer Pod | confidential-containers (no volumes) | logRouter only |
| `ccco-bmtl` | IBM CCCO Baremetal | confidential-containers + volumes | logRouter + volumes + host-attestation |

#### Examples

**Generate combined contract template (default):**
```bash
contract-cli contract-template
```

**Generate workload-only template:**
```bash
contract-cli contract-template --type workload
```

**Generate env-only template:**
```bash
contract-cli contract-template --type env
```

**Generate template for CCRT:**
```bash
contract-cli contract-template --type contract --os ccrt
```

**Save combined template to file:**
```bash
contract-cli contract-template --out contract-template.yaml
```

**Generate CCRV workload template:**
```bash
contract-cli contract-template \
  --type workload \
  --os ccrv \
  --out ccrv-workload-template.yaml
```

**Generate CCCO Peer Pod workload template:**
```bash
contract-cli contract-template \
  --type workload \
  --os ccco-peerpod \
  --out ccco-peerpod-workload.yaml
```

**Generate CCCO Baremetal combined template:**
```bash
contract-cli contract-template \
  --type contract \
  --os ccco-bmtl \
  --out ccco-bmtl-contract.yaml
```

---

### decrypt-attestation

Decrypt encrypted attestation records generated by Confidential Computing instances. Attestation records are typically found at `/var/hyperprotect/se-checksums.txt.enc` and contain cryptographic hashes for verifying workload integrity.

Optionally verify the signature of decrypted attestation records by providing both `--signature` and `--attestation-cert` flags together.

#### Usage

```bash
contract-cli decrypt-attestation [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to encrypted attestation file (use '-' for standard input) |
| `--priv` | string | Yes | Path to private key used for decryption |
| `--password` | string | No | Password for encrypted private key |
| `--out` | string | No | Path to save decrypted attestation records |
| `--signature` | string | No* | Path to signature file (se-signature.bin) |
| `--attestation-cert` | string | No* | Path to IBM attestation certificate file (PEM format) |
| `-h, --help` | - | No | Display help information |

\* Both `--signature` and `--attestation-cert` must be provided together if signature verification is desired

#### Examples

**Decrypt to console:**
```bash
contract-cli decrypt-attestation \
  --in se-checksums.txt.enc \
  --priv private.pem
```

**Decrypt and save to file:**
```bash
contract-cli decrypt-attestation \
  --in se-checksums.txt.enc \
  --priv private.pem \
  --out decrypted-attestation.txt
```

**Decrypt and verify signature:**
```bash
contract-cli decrypt-attestation \
  --in se-checksums.txt.enc \
  --priv private.pem \
  --signature se-signature.bin \
  --attestation-cert hpse-attestation.crt
```

**Decrypt, verify signature, and save to file:**
```bash
contract-cli decrypt-attestation \
  --in se-checksums.txt.enc \
  --priv private.pem \
  --out decrypted-attestation.txt \
  --signature se-signature.bin \
  --attestation-cert hpse-attestation.crt
```

**Using standard input:**
```bash
cat se-checksums.txt.enc | contract-cli decrypt-attestation \
  --in - \
  --priv private.pem
```

**Using password-protected private key:**
```bash
contract-cli decrypt-attestation \
  --in se-checksums.txt.enc \
  --priv private-encrypted.pem \
  --password "your-secure-password" \
  --out decrypted-attestation.txt
```

**Using password from environment variable:**
```bash
export PRIVATE_KEY_PASSWORD="your-secure-password"
contract-cli decrypt-attestation \
  --in se-checksums.txt.enc \
  --priv private-encrypted.pem \
  --password "$PRIVATE_KEY_PASSWORD" \
  --out decrypted-attestation.txt
```

---

### download-certificate

Download encryption certificates from the IBM Confidential Computing Repository. Retrieves the latest or specific versions of encryption certificates required for contract encryption and workload deployment.

#### Usage

```bash
contract-cli download-certificate [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--version` | strings | Yes | Specific certificate versions to download (comma-separated, e.g., 1.0.21,1.0.22) |
| `--format` | string | No | Output format for data (json, yaml, or text) |
| `--out` | string | No | Path to save downloaded encryption certificates |
| `-h, --help` | - | No | Display help information |

#### Examples

**Download latest certificate:**
```bash
contract-cli download-certificate
```

**Download specific version:**
```bash
contract-cli download-certificate --version 1.0.23
```

**Download multiple versions:**
```bash
contract-cli download-certificate --version 1.0.21,1.0.22,1.0.23
```

**Save to file in YAML format:**
```bash
contract-cli download-certificate \
  --version 1.0.23 \
  --format yaml \
  --out certificates.yaml
```

---

### sign-contract

Generates a signed contract from a contract with encrypted workload and env sections.

#### Usage

```bash
contract-cli sign-contract [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to encrypted IBM Confidential Computing contract YAML file (use '-' for standard input) |
| `--priv` | string | Yes | Path to private key for signing |
| `--password` | string | No | Password for encrypted private key |
| `--out` | string | No | Path to save signed and encrypted contract |
| `-h, --help` | - | No | Display help information |

#### Examples

**Sign a contract:**
```bash
contract-cli sign-contract --in contract.yaml --priv private.pem
```

**Sign and save to file:**
```bash
contract-cli sign-contract --in contract.yaml --priv private.pem --out signed-contract.yaml
```

**Using standard input:**
```bash
cat contract.yaml | contract-cli sign-contract --in - --priv private.pem
```

**Using password-protected private key:**
```bash
contract-cli sign-contract \
  --in contract.yaml \
  --priv private-encrypted.pem \
  --password "your-secure-password" \
  --out signed-contract.yaml
```

---

### encrypt

Generate a signed and encrypted contract for IBM Confidential Computing deployment. Supports optional contract expiry for enhanced security.

#### Usage

```bash
contract-cli encrypt [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to unencrypted IBM Confidential Computing contract YAML file (use '-' for standard input) |
| `--priv` | string | No* | Path to private key for signing |
| `--password` | string | No | Password for encrypted private key |
| `--cert` | string | No | Path to encryption certificate (uses latest embedded certificate for the provided IBM Confidential Computing platform, if not specified) |
| `--os` | string | No | Target IBM Confidential Computing platform: `ccrt`, `ccrv`, `ccco`, or `hpvs` (default: `hpvs`) |
| `--ver` | string | No | Specific encryption certificate version (e.g., `26.2.0`). Uses latest version if not specified. Use `list-encryptioncert-versions` to see available versions |
| `--out` | string | No | Path to save signed and encrypted contract |
| `--contract-expiry` | bool | No | Enable contract expiry feature |
| `--cacert` | string | No** | Path to CA certificate (required with expiry) |
| `--cakey` | string | No** | Path to CA key (required with expiry) |
| `--csr` | string | No** | Path to CSR file (required with expiry) |
| `--csrParam` | string | No** | Path to CSR parameters JSON |
| `--expiry` | int | No** | Contract validity in days (required with expiry) |
| `-h, --help` | - | No | Display help information |

\* Generated automatically if not provided
\** Required when `--contract-expiry` is enabled

#### Examples

**Basic encryption:**
```bash
contract-cli encrypt --in contract.yaml --priv private.pem
```

**With custom certificate:**
```bash
contract-cli encrypt \
  --in contract.yaml \
  --priv private.pem \
  --cert encryption.crt
```

**Save to file:**
```bash
contract-cli encrypt \
  --in contract.yaml \
  --priv private.pem \
  --out encrypted-contract.yaml
```

**With specific certificate version:**
```bash
contract-cli encrypt \
  --in contract.yaml \
  --priv private.pem \
  --os ccrt \
  --ver 26.2.0
```

**With contract expiry:**
```bash
contract-cli encrypt \
  --contract-expiry \
  --in contract.yaml \
  --priv private.pem \
  --cacert ca.crt \
  --cakey ca.key \
  --csr csr.pem \
  --expiry 90
```

**For CCRV:**
```bash
contract-cli encrypt \
  --in contract.yaml \
  --priv private.pem \
  --os ccrv
```

**Using password-protected private key:**
```bash
contract-cli encrypt \
  --in contract.yaml \
  --priv private-encrypted.pem \
  --password "your-secure-password" \
  --out encrypted-contract.yaml
```

**For CCCO:**
```bash
contract-cli encrypt \
  --in contract.yaml \
  --priv private.pem \
  --os ccco
```

**Using standard input:**
```bash
echo "test-string" | contract-cli encrypt \
  --in - \
  --priv private.pem
```

---

### encrypt-string

Encrypt strings using the IBM Confidential Computing encryption format. Output format: `hyper-protect-basic.<encrypted-password>.<encrypted-string>`. Use this to encrypt sensitive data like passwords or API keys for contracts.

#### Usage

```bash
contract-cli encrypt-string [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | String data to encrypt (use '-' for standard input) |
| `--format` | string | No | Input data format (text or json) |
| `--cert` | string | No | Path to encryption certificate (uses latest embedded certificate for the provided IBM Confidential Computing platform, if not specified) |
| `--os` | string | No | Target IBM Confidential Computing platform: `ccrt`, `ccrv`, `ccco`, or `hpvs` (default: `hpvs`) |
| `--ver` | string | No | Specific encryption certificate version (e.g., `26.2.0`). Uses latest version if not specified. Use `list-encryptioncert-versions` to see available versions |
| `--out` | string | No | Path to save encrypted output |
| `-h, --help` | - | No | Display help information |

#### Examples

**Encrypt plain text:**
```bash
contract-cli encrypt-string --in "my-secret-password"
```

**Encrypt JSON:**
```bash
contract-cli encrypt-string \
  --in '{"apiKey": "secret123"}' \
  --format json
```

**With custom certificate:**
```bash
contract-cli encrypt-string \
  --in "my-secret" \
  --cert encryption.crt
```

**Save to file:**
```bash
contract-cli encrypt-string \
  --in "my-secret" \
  --out encrypted-secret.txt
```

**With specific certificate version:**
```bash
contract-cli encrypt-string \
  --in "my-secret-password" \
  --os ccrv \
  --ver 25.11.0
```

**Using standard input:**
```bash
# Encrypt echo statement
echo "my-secret-password" | contract-cli encrypt-string --in -

# Encrypt file content
cat workload.yaml | contract-cli encrypt-string --in -
```

---

### get-certificate

Extract a specific encryption certificate version from download-certificate output. Parses the JSON output from download-certificate and extracts the certificate for the specified version.

#### Usage

```bash
contract-cli get-certificate [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to download-certificate JSON output (use '-' for standard input) |
| `--version` | string | Yes | Certificate version to extract (e.g., 1.0.23) |
| `--out` | string | No | Path to save extracted encryption certificate |
| `-h, --help` | - | No | Display help information |

#### Examples

**Extract specific version:**
```bash
contract-cli get-certificate \
  --in certificates.json \
  --version 1.0.23
```

**Save to file:**
```bash
contract-cli get-certificate \
  --in certificates.json \
  --version 1.0.23 \
  --out cert-1.0.23.crt
```

**Using standard input:**
```bash
cat "cert.json" | contract-cli get-certificate --in - --version 1.0.23
```

---

### image

Retrieve IBM Confidential Computing Container Runtime image details from IBM Cloud. Parses image information from IBM Cloud API, CLI, or Terraform output to extract image ID, name, checksum, and version. Supports filtering by specific version.

#### Usage

```bash
contract-cli image [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to IBM Cloud images JSON (from API, CLI, or Terraform) (use '-' for standard input) |
| `--version` | string | No | Specific version to retrieve (returns latest if not specified) |
| `--format` | string | No | Output format for data (json, yaml, or text) |
| `--out` | string | No | Path to save image details |
| `-h, --help` | - | No | Display help information |

#### Examples

**Get latest image:**
```bash
contract-cli image --in ibm-cloud-images.json
```

**Get specific version:**
```bash
contract-cli image \
  --in ibm-cloud-images.json \
  --version "1.0.23"
```

**Output in YAML:**
```bash
contract-cli image \
  --in ibm-cloud-images.json \
  --format yaml
```

**Save to file:**
```bash
contract-cli image \
  --in ibm-cloud-images.json \
  --out hpcr-image.json
```

**Using standard input:**
```bash
cat "ibm-cloud-images.json" | contract-cli image --in -
```


### list-encryptioncert-versions

List all available embedded encryption certificate versions for IBM Confidential Computing platforms. This command helps you discover which certificate versions are available before using the `--ver` flag with encryption commands. The embedded certificates are bundled with the CLI and don't require downloading from IBM Cloud.

#### Usage

```bash
contract-cli list-encryptioncert-versions [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--os` | string | No | Filter by platform (ccrt, ccrv, ccco, or hpvs). Shows all platforms if not specified |
| `--format` | string | No | Output format: `json` or `yaml` (defaults to `json` if not specified) |
| `--out` | string | No | Path to save output (prints to stdout if not specified) |
| `-h, --help` | - | No | Display help information |

#### Examples

**List all available encryption certificate versions in JSON format (default):**
```bash
contract-cli list-encryptioncert-versions
```

Output:
```json
{"ccco":["25.12.0","25.10.0"],"ccrt":["26.5.0","26.2.0"],"ccrv":["26.4.1","25.11.0"],"hpvs":["26.5.0","26.2.0"]}
```

**List all available encryption certificate versions in YAML format:**
```bash
contract-cli list-encryptioncert-versions --format yaml
```

Output:
```yaml
ccco:
  - 25.12.0
  - 25.10.0
ccrt:
  - 26.5.0
  - 26.2.0
ccrv:
  - 26.4.1
  - 25.11.0
hpvs:
  - 26.5.0
  - 26.2.0
```

**List versions for a specific platform in JSON:**
```bash
contract-cli list-encryptioncert-versions --os ccrt --format json
```

Output:
```json
{"ccrt":["26.5.0","26.2.0"]}
```

**List versions for HPVS platform in YAML:**
```bash
contract-cli list-encryptioncert-versions --os hpvs --format yaml
```

Output:
```yaml
hpvs:
  - 26.5.0
  - 26.2.0
```

**Save output to file:**
```bash
contract-cli list-encryptioncert-versions --os ccrv --format yaml --out ccrv-versions.yaml
```

#### Use Cases

1. **Discover Available Versions**: Find out which encryption certificate versions are embedded in your CLI installation
2. **Version Selection**: Choose a specific version for encryption operations using the `--ver` flag
3. **Compatibility Check**: Verify that a required certificate version is available before running automation scripts
4. **Documentation**: Generate a list of supported versions for your deployment documentation

#### Related Commands

- [`encrypt`](#encrypt) - Use `--ver` flag to specify certificate version for contract encryption
- [`encrypt-string`](#encrypt-string) - Use `--ver` flag to specify certificate version for string encryption
- [`base64-tgz`](#base64-tgz) - Use `--ver` flag with `--output encrypt` to specify encryption certificate version

---

---

### validate-contract

Validate an unencrypted contract against the IBM Confidential Computing schema. Checks contract structure, required fields, and data types before encryption to help catch errors early in the development process.

#### Usage

```bash
contract-cli validate-contract [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to unencrypted IBM Confidential Computing contract YAML file (use '-' for standard input) |
| `--os` | string | No | Target IBM Confidential Computing platform: `ccrt`, `ccrv`, `ccco`, or `hpvs` (default: `hpvs`) |
| `--type` | string | No | Contract section to validate: `workload`, `env`, or `''` for both (default: `''`) |
| `-h, --help` | - | No | Display help information |

#### Examples

**Validate full contract (both sections):**
```bash
contract-cli validate-contract --in contract.yaml --os ccrt
```

**Validate only the workload section:**
```bash
contract-cli validate-contract --in contract.yaml --os ccrt --type workload
```

**Validate only the env section:**
```bash
contract-cli validate-contract --in contract.yaml --os ccrt --type env
```

**Validate CCRV contract:**
```bash
contract-cli validate-contract --in contract.yaml --os ccrv
```

**Validate CCCO contract:**
```bash
contract-cli validate-contract --in contract.yaml --os ccco
```

**Using standard input:**
```bash
cat contract.yaml | contract-cli validate-contract --in - --os ccrt
```

---

### validate-network

Validate network-config YAML file against the schema. Validates network configuration for on-premise deployments, ensuring all required fields are present and properly formatted.

#### Usage

```bash
contract-cli validate-network [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to network-config YAML file (use '-' for standard input) |
| `-h, --help` | - | No | Display help information |

#### Examples

**Validate network configuration:**
```bash
contract-cli validate-network --in network-config.yaml
```

**Using standard input:**
```bash
cat network-config.yaml | contract-cli validate-network --in -
```

---


### validate-encryption-certificate

Validates encryption certificate for on-premise, VPC deployment. It will check encryption certificate validity, ensuring all required fields are present and properly formatted.

#### Usage

```bash
contract-cli validate-encryption-certificate [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to encryption certificate file (use '-' for standard input) |
| `-h, --help` | - | No | Display help information |

#### Examples

**Validate encryption certificate configuration:**
```bash
contract-cli validate-encryption-certificate --in encryption-cert.crt
```

**Using standard input:**
```bash
cat encryption-cert.crt | contract-cli validate-encryption-certificate --in -
```

---

### sealed-secret
Generate sealed secrets for IBM Confidential Computing Containers for Red Hat OpenShift Container Platform (CCCO).

#### Usage

```bash
contract-cli sealed-secret [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Secret for sealing (provide as string or file path, use '-' for standard input) |
| `--type` | string | Yes | Type of secret: 'env' for env section of contract or 'workload' for workload section of contract |
| `--out` | string | No | Path to save sealed secret output (prints to stdout if not specified) |
| `--encryptionkey` | string | No | Path to RSA private key for encryption (generates new key if not provided) |
| `--signingkey` | string | No | Path to RSA private key for signing (generates new key if not provided) |
| `-h, --help` | - | No | Display help information |

#### Examples

**Generate sealed secret for env section:**
```bash
contract-cli sealed-secret \
  --in "value123" \
  --type env \
  --out sealed-secret.txt
```

**Generate sealed secret for workload section:**
```bash
contract-cli sealed-secret \
  --in workload-secret-data \
  --type workload \
  --out sealed-workload.txt
```

**Generate sealed secret from file:**
```bash
contract-cli sealed-secret \
  --in secrets.txt \
  --type env \
  --out sealed-secret.txt
```

**Generate sealed secret with custom encryption and signing keys:**
```bash
contract-cli sealed-secret \
  --in "value123" \
  --type env \
  --encryptionkey encryption.key \
  --signingkey signing.key \
  --out sealed-secret.txt
```

**Read secret from stdin:**
```bash
echo "value123" | contract-cli sealed-secret \
  --in - \
  --type env
```

**Output format:**
The command outputs:
- The sealed secret data (for use in contract)
- `SECRET_DECRYPTION_KEY` - Private key for decryption (keep secure)
- `SECRET_VERIFICATION_KEY` - Public key for verification

---

### initdata
Create initdata annotation from signed and encrypted contract for IBM Confidential Computing Containers for Red Hat OpenShift Container Platform. Supports both Peer Pod and Baremetal solutions.

#### Usage

```bash
contract-cli initdata [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to signed & encrypted contract YAML file (use '-' for standard input) |
| `--sehdr` | string | No | Path to SE header binary file (.bin) for baremetal solution |
| `--out` | string | No | Path to store gzipped & encoded initdata value |
| `-h, --help` | - | No | Display help information |

#### Examples

**Create initdata for Peer Pod solution without SE header binary:**
```bash
contract-cli initdata --in signed_encrypted_contract.yaml
```

**Create initdata for Baremetal solution with SE header binary:**
```bash
contract-cli initdata \
  --in signed_encrypted_contract.yaml \
  --sehdr se-header.bin \
  --out initdata.txt
```

**Save output to file for peerpod solution without SE header binary:**
```bash
contract-cli initdata \
  --in signed_encrypted_contract.yaml \
  --out initdata-annotation.txt
```

**Save output to file for baremetal solution with SE header binary:**
```bash
contract-cli initdata \
  --in signed_encrypted_contract.yaml \
  --sehdr se-header.bin \
  --out initdata-annotation.txt
```


**Using standard input:**
```bash
cat signed_encrypted_contract.yaml | contract-cli initdata --in -
```

#### Notes

- With `--sehdr`, the command generates initdata for baremetal solution
- Without `--sehdr`, the command generates initdata for Peer Pod solution
- The SE header binary file is automatically encoded to base64 before being included in the initdata
- Output is gzipped and base64 encoded, ready to use as an initdata annotation

---

## Common Workflows

### Complete Contract Generation Workflow

```bash
# Step 1: Generate key pair
openssl genrsa -out private.pem 4096

# Step 2: Download encryption certificate
contract-cli download-certificate --version 1.0.23 --out certs.json
contract-cli get-certificate --in certs.json --version 1.0.23 --out cert.crt

# Step 3: Create docker-compose archive
contract-cli base64-tgz --in ./compose-folder --output encrypted --cert cert.crt --out archive.txt

# Step 4: Create contract YAML (with archive from step 3)
cat > contract.yaml <<EOF
env: |
  type: env
  logging:
    logRouter:
      hostname: logs.example.com
workload: |
  type: workload
  compose:
    archive: $(cat archive.txt)
EOF

# Step 5: Validate contract
contract-cli validate-contract --in contract.yaml --os hpvs

# Step 6: Generate signed and encrypted contract
contract-cli encrypt --in contract.yaml --priv private.pem --cert cert.crt --out final-contract.yaml
```

### Working with Attestation Records

```bash
# Decrypt attestation from running instance
contract-cli decrypt-attestation \
  --in se-checksums.txt.enc \
  --priv private.pem \
  --out attestation.txt

# Decrypt and verify signature
contract-cli decrypt-attestation \
  --in se-checksums.txt.enc \
  --priv private.pem \
  --signature se-signature.bin \
  --attestation-cert hpse-attestation.crt \
  --out attestation.txt

# View decrypted attestation
cat attestation.txt
```

### Certificate Management

```bash
# Download all available certificates
contract-cli download-certificate --out all-certs.json

# Extract specific version
contract-cli get-certificate \
  --in all-certs.json \
  --version 1.0.23 \
  --out cert-1.0.23.crt
```

---

## CI/CD Integration

The CLI supports `stdin` input (`--in -`) for all commands, making it easy to integrate into CI/CD pipelines.

### GitHub Actions Example

```yaml
name: Generate Contract
on:
  push:
    branches: [main]

jobs:
  generate-contract:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Download contract-cli
        run: |
          curl -sL https://github.com/ibm-hyper-protect/contract-cli/releases/latest/download/contract-cli-linux-amd64 -o contract-cli
          chmod +x contract-cli

      - name: Validate contract
        run: ./contract-cli validate-contract --in contract.yaml --os hpvs

      - name: Generate signed and encrypted contract
        run: |
          ./contract-cli encrypt \
            --in contract.yaml \
            --priv "${{ secrets.PRIVATE_KEY_PATH }}" \
            --out encrypted-contract.yaml

      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: encrypted-contract
          path: encrypted-contract.yaml
```

### Shell Script Pattern

```bash
#!/bin/bash
set -euo pipefail

# Validate → Encrypt → Deploy pattern
if ! contract-cli validate-contract --in contract.yaml --os hpvs; then
  echo "Contract validation failed" >&2
  exit 1
fi

contract-cli encrypt \
  --in contract.yaml \
  --priv private.pem \
  --out encrypted-contract.yaml

echo "Contract generated successfully"
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | General error (invalid input, missing files, encryption failure, etc.) |

All error messages are written to `stderr`. Successful output is written to `stdout` (unless `--out` is specified).

---

## Troubleshooting

### OpenSSL Not Found

**Error:**
```
Error: openssl binary not found in PATH
```

**Solution:**
- Install OpenSSL for your platform
- Or set `OPENSSL_BIN` environment variable to the full path of OpenSSL

### Invalid Contract Schema

**Error:**
```
Error: contract validation failed
```

**Solution:**
- Run `validate-contract` to see specific schema errors
- Check contract structure matches IBM Confidential Computing requirements
- Ensure all required fields are present

### Certificate Version Not Found

**Error:**
```
Error: certificate version not found
```

**Solution:**
- Run `download-certificate` without `--version` to see available versions
- Verify the version number format (e.g., `1.0.23`)

### Permission Denied

**Error:**
```
Error: permission denied reading file
```

**Solution:**
- Check file permissions: `chmod 600 private.pem`
- Ensure you have read access to input files
- Verify output directory is writable

---

## Examples

The [`samples/`](../samples/) directory contains working examples:

- **[Contract](../samples/contract.yaml)** - Basic contract structure
- **[Contract with Expiry](../samples/contract-expiry/)** - Contract with expiration
- **[Attestation Records](../samples/attestation/)** - Example attestation files
- **[Certificate Examples](../samples/certificate/)** - Encryption certificate samples
- **[Network Configuration](../samples/network/)** - Network config examples
- **[Docker Compose](../samples/tgz/)** - Compose file examples
- **[Signed & Encrypted Contract](../samples/hpcc/signed-encrypt-hpcc.yaml)** - Signed & Encrypted contract
- **[Contract Signing](../samples/sign/)** - Contract signing examples

---

## Additional Resources

- **[Main README](../README.md)** - Project overview and quick start
- **[Contributing Guide](../CONTRIBUTING.md)** - How to contribute
- **[Security Policy](../SECURITY.md)** - Security best practices
- **[Changelog](../CHANGELOG.md)** - Release history and version notes


### IBM Confidential Computing Documentation

- [Confidential computing with LinuxONE](https://cloud.ibm.com/docs/vpc?topic=vpc-about-se)
- [IBM Confidential Computing Container Runtime](https://www.ibm.com/docs/en/cccr/2.2.x)
- [IBM Confidential Computing Container Runtime for Red Hat Virtualization Solutions](https://www.ibm.com/docs/en/ccrv/1.1.x)
- [IBM Confidential Computing Containers for Red Hat OpenShift](https://www.ibm.com/docs/en/ccro/1.1.x)

### Related Projects

- [contract-go](https://github.com/ibm-hyper-protect/contract-go) - Go library
- [terraform-provider-hpcr](https://github.com/ibm-hyper-protect/terraform-provider-hpcr) - Terraform provider
- [k8s-operator-hpcr](https://github.com/ibm-hyper-protect/k8s-operator-hpcr) - Kubernetes operator

---

**Need Help?**

- [Open an issue](https://github.com/ibm-hyper-protect/contract-cli/issues/new/choose)
- [Ask a question](https://github.com/ibm-hyper-protect/contract-cli/discussions)
- Check the [troubleshooting](#troubleshooting) section above

