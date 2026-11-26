# AGENTS.md

This file provides guidance to agents when working with code in this repository.

## Non-Obvious Architectural Constraints

- The `cmd/` + `lib/` two-layer split is mandatory: no business logic in `cmd/`, no flag wiring in `lib/`. `common/` is the only shared package and is for generic helpers only (plugin logic lives in `lib/plugin/`).
- The CLI is an adapter: validate flags → prepare inputs → call `contract-go/v2` (`contract.*`) → write output.
- File-input flags support stdin via `--in -` (`common.ReadDataFromStdin` + `common.ValidateStdinInput`); new input commands should follow this.
- No global state except `rootCmd` and version vars; tests drive `rootCmd.SetArgs(...)` + `Execute()`.
- Plugin dispatch is hooked into `rootCmd.Args` (not a Cobra subcommand): unknown first arg → `~/.contract-cli/plugins/contract-plugin-<arg>` → `syscall.Exec` (Unix) / child process (Windows) with `CONTRACT_CLI_PLUGIN=1`. `registerPluginCommands` also adds pass-through stub commands (flag parsing disabled, shown under "Plugin Commands" in `--help`) that call `plugin.Exec`; built-ins always shadow plugins.
- GitHub plugin install expects release assets named `contract-plugin-<name>-<os>-<arch>[.exe]` plus a SHA-256 checksum file in the same release; it fails closed unless `--insecure`. Install provenance is the `.sources.json` entry: local installs have none and are skipped by `plugin update`.
