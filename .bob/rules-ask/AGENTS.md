# AGENTS.md

This file provides guidance to agents when working with code in this repository.

## Non-Obvious Documentation Context

- `lib/` is the internal business-logic layer of the CLI, not a public library. Command long/short descriptions and flag help live as exported string constants in `lib/<name>/<name>.go`, not in `cmd/` (except `cmd/root.go` and `cmd/plugin.go` subcommands, which define some text inline).
- `lib/plugin/` is the only lib package with several source files (`dispatch*.go` split by OS via build tags).
- `samples/` holds real cert/key/contract fixtures used by tests (intentionally committed).
- `build/` is tracked only via `.gitkeep`; contents are gitignored test output.
- `openssl` runtime dependency is documented in `CONTRIBUTING.md`, not the README install section.
- `--os` platform ids (`ccrt`, `ccrv`, `ccco`, `hpvs`): `hpvs` → `hyper-protect-basic`, all others → `contract-basic` encrypted-string prefix.
- Plugin docs: `docs/plugins.md`; design notes in `plugin-support-plan.md`. Plugins are `contract-plugin-<name>` binaries in `~/.contract-cli/plugins`.
