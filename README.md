# Contract CLI

[![contract-cli CI](https://github.com/ibm-hyper-protect/contract-cli/actions/workflows/build.yaml/badge.svg)](https://github.com/ibm-hyper-protect/contract-cli/actions/workflows/build.yaml)
[![Latest Release](https://img.shields.io/github/v/release/ibm-hyper-protect/contract-cli?include_prereleases)](https://github.com/ibm-hyper-protect/contract-cli/releases/latest)
[![User Documentation](https://img.shields.io/badge/User%20Documentation-GitHub%20Pages-blue.svg)](https://ibm-hyper-protect.github.io/contract-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/ibm-hyper-protect/contract-cli)](https://goreportcard.com/report/ibm-hyper-protect/contract-cli)
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

- **Certificate Operations**
  - Download encryption certificates from IBM Cloud
  - Extract specific encryption certificates by version
  - Validate expiry of encryption certificate

- **Contract Generation**
  - Generate Base64-encoded data from text, JSON, and docker compose / podman play archives
  - Create signed and signed & encrypted contracts
  - Support contract expiry with CA certificates
  - Validate contract schemas
  - Create Gzipped & Encoded initdata for IBM Confidential Computing Containers Peer Pod

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
  - Support HPVS, HPCR RHVS, and HPCC Peer Pod configurations


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
# Validate contract schema
contract-cli validate-contract \
  --in contract.yaml \
  --os hpvs
```

### Create initdata annotation from signed & encrypted contract

```bash
# Create initdata annotation
contract-cli initdata \
  --in signed_encrypted_contract.yaml
```

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
  decrypt-attestation             Decrypt encrypted attestation records
  download-certificate            Download encryption certificates
  encrypt                         Generate signed and encrypted contract
  encrypt-string                  Encrypt string in IBM Confidential Computing format
  get-certificate                 Extract specific certificate version from download output
  help                            Help about any command
  image                           Get IBM Confidential Computing Container Runtime image details from IBM Cloud
  initdata                        Gzip and Encoded initdata annotation
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
