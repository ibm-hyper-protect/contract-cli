# AGENTS.md

This file provides guidance to agents when working with code in this repository.

## Non-Obvious Architectural Constraints

- The `cmd/` + `lib/` two-layer split is mandatory: no business logic in `cmd/`, no flag wiring in `lib/`. `common/` is the only shared package.
- `lib/<name>` favors exported constants (flag names/descriptions) over config structs.
- The CLI is an adapter: validate flags → prepare inputs → call `contract-go/v2` (`contract.*`) → write output.
- File-input flags support stdin via `--in -` (`common.ReadDataFromStdin` + `common.ValidateStdinInput`); new input commands should follow this.
- Omitting `--priv` implies an `openssl` runtime dependency (key generation).
- No global state except `rootCmd` and version vars; tests drive `rootCmd.SetArgs(...)` + `Execute()`.
- `CGO_ENABLED=0` is required for cross-compilation — no C dependencies.
- Plugin dispatch is hooked into `rootCmd.Args` (not a Cobra subcommand): unknown first arg → look up `~/.contract-cli/plugins/contract-plugin-<arg>` → `syscall.Exec` (Unix) / child process (Windows) with `CONTRACT_CLI_PLUGIN=1`. Built-in commands always shadow plugins.
- GitHub plugin install expects release assets named `contract-plugin-<name>-<os>-<arch>[.exe]`; local-path installs are skipped by `plugin update`.
