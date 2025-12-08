# Contract CLI

[![contract-cli CI](https://github.com/ibm-hyper-protect/contract-cli/actions/workflows/build.yaml/badge.svg)](https://github.com/ibm-hyper-protect/contract-cli/actions/workflows/build.yaml)
[![Latest Release](https://img.shields.io/github/v/release/ibm-hyper-protect/contract-cli?include_prereleases)](https://github.com/ibm-hyper-protect/contract-cli/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/ibm-hyper-protect/contract-cli)](https://goreportcard.com/report/ibm-hyper-protect/contract-cli)
[![GitHub All Releases](https://img.shields.io/github/downloads/ibm-hyper-protect/contract-cli/total.svg)](https://github.com/ibm-hyper-protect/contract-cli/releases/latest)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A command-line tool for automating the provisioning and management of IBM Hyper Protect confidential computing workloads.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Documentation](#documentation)
- [Supported Platforms](#supported-platforms)
- [Examples](#examples)
- [Related Projects](#related-projects)
- [Contributing](#contributing)
- [License](#license)
- [Support](#support)

## Overview

The Contract CLI automates the provisioning of IBM Hyper Protect confidential computing solutions:

- **Hyper Protect Virtual Servers (HPVS)** - Secure virtual servers on IBM Cloud
- **Hyper Protect Container Runtime (HPCR)** for RedHat Virtualization (RHVS)
- **Hyper Protect Confidential Container (HPCC)** for Red Hat OpenShift Peer Pods

This CLI tool leverages [ibm-hyper-protect/contract-go](https://github.com/ibm-hyper-protect/contract-go) for all cryptographic operations and contract management functionality, providing a user-friendly command-line interface for deploying workloads in secure enclaves on IBM LinuxONE.

### What are Hyper Protect Services?

IBM Hyper Protect services provide confidential computing capabilities that protect data in use by leveraging Secure Execution feature of Z.

Learn more:

- [Confidential computing with LinuxONE](https://cloud.ibm.com/docs/vpc?topic=vpc-about-se)
- [IBM Hyper Protect Virtual Servers](https://www.ibm.com/docs/en/hpvs/2.2.x)
- [IBM Hyper Protect Confidential Container for Red Hat OpenShift](https://www.ibm.com/docs/en/hpcc/1.1.x)

## Features

- **Attestation Management**
  - Decrypt encrypted attestation records

- **Certificate Operations**
  - Download HPVS encryption certificates from IBM Cloud
  - Extract specific encryption certificates by version
  - Validate expiry of encryption certificate

- **Contract Generation**
  - Generate Base64-encoded data from text, JSON, and docker compose / podman play archives
  - Create signed and encrypted contracts
  - Support contract expiry with CA certificates
  - Validate contract schemas

- **Archive Management**
  - Generate Base64 tar archives of `docker-compose.yaml` or `pods.yaml`
  - Support encrypted base64 tar generation

- **String Encryption**
  - Encrypt strings using Hyper Protect format
  - Support both text and JSON input

- **Image Selection**
  - Retrieve latest HPCR image details from IBM Cloud
  - Filter images by semantic versioning

- **Network Validation**
  - Validate network-config schemas for on-premise deployments
  - Support HPVS, HPCR RHVS, and HPCC Peer Pod configurations

## Installation

Download the CLI tool for your operating system from the [releases page](https://github.com/ibm-hyper-protect/contract-cli/releases/latest).

### Supported Platforms

The CLI is available for the following platforms:

| OS | Architecture | Binary |
|----|--------------|--------|
| Linux | amd64, arm64, s390x, ppc64le | `contract-cli-linux-*` |
| macOS | amd64, arm64 | `contract-cli-darwin-*` |
| Windows | amd64, arm64 | `contract-cli-windows-*.exe` |

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
  --contract contract.yaml \
  --key private.pem \
  --output encrypted-contract.yaml
```

### Download and Use Encryption Certificates

```bash
# Download the latest encryption certificates
contract-cli download-certificate \
  --output certificates.json

# Extract a specific version
contract-cli get-certificate \
  --input certificates.json \
  --version "1.0.0" \
  --output cert-1.0.0.crt
```

### Validate encryption certitifacate

```bash
# validate downloaded encryption certificate
contract-cli validate-encryption-certificate \
  --in encryption-cert.crt
```

### Validate a Contract Before Encryption

```bash
# Validate contract schema
contract-cli validate-contract \
  --contract contract.yaml \
  --type hpvs
```

## Usage

```bash
$ contract-cli --help
Contract CLI automates contract generation and management for IBM Hyper Protect services.

Supports:
  - Hyper Protect Virtual Servers (HPVS) for VPC
  - Hyper Protect Container Runtime (HPCR) for RHVS
  - Hyper Protect Confidential Container (HPCC) Peer Pods

Documentation: https://github.com/ibm-hyper-protect/contract-cli/blob/main/docs/README.md

Usage:
  contract-cli [flags]
  contract-cli [command]

Available Commands:
  base64                          Encode input as Base64
  base64-tgz                      Create Base64 tar archive of container configurations
  decrypt-attestation             Decrypt encrypted attestation records
  download-certificate            Download encryption certificates
  encrypt                         Generate signed and encrypted contract
  encrypt-string                  Encrypt string in Hyper Protect format
  get-certificate                 Extract specific certificate version from download output
  help                            Help about any command
  image                           Get HPCR image details from IBM Cloud
  validate-contract               Validate contract schema
  validate-encryption-certificate validate encryption certificate
  validate-network                Validate network configuration schema

Flags:
  -h, --help      help for contract-cli
  -v, --version   version for contract-cli

Use "contract-cli [command] --help" for more information about a command.
```

## Documentation

Comprehensive documentation is available at:

- **[User Documentation](docs/README.md)** - Detailed command reference and usage examples
- **[Command Reference](docs/README.md)** - Complete guide for all CLI commands

## Supported Platforms

| Platform | Description | Support Status |
|----------|-------------|----------------|
| HPVS | Hyper Protect Virtual Servers | Supported |
| HPCR-RHVS | Hyper Protect Container Runtime for Red Hat Virtualization | Supported |
| HPCC-PeerPod | Hyper Protect Confidential Container Peer Pods | Supported |

## Examples

The [`samples/`](samples/) directory contains example configurations:

- [Simple Contract](samples/simple_contract.yaml)
- [Contract with Expiry](samples/contract_expiry.yaml)
- [Attestation Records](samples/attestation/)
- [Network Configuration](samples/network/)
- [Docker Compose Examples](samples/tgz/)

## Related Projects

This CLI tool is part of the IBM Hyper Protect ecosystem:

| Project | Description |
|---------|-------------|
| [contract-go](https://github.com/ibm-hyper-protect/contract-go) | Core Go library for Hyper Protect contracts |
| [terraform-provider-hpcr](https://github.com/ibm-hyper-protect/terraform-provider-hpcr) | Terraform provider for Hyper Protect contracts |
| [k8s-operator-hpcr](https://github.com/ibm-hyper-protect/k8s-operator-hpcr) | Kubernetes operator for contract management |
| [linuxone-vsi-automation-samples](https://github.com/ibm-hyper-protect/linuxone-vsi-automation-samples) | Terraform examples for HPVS and HPCR RHVS |
| [hyper-protect-virtual-server-samples](https://github.com/ibm-hyper-protect/hyper-protect-virtual-server-samples) | HPVS feature samples and scripts |

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
