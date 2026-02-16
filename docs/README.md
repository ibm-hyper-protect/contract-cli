# Contract CLI Documentation

Complete command reference and usage guide for the Hyper Protect Contract CLI.

## Table of Contents

- [Introduction](#introduction)
- [Installation](#installation)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
  - [base64](#base64)
  - [base64-tgz](#base64-tgz)
  - [decrypt-attestation](#decrypt-attestation)
  - [download-certificate](#download-certificate)
  - [encrypt](#encrypt)
  - [encrypt-string](#encrypt-string)
  - [get-certificate](#get-certificate)
  - [image](#image)
  - [validate-contract](#validate-contract)
  - [validate-network](#validate-network)
  - [initdata](#initdata)
- [Common Workflows](#common-workflows)
- [Troubleshooting](#troubleshooting)
- [Examples](#examples)

## Introduction

The Contract CLI automates the process of generating and managing contracts for provisioning IBM Hyper Protect services including Hyper Protect Virtual Servers (HPVS) for VPC, Hyper Protect Container Runtime (HPCR) for RHVS, and Hyper Protect Confidential Container (HPCC) for Peer Pods. It provides a comprehensive set of commands for:

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
| `--in` | string | Yes | Input data to encode (text or JSON) |
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

---

### base64-tgz

Generate Base64-encoded tar.gz archive of docker-compose.yaml or pods.yaml. Creates a compressed archive of your container configuration files, encoded as Base64 for inclusion in Hyper Protect contracts. Supports both plain and encrypted output.

#### Usage

```bash
contract-cli base64-tgz [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to folder containing `docker-compose.yaml` or `pods.yaml` |
| `--output` | string | No | Output type: `plain` or `encrypted` (default: `plain`) |
| `--cert` | string | No | Path to encryption certificate (for encrypted output) |
| `--os` | string | No | Target Hyper Protect platform: `hpvs`, `hpcr-rhvs`, or `hpcc-peerpod` (default: `hpvs`) |
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
  --os hpcr-rhvs
```

**For HPCC Peer Pods:**
```bash
contract-cli base64-tgz \
  --in ./pods-folder \
  --output encrypted \
  --os hpcc-peerpod
```

**Save to file:**
```bash
contract-cli base64-tgz --in ./compose-folder --out archive.txt
```

---

### decrypt-attestation

Decrypt encrypted attestation records generated by Hyper Protect instances. Attestation records are typically found at `/var/hyperprotect/se-checksums.txt.enc` and contain cryptographic hashes for verifying workload integrity.

#### Usage

```bash
contract-cli decrypt-attestation [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to encrypted attestation file |
| `--priv` | string | Yes | Path to private key used for decryption |
| `--out` | string | No | Path to save decrypted attestation records |
| `-h, --help` | - | No | Display help information |

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

---

### download-certificate

Download encryption certificates from the IBM Hyper Protect Repository. Retrieves the latest or specific versions of HPCR encryption certificates required for contract encryption and workload deployment.

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

### encrypt

Generate a signed and encrypted contract for Hyper Protect deployment. Supports optional contract expiry for enhanced security.

#### Usage

```bash
contract-cli encrypt [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to unencrypted Hyper Protect contract YAML file |
| `--priv` | string | No* | Path to private key for signing |
| `--cert` | string | No | Path to encryption certificate (uses latest if not specified) |
| `--os` | string | No | Target Hyper Protect platform: `hpvs`, `hpcr-rhvs`, or `hpcc-peerpod` (default: `hpvs`) |
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

**For HPCR-RHVS:**
```bash
contract-cli encrypt \
  --in contract.yaml \
  --priv private.pem \
  --os hpcr-rhvs
```

**For HPCC Peer Pods:**
```bash
contract-cli encrypt \
  --in contract.yaml \
  --priv private.pem \
  --os hpcc-peerpod
```

---

### encrypt-string

Encrypt strings using the Hyper Protect encryption format. Output format: `hyper-protect-basic.<encrypted-password>.<encrypted-string>`. Use this to encrypt sensitive data like passwords or API keys for contracts.

#### Usage

```bash
contract-cli encrypt-string [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | String data to encrypt |
| `--format` | string | No | Input data format (text or json) |
| `--cert` | string | No | Path to encryption certificate (uses latest if not specified) |
| `--os` | string | No | Target Hyper Protect platform: `hpvs`, `hpcr-rhvs`, or `hpcc-peerpod` (default: `hpvs`) |
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
| `--in` | string | Yes | Path to download-certificate JSON output |
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

---

### image

Retrieve Hyper Protect Container Runtime (HPCR) image details from IBM Cloud. Parses image information from IBM Cloud API, CLI, or Terraform output to extract image ID, name, checksum, and version. Supports filtering by specific HPCR version.

#### Usage

```bash
contract-cli image [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to IBM Cloud images JSON (from API, CLI, or Terraform) |
| `--version` | string | No | Specific HPCR version to retrieve (returns latest if not specified) |
| `--format` | string | No | Output format for data (json, yaml, or text) |
| `--out` | string | No | Path to save HPCR image details |
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

---

### validate-contract

Validate an unencrypted contract against the Hyper Protect schema. Checks contract structure, required fields, and data types before encryption to help catch errors early in the development process.

#### Usage

```bash
contract-cli validate-contract [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to unencrypted Hyper Protect contract YAML file |
| `--os` | string | No | Target Hyper Protect platform: `hpvs`, `hpcr-rhvs`, or `hpcc-peerpod` (default: `hpvs`) |
| `-h, --help` | - | No | Display help information |

#### Examples

**Validate HPVS contract:**
```bash
contract-cli validate-contract --in contract.yaml --os hpvs
```

**Validate HPCR-RHVS contract:**
```bash
contract-cli validate-contract --in contract.yaml --os hpcr-rhvs
```

**Validate HPCC Peer Pods contract:**
```bash
contract-cli validate-contract --in contract.yaml --os hpcc-peerpod
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
| `--in` | string | Yes | Path to network-config YAML file |
| `-h, --help` | - | No | Display help information |

#### Examples

**Validate network configuration:**
```bash
contract-cli validate-network --in network-config.yaml
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
| `--in` | string | Yes | Path to encryption certificate file |
| `-h, --help` | - | No | Display help information |

#### Examples

**Validate encryption certifacte configuration:**
```bash
contract-cli validate-encryption-certificate --in encryption-cert.crt
```

---

### initdata
Create initdata annotation from signed and encrypted contract for Hyper Protect Confidential Containers PeerPod solution

#### Usage

```bash
contract-cli initdata [flags]
```

#### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--in` | string | Yes | Path to signed & encrypted contract YAML file |
| `--out` | string | No | Path to store gzipped & encoded initdata value |
| `-h, --help` | - | No | Display help information |

#### Examples

**Create Hpcc Initdata from signed & encrypted contract**
```bash
contract-cli initdata --in signed_encrypted_contract.yaml
```

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
- Check contract structure matches HPVS/HPCR requirements
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

- **[Simple Contract](../samples/simple_contract.yaml)** - Basic contract structure
- **[Contract with Expiry](../samples/contract_expiry.yaml)** - Contract with expiration
- **[Attestation Records](../samples/attestation/)** - Example attestation files
- **[Network Configuration](../samples/network/)** - Network config examples
- **[Docker Compose](../samples/tgz/)** - Compose file examples
- **[Signed & Encrypted Contract](../samples/hpcc/signed-encrypt-hpcc.yaml)** - Signed & Encrypted hpcc contract

---

## Additional Resources

- **[Main README](../README.md)** - Project overview and quick start
- **[Contributing Guide](../CONTRIBUTING.md)** - How to contribute
- **[Security Policy](../SECURITY.md)** - Security best practices
- **[IBM Hyper Protect Documentation](https://www.ibm.com/docs/en/hpvs/2.2.x)** - Official IBM docs

---

**Need Help?**

- [Open an issue](https://github.com/ibm-hyper-protect/contract-cli/issues/new/choose)
- [Ask a question](https://github.com/ibm-hyper-protect/contract-cli/discussions)
- Check the [troubleshooting](#troubleshooting) section above
