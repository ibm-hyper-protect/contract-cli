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

# Contract CLI

[![Build & Test](https://github.com/ibm-hyper-protect/contract-cli/actions/workflows/build.yml/badge.svg)](https://github.com/ibm-hyper-protect/contract-cli/actions/workflows/build.yml)
[![Release](https://github.com/ibm-hyper-protect/contract-cli/actions/workflows/release.yml/badge.svg)](https://github.com/ibm-hyper-protect/contract-cli/actions/workflows/release.yml)
[![Latest Release](https://img.shields.io/github/v/release/ibm-hyper-protect/contract-cli?include_prereleases)](https://github.com/ibm-hyper-protect/contract-cli/releases/latest)
[![User Documentation](https://img.shields.io/badge/User%20Documentation-GitHub%20Pages-blue.svg)](https://ibm-hyper-protect.github.io/contract-cli)
[![GitHub All Releases](https://img.shields.io/github/downloads/ibm-hyper-protect/contract-cli/total.svg)](https://github.com/ibm-hyper-protect/contract-cli/releases/latest)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A command-line tool for automating the provisioning and management of IBM Confidential Computing workloads on IBM Z and LinuxONE.

## Table of Contents

- [Overview](#overview)
- [Who Is This For?](#who-is-this-for)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Documentation](#documentation)
- [Examples](#examples)
- [Related Projects](#related-projects)
- [Contributing](#contributing)
- [License](#license)
- [Support](#support)

## Overview

The Contract CLI automates the provisioning of IBM Confidential Computing solutions:

- **IBM Confidential Computing Container Runtime** (formerly known as Hyper Protect Virtual Servers) — Deploy confidential computing workloads on IBM Z and LinuxONE using IBM Secure Execution for Linux
- **IBM Confidential Computing Container Runtime for Red Hat Virtualization Solutions** (formerly known as Hyper Protect Container Runtime for Red Hat Virtualization Solutions) — Purpose-built for hosting critical, centralized services within tightly controlled virtualized environments on IBM Z
- **IBM Confidential Computing Containers for Red Hat OpenShift Container Platform** (formerly known as IBM Hyper Protect Confidential Container for Red Hat OpenShift Container Platform) — Deploy isolated workloads using IBM Secure Execution for Linux, integrated with Red Hat OpenShift Container Platform

This CLI tool leverages [ibm-hyper-protect/contract-go](https://github.com/ibm-hyper-protect/contract-go) for all cryptographic operations and contract management functionality, providing a user-friendly command-line interface for deploying workloads in secure enclaves on IBM Z and LinuxONE.

### What is IBM Confidential Computing?

IBM Confidential Computing services protect data in use by leveraging the IBM Secure Execution for Linux feature on IBM Z and LinuxONE hardware. Each deployment is configured through a **contract** — an encrypted YAML definition file that specifies workload, environment, and attestation settings.

Learn more:

- [Confidential computing with LinuxONE](https://cloud.ibm.com/docs/vpc?topic=vpc-about-se)
- [IBM Confidential Computing Container Runtime](https://www.ibm.com/docs/en/cccr/2.2.x)
- [IBM Confidential Computing Container Runtime for Red Hat Virtualization Solutions](https://www.ibm.com/docs/en/ccrv/1.1.x)
- [IBM Confidential Computing Containers for Red Hat OpenShift](https://www.ibm.com/docs/en/ccro/1.1.x)

### Who Is This For?

This CLI is for **developers, DevOps engineers, and platform teams** who need to generate, sign, and encrypt deployment contracts for IBM Confidential Computing services. Common use cases include:

- **Scripting & Automation** — Generate contracts in CI/CD pipelines
- **Certificate Management** — Download and verify IBM encryption certificates
- **Attestation** — Decrypt and verify workload integrity records
- **Validation** — Validate contracts and network configurations before deployment

> **Go developers** who need programmatic access should use the [contract-go](https://github.com/ibm-hyper-protect/contract-go) library directly. For infrastructure-as-code workflows, see the [terraform-provider-hpcr](https://github.com/ibm-hyper-protect/terraform-provider-hpcr) Terraform provider.

## Features

- **Attestation Management**
  - Decrypt encrypted attestation records
  - Verify signature of decrypted attestation records against IBM attestation certificate
  - Support for password-protected private keys

- **Certificate Operations**
  - Download encryption certificates from IBM Cloud
  - Extract specific encryption certificates by version
  - Validate expiry of encryption certificate
  - **List available embedded encryption certificate versions** for all platforms or specific platform
  - **Use specific certificate encryption version** for encryption operations

- **Contract Generation**
  - Generate Base64-encoded data from text, JSON, and docker compose / podman play archives
  - Create signed and signed & encrypted contracts
  - **Generate contract YAML templates** for workload, env, or combined contract sections
  - Support contract expiry with CA certificates
  - Support for password-protected private keys in signing and encryption operations
  - **Specify encryption certificate version** for encryption operations with `--ver` flag
  - Validate contract schemas
  - Create Gzipped & Encoded initdata for IBM Confidential Computing Containers (Peer Pod and Baremetal solutions)

- **Sealed Secret Management**
  - Generate sealed secrets for CCCO workload and environment sections
  - Automatic key generation or use custom encryption/signing keys
  - Output sealed secrets with decryption and verification keys

- **Archive Management**
  - Generate Base64 tar archives of `docker-compose.yaml` or `pods.yaml`
  - Support encrypted base64 tar generation

- **String Encryption**
  - Encrypt strings using IBM Confidential Computing format
  - Support both text and JSON input

- **Image Selection**
  - Retrieve latest IBM Confidential Computing Container Runtime image details from IBM Cloud
  - Filter images by semantic versioning

- **Network Validation**
  - Validate network-config schemas for on-premise deployments
  - Support ccrt, ccrv, and ccco configurations


## Installation

### Homebrew (macOS / Linux)

```bash
brew tap ibm-hyper-protect/contract-cli https://github.com/ibm-hyper-protect/contract-cli
brew install contract-cli

# Install a specific version
brew install contract-cli@1.2.0
```

### Debian / Ubuntu (apt)

Download the `.deb` package from the [releases page](https://github.com/ibm-hyper-protect/contract-cli/releases/latest) and install:

```bash
# Download (replace VERSION and ARCH as needed)
curl -LO https://github.com/ibm-hyper-protect/contract-cli/releases/latest/download/contract-cli_VERSION_linux_amd64.deb

# Install
sudo dpkg -i contract-cli_*.deb
```

### Fedora / RHEL / Rocky / Alma (dnf/yum)

Download the `.rpm` package from the [releases page](https://github.com/ibm-hyper-protect/contract-cli/releases/latest) and install:

```bash
# Download (replace VERSION and ARCH as needed)
curl -LO https://github.com/ibm-hyper-protect/contract-cli/releases/latest/download/contract-cli_VERSION_linux_amd64.rpm

# Install
sudo rpm -i contract-cli_*.rpm
```

### Docker

```bash
# Run using the Docker image
docker run --rm ghcr.io/ibm-hyper-protect/contract-cli --version

# Example: encrypt a contract
docker run --rm -v "$(pwd):/work" -w /work \
  ghcr.io/ibm-hyper-protect/contract-cli encrypt \
  --in contract.yaml --priv private.pem --out encrypted.yaml
```

Available tags:
- `ghcr.io/ibm-hyper-protect/contract-cli:latest` — latest release
- `ghcr.io/ibm-hyper-protect/contract-cli:<version>` — specific version

Multi-architecture support: `amd64`, `arm64`, `s390x`, `ppc64le`.

### Windows (Winget)

> **Note:** Winget package submission is in progress. In the meantime, download the Windows binary from the [releases page](https://github.com/ibm-hyper-protect/contract-cli/releases/latest).

```powershell
winget install ibmcc-contract-cli

# Install a specific version
winget install ibmcc-contract-cli --version 1.2.0
```

### Direct Binary Download

Download the CLI tool for your operating system from the [releases page](https://github.com/ibm-hyper-protect/contract-cli/releases/latest).

#### Verify Download (Recommended)

After downloading, verify the binary using the checksum file:

```bash
# Download the checksums file
curl -LO https://github.com/ibm-hyper-protect/contract-cli/releases/latest/download/checksums.txt

# Verify (Linux/macOS)
sha256sum --check checksums.txt --ignore-missing

# Verify cosign signature (if cosign is installed)
cosign verify-blob \
  --key https://github.com/ibm-hyper-protect/contract-cli/releases/latest/download/checksums.txt.pem \
  --signature https://github.com/ibm-hyper-protect/contract-cli/releases/latest/download/checksums.txt.sig \
  checksums.txt
```

### Supported Platforms

The CLI is available for the following platforms:

| OS | Architecture | Package Formats |
|----|--------------|-----------------|
| Linux | amd64, arm64, s390x, ppc64le | Binary, `.deb`, `.rpm`, `.tar.gz`, Docker |
| macOS | amd64, arm64 | Binary, `.tar.gz`, Homebrew |
| Windows | amd64, arm64 | Binary, `.zip`, Winget |

### Prerequisites

- **OpenSSL** - Required for encryption operations
  - On Linux: `apt-get install openssl` or `yum install openssl`
  - On macOS: `brew install openssl`
  - On Windows: [Download OpenSSL](https://slproweb.com/products/Win32OpenSSL.html)

#### Optional: Custom OpenSSL Path

If OpenSSL is not in your system PATH, set the `OPENSSL_BIN` environment variable:

```bash
# Linux/macOS
export OPENSSL_BIN=/usr/bin/openssl

# Windows (PowerShell)
$env:OPENSSL_BIN="C:\Program Files\OpenSSL-Win64\bin\openssl.exe"

# Docker
docker run --rm -e OPENSSL_BIN=/usr/bin/openssl \
  -v "$(pwd):/work" -w /work \
  ghcr.io/ibm-hyper-protect/contract-cli encrypt \
  --in contract.yaml --priv private.pem --out encrypted.yaml
```

## Quick Start

### Generate a Contract Template

```bash
# Generate a combined contract template (workload + env) to stdout
contract-cli contract-template

# Generate a workload-only template for HPVS
contract-cli contract-template --type workload --os hpvs

# Generate an env-only template for CCRT
contract-cli contract-template --type env --os ccrt

# Generate an workload-only template for ccco-peerpod
contract-cli contract-template --type workload --os ccco-peerpod

# Save the full contract template to a file
contract-cli contract-template --out contract-template.yaml
```

### Generate a Signed and Encrypted Contract

```bash
# Create a contract YAML file
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

# Generate RSA key pair
openssl genrsa -out private.pem 4096

# Generate signed and encrypted contract
contract-cli encrypt \
  --in contract.yaml \
  --priv private.pem \
  --out encrypted-contract.yaml

# Or with a password-protected private key
openssl genrsa -aes256 -out private-encrypted.pem 4096
contract-cli encrypt \
  --in contract.yaml \
  --priv private-encrypted.pem \
  --password "your-secure-password" \
  --out encrypted-contract.yaml
```

### List Available Certificate Versions

```bash
# List all available embedded certificate versions (JSON format by default)
contract-cli list-encryptioncert-versions

# List versions in YAML format
contract-cli list-encryptioncert-versions --format yaml

# List versions for a specific platform in JSON
contract-cli list-encryptioncert-versions --os ccrt --format json

# List versions for a specific platform in YAML
contract-cli list-encryptioncert-versions --os ccrv --format yaml

# List HPVS certificate versions
contract-cli list-encryptioncert-versions --os hpvs --format json

# Save output to file
contract-cli list-encryptioncert-versions --os ccco --format yaml --out ccco-versions.yaml
```

### Use Specific Certificate Version for Encryption

```bash
# Encrypt contract with a specific certificate version
contract-cli encrypt \
  --in contract.yaml \
  --priv private.pem \
  --os ccrt \
  --ver 26.2.0 \
  --out encrypted-contract.yaml

# Encrypt string with specific certificate version
contract-cli encrypt-string \
  --in "sensitive data" \
  --os ccrv \
  --ver 25.11.0

# Create encrypted base64 tar with specific certificate version for HPVS
contract-cli base64-tgz \
  --in docker-compose.yaml \
  --os hpvs \
  --ver 26.5.0 \
  --output encrypt
```

### Download and Use Encryption Certificates

```bash
# Download the latest encryption certificates
contract-cli download-certificate \
  --out certificates.json

# Extract a specific version
contract-cli get-certificate \
  --in certificates.json \
  --version "1.0.0" \
  --out cert-1.0.0.crt
```

### Validate Encryption Certificate

```bash
# validate downloaded encryption certificate
contract-cli validate-encryption-certificate \
  --in encryption-cert.crt
```

### Validate a Contract Before Encryption

```bash
# Validate full contract (both workload and env sections)
contract-cli validate-contract \
  --in contract.yaml \
  --os ccrt

# Validate only the workload section
contract-cli validate-contract \
  --in contract.yaml \
  --os ccrt \
  --type workload

# Validate only the env section
contract-cli validate-contract \
  --in contract.yaml \
  --os ccrt \
  --type env
```

### Create initdata annotation from signed & encrypted contract

```bash
# Create initdata annotation for Peer Pod solution
contract-cli initdata \
  --in signed_encrypted_contract.yaml

# Create initdata annotation for Baremetal solution with SE header binary
contract-cli initdata \
  --in signed_encrypted_contract.yaml \
  --sehdr se-header.bin \
  --out initdata.txt
```

### Using Password-Protected Private Keys

The CLI supports password-protected private keys for enhanced security:

```bash
# Create a password-protected private key
openssl genrsa -aes256 -out private-encrypted.pem 4096

# Decrypt attestation with password-protected key
contract-cli decrypt-attestation \
  --in encrypted-attestation.txt \
  --priv private-encrypted.pem \
  --password "your-secure-password" \
  --out decrypted-attestation.txt

# Sign contract with password-protected key
contract-cli sign-contract \
  --in contract.yaml \
  --priv private-encrypted.pem \
  --password "your-secure-password" \
  --out signed-contract.yaml

# Encrypt contract with password-protected key
contract-cli encrypt \
  --in contract.yaml \
  --priv private-encrypted.pem \
  --password "your-secure-password" \
  --out encrypted-contract.yaml
```

**Security Note**: For production use, consider using environment variables or secure secret management systems instead of passing passwords directly on the command line:

```bash
# Using environment variable
export PRIVATE_KEY_PASSWORD="your-secure-password"
contract-cli encrypt \
  --in contract.yaml \
  --priv private-encrypted.pem \
  --password "$PRIVATE_KEY_PASSWORD" \
  --out encrypted-contract.yaml

# Or read from a secure file
contract-cli encrypt \
  --in contract.yaml \
  --priv private-encrypted.pem \
  --password "$(cat /secure/path/password.txt)" \
  --out encrypted-contract.yaml
```

### Generate Sealed Secrets for CCCO

```bash
# Generate sealed secret for environment variables
contract-cli sealed-secret \
  --in "value123" \
  --type env \
  --out sealed-secret.txt

# Generate sealed secret for workload data
contract-cli sealed-secret \
  --in workload-secret-data \
  --type workload \
  --out sealed-workload.txt

# Generate sealed secret from file
contract-cli sealed-secret \
  --in secrets.txt \
  --type env \
  --out sealed-secret.txt

# Generate sealed secret with custom encryption and signing keys

openssl genrsa -out encryption.pem 2048
openssl genrsa -out signing.pem 2048

contract-cli sealed-secret \
  --in "value123" \
  --type env \
  --encryptionkey encryption.pem \
  --signingkey signing.pem \
  --out sealed-secret.txt

# Read secret from stdin
echo "value123" | contract-cli sealed-secret \
  --in - \
  --type env
```

The sealed secret output includes:
- The sealed secret data (for use in contract)
- `SECRET_DECRYPTION_KEY` - Private key for decryption (keep secure)
- `SECRET_VERIFICATION_KEY` - Public key for verification

## Usage

```bash
$ contract-cli --help
Contract CLI automates contract generation and management for IBM Confidential Computing services.

Supports:
  - IBM Confidential Computing Container Runtime
  - IBM Confidential Computing Container Runtime for Red Hat Virtualization Solutions
  - IBM Confidential Computing Containers for Red Hat OpenShift Container Platform

Documentation: https://ibm-hyper-protect.github.io/contract-cli/

Usage:
  contract-cli [flags]
  contract-cli [command]

Available Commands:
  base64                          Encode input as Base64
  base64-tgz                      Create Base64 tar archive of container configurations
  contract-template               Generate a contract template
  decrypt-attestation             Decrypt encrypted attestation records
  download-certificate            Download encryption certificates
  encrypt                         Generate signed and encrypted contract
  encrypt-string                  Encrypt string in IBM Confidential Computing format
  get-certificate                 Extract specific certificate version from download output
  help                            Help about any command
  sealed-secret                   Generate sealed secret for CCCO
  image                           Get IBM Confidential Computing Container Runtime image details from IBM Cloud
  initdata                        Gzip and Encoded initdata annotation
  list-encryptioncert-versions    List available encryption certificate versions
  sign-contract                   Sign an encrypted contract
  validate-contract               Validate contract schema
  validate-encryption-certificate Validate encryption certificate
  validate-network                Validate network configuration schema

Flags:
  -h, --help      help for contract-cli
  -v, --version   version for contract-cli

Use "contract-cli [command] --help" for more information about a command.
```

## Documentation

Comprehensive documentation is available at:

- **[Command Reference & User Guide](https://ibm-hyper-protect.github.io/contract-cli/)** - Detailed command reference, workflows, and usage examples
- **[Changelog](CHANGELOG.md)** - Release history and version notes

## Examples

The [`samples/`](samples/) directory contains example configurations:

- [Contract](samples/contract.yaml)
- [Contract with Expiry](samples/contract-expiry/)
- [Attestation Records](samples/attestation/)
- [Certificate Examples](samples/certificate/)
- [Network Configuration](samples/network/)
- [Docker Compose Examples](samples/tgz/)
- [Signed & Encrypted Contract](samples/hpcc/signed-encrypt-hpcc.yaml)
- [Contract Signing](samples/sign/)

## Related Projects

This CLI tool is part of the IBM Confidential Computing ecosystem:

| Project | Description | When to Use |
|---------|-------------|-------------|
| [contract-go](https://github.com/ibm-hyper-protect/contract-go) | Core Go library for IBM Confidential Computing contracts | When you need programmatic access from Go code |
| [terraform-provider-hpcr](https://github.com/ibm-hyper-protect/terraform-provider-hpcr) | Terraform provider for IBM Confidential Computing contracts | When managing infrastructure as code with Terraform |
| [k8s-operator-hpcr](https://github.com/ibm-hyper-protect/k8s-operator-hpcr) | Kubernetes operator for contract management | When managing contracts in Kubernetes clusters |
| [linuxone-vsi-automation-samples](https://github.com/ibm-hyper-protect/linuxone-vsi-automation-samples) | Terraform & CLI examples for IBM Confidential Computing deployments | For deployment automation reference |
| [hyper-protect-virtual-server-samples](https://github.com/ibm-hyper-protect/hyper-protect-virtual-server-samples) | IBM Confidential Computing feature samples and scripts | For feature samples and reference scripts |

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details on:

- Opening issues
- Submitting pull requests
- Code style and conventions
- Testing requirements

Please also read our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

### Reporting Issues

We use GitHub issue templates to help us understand and address your concerns efficiently:

- **[Report a Bug](https://github.com/ibm-hyper-protect/contract-cli/issues/new?template=bug_report.yml)** - Found a bug? Let us know!
- **[Request a Feature](https://github.com/ibm-hyper-protect/contract-cli/issues/new?template=feature_request.yml)** - Have an idea for improvement?
- **[Ask a Question](https://github.com/ibm-hyper-protect/contract-cli/issues/new?template=question.yml)** - Need help using the CLI?

### Security

- **Security Vulnerabilities**: Report via [GitHub Security Advisories](https://github.com/ibm-hyper-protect/contract-cli/security/advisories/new) - **DO NOT** create public issues
- See our complete [Security Policy](SECURITY.md) for details

### Community

- **[Discussions](https://github.com/ibm-hyper-protect/contract-cli/discussions)** - General questions and community discussion
- **[Documentation](docs/README.md)** - Comprehensive CLI documentation
- **[Maintainers](MAINTAINERS.md)** - Current maintainer list and contact info

## Contributors

![Contributors](https://contrib.rocks/image?repo=ibm-hyper-protect/contract-cli)
