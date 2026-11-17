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

# Contract CLI Plugin Developer Guide

This guide explains how to create, release, and distribute a `contract-cli` plugin.

## Table of Contents

- [Overview](#overview)
- [Repository & Binary Naming](#repository--binary-naming)
- [GitHub Release Asset Naming](#github-release-asset-naming)
- [How Dispatch Works](#how-dispatch-works)
- [Handshake Requirement](#handshake-requirement)
- [Version Requirement](#version-requirement)
- [Minimal Go + Cobra Plugin Example](#minimal-go--cobra-plugin-example)
- [Local Development Workflow](#local-development-workflow)
- [User Quickstart](#user-quickstart)

---

## Overview

`contract-cli` supports a **binary plugin model** modelled on `ibmcloud CLI` and `kubectl`:

- Plugins are standalone executables named `contract-plugin-<name>`.
- They are stored in `~/.contract-cli/plugins/`.
- Any unknown sub-command is transparently delegated to the matching plugin binary.
- Plugins **cannot be executed directly** — they must be invoked through `contract-cli`.

Example: a plugin named `ccrv` is invoked as:

```sh
contract-cli ccrv <plugin-subcommands> [flags]
```

---

## Repository & Binary Naming

| Item | Convention |
|---|---|
| GitHub repository | `contract-plugin-<name>` under any GitHub org |
| Binary name (Linux / macOS) | `contract-plugin-<name>` |
| Binary name (Windows) | `contract-plugin-<name>.exe` |

The plugin short-name (`<name>`) becomes the sub-command in `contract-cli`. For example, a repo named `contract-plugin-ccrv` is invoked as `contract-cli ccrv`.

---

## GitHub Release Asset Naming

When `contract-cli plugin install` downloads a release, it looks for an asset named:

```
contract-plugin-<name>-<os>-<arch>[.exe]
```

| Placeholder | Allowed values |
|---|---|
| `<os>` | `linux`, `darwin`, `windows` |
| `<arch>` | `amd64`, `arm64`, `s390x`, `ppc64le` |

Example assets for a plugin named `ccrv`:

```
contract-plugin-ccrv-linux-amd64
contract-plugin-ccrv-linux-arm64
contract-plugin-ccrv-linux-s390x
contract-plugin-ccrv-linux-ppc64le
contract-plugin-ccrv-darwin-amd64
contract-plugin-ccrv-darwin-arm64
contract-plugin-ccrv-windows-amd64.exe
```

Using [GoReleaser](https://goreleaser.com/) is the recommended way to build and upload these assets automatically.

### Checksum File (required)

By default `contract-cli plugin install` and `plugin update` verify the downloaded binary against a SHA-256 checksum file published in the **same release**. One of the following must be present (checked in this order):

| Asset name | Notes |
|---|---|
| `<asset>.sha256` | Per-asset file, e.g. `contract-plugin-ccrv-linux-amd64.sha256` |
| `checksums.txt`, `SHA256SUMS`, `SHA256SUMS.txt`, or `*_checksums.txt` | Aggregate file (GoReleaser's default is `<project>_<version>_checksums.txt`) |

Format is the standard `sha256sum` output (`<64-hex-digest>  <asset-name>`, one per line). A per-asset file may also contain just the bare digest.

If no checksum file is found, the asset is not listed in it, or the digest does not match, the install fails and nothing is written to the plugin directory (an already-installed version is left untouched).

To skip verification (not recommended), pass `--insecure`:

```sh
contract-cli plugin install github.com/org/contract-plugin-ccrv --insecure
contract-cli plugin update ccrv --insecure
```

`--insecure` has no effect on local-path installs, which are never verified.

---

## How Dispatch Works

When `contract-cli ccrv generate --env prod` is run:

1. Cobra checks whether `ccrv` is a registered built-in command → **No**.
2. `contract-cli` looks for `~/.contract-cli/plugins/contract-plugin-ccrv`.
3. If found and executable, it calls `exec contract-plugin-ccrv generate --env prod` with:
   - All remaining arguments passed through unchanged.
   - `stdin`, `stdout`, and `stderr` inherited.
   - The environment variable `CONTRACT_CLI_PLUGIN=1` injected.
   - The plugin's exit code propagated back to the caller.
4. If no matching binary exists, `contract-cli` prints the standard "unknown command" error.

---

## Handshake Requirement

**Every plugin must reject direct execution.**

`contract-cli` injects `CONTRACT_CLI_PLUGIN=1` into the child process environment on every managed dispatch. Your plugin must check for this variable at the very top of `main()` and exit with an error if it is absent.

### Required boilerplate

```go
package main

import (
    "fmt"
    "os"

    "github.com/ibm-hyper-protect/contract-cli/lib/plugin"
)

func main() {
    // Reject direct execution — must be first.
    if err := plugin.CheckHandshake(); err != nil {
        fmt.Fprintln(os.Stderr, "Error:", err)
        os.Exit(1)
    }

    // ... rest of your plugin
}
```

When run directly (without `contract-cli`), the plugin prints:

```
Error: this plugin must be run through contract-cli
```

and exits with code `1`.

> **Import path**: `github.com/ibm-hyper-protect/contract-cli/lib/plugin`
> The constants `plugin.HandshakeEnvVar` (`"CONTRACT_CLI_PLUGIN"`) and `plugin.HandshakeEnvVal` (`"1"`) are also exported if you need to check them directly.

---

## Version Requirement

`contract-cli plugin list` calls each plugin binary with `--version` to display its version. Your plugin **must** implement a `--version` flag that prints a **single line** to stdout.

**Never hardcode the version string.** Declare it as a package-level variable and inject it at build time via `-ldflags`:

```go
// version is set at build time: go build -ldflags "-X 'main.version=v1.2.0'"
var version = "dev"

var rootCmd = &cobra.Command{
    Use:     "contract-plugin-ccrv",
    Version: version,   // resolved from the ldflags-injected var, not a literal
}
```

Cobra uses `rootCmd.Version` to automatically handle `--version`, printing:

```
contract-plugin-ccrv version v1.2.0
```

`contract-cli plugin list` captures the first line of that output and displays `v1.2.0`.

---

## Minimal Go + Cobra Plugin Example

The following is a complete, working `contract-plugin-ccrv` skeleton:

```go
package main

import (
    "fmt"
    "os"

    "github.com/ibm-hyper-protect/contract-cli/lib/plugin"
    "github.com/spf13/cobra"
)

// version is injected at build time via:
//   go build -ldflags "-X 'main.version=$(VERSION)'"
// Defaults to "dev" for local builds.
var version = "dev"

var rootCmd = &cobra.Command{
    Use:     "contract-plugin-ccrv",
    Short:   "CCRV plugin for contract-cli",
    Version: version,
}

var generateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Generate a CCRV contract",
    Run: func(cmd *cobra.Command, args []string) {
        env, _ := cmd.Flags().GetString("env")
        fmt.Printf("Generating CCRV contract for env: %s\n", env)
    },
}

func init() {
    generateCmd.Flags().String("env", "dev", "Target environment")
    rootCmd.AddCommand(generateCmd)
}

func main() {
    // Reject direct execution — must be first.
    if err := plugin.CheckHandshake(); err != nil {
        fmt.Fprintln(os.Stderr, "Error:", err)
        os.Exit(1)
    }

    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

Build and test locally:

```sh
go build -ldflags "-X 'main.version=dev'" -o contract-plugin-ccrv .
cp contract-plugin-ccrv ~/.contract-cli/plugins/
contract-cli ccrv --help
contract-cli ccrv generate --env prod
```

---

## Makefile

Add a `Makefile` to your plugin repo that mirrors the `contract-cli` pattern:

```makefile
APP_NAME   := contract-plugin-ccrv
OUTPUT     := build
VERSION    ?= dev
BUILD_DATE ?= $(shell date -u)

build:
	CGO_ENABLED=0 go build \
	    -ldflags "-X 'main.version=$(VERSION)'" \
	    -o $(OUTPUT)/$(APP_NAME) .

test:
	go test ./... -v

fmt:
	go fmt ./...
```

Build with a version:

```sh
make build VERSION=v1.2.0
```

---

## GoReleaser Configuration

Use [GoReleaser](https://goreleaser.com/) to automatically build and upload the correct platform assets to GitHub Releases. The asset naming convention (`contract-plugin-ccrv-linux-amd64` etc.) **must match exactly** what `contract-cli plugin install` looks for.

```yaml
# .goreleaser.yaml
project_name: contract-plugin-ccrv

builds:
  - id: contract-plugin-ccrv
    binary: contract-plugin-ccrv
    main: .
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
      - s390x
      - ppc64le
    ldflags:
      - -X 'main.version={{ .Version }}'

archives:
  - format: binary
    name_template: "{{ .Binary }}-{{ .Os }}-{{ .Arch }}"

# Publishes <project>_<version>_checksums.txt (SHA-256), which plugin install verifies.
checksum:
  algorithm: sha256
```

On a `git tag v1.2.0` + `goreleaser release`, this uploads:

```
contract-plugin-ccrv-linux-amd64
contract-plugin-ccrv-linux-arm64
contract-plugin-ccrv-linux-s390x
contract-plugin-ccrv-linux-ppc64le
contract-plugin-ccrv-darwin-amd64
contract-plugin-ccrv-darwin-arm64
contract-plugin-ccrv-windows-amd64.exe
contract-plugin-ccrv_1.2.0_checksums.txt
```

Each binary responds to `--version` with the tag value (e.g. `v1.2.0`), which `contract-cli plugin list` picks up.

---

## Local Development Workflow

1. **Build your plugin binary:**

   ```sh
   go build -o contract-plugin-<name> .
   ```

2. **Install it locally:**

   ```sh
   contract-cli plugin install ./contract-plugin-<name>
   # or manually:
   cp contract-plugin-<name> ~/.contract-cli/plugins/
   chmod +x ~/.contract-cli/plugins/contract-plugin-<name>
   ```

3. **Verify it appears in the list:**

   ```sh
   contract-cli plugin list
   ```

4. **Invoke it:**

   ```sh
   contract-cli <name> --help
   contract-cli <name> <subcommand> [flags]
   ```

5. **Remove it when done:**

   ```sh
   contract-cli plugin uninstall <name>
   ```

---

## User Quickstart

Install the `ccrv` plugin and use it:

```sh
# Install latest release
contract-cli plugin install github.com/org/contract-plugin-ccrv

# Install a specific version
contract-cli plugin install github.com/org/contract-plugin-ccrv@v1.2.0

# List installed plugins and their versions
contract-cli plugin list

# Use the plugin (all sub-commands pass through)
contract-cli ccrv --help
contract-cli ccrv generate --env prod

# Update to latest
contract-cli plugin update ccrv

# Update all plugins
contract-cli plugin update

# Remove the plugin
contract-cli plugin uninstall ccrv
```
