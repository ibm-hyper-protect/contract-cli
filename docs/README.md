# Hyper Protect Contract CLI

## Introduction

The CLI has been developed to automate the process for generating contracts for provisioning Hyper Protect Virtual Servers for VPC and Hyper Protect Container Runtime.

## Usage

### Base64

This feature will help you to generate base64 of your plain text or JSON input.

```bash
$ contract-cli base64 --help
Generate base64 of input text

Usage:
  contract-cli base64 [flags]

Flags:
      --format string   Format of input data (options: text/json) (default "text")
  -h, --help            help for base64
      --in string       Input data that needs to be converted to string
      --out string      Path to store Base64 output
```

To generate base64 of plain text.
```bash
$ contract-cli base64 --in <input-text> --format plain
```

To generate base64 of JSON data.
```bash
$ contract-cli base64 --in <input-json-data> --format json
```

If you want to redirect the result to a file.
```bash
$ contract-cli base64 --in <input-text> --format plain --out <path-to-output-file>
```

The following is an example to generate base64.
```bash
$ contract-cli base64 --in sampleText --format plain # For plain text
$ contract-cli base64 --in {"type": "workload"} --format json # For JSON text
```
