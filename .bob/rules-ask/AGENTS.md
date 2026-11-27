# AGENTS.md

This file provides guidance to agents when working with code in this repository.

## Non-Obvious Documentation Context

- `lib/` is the internal business-logic layer of the CLI, not a public library. Command long/short descriptions and flag help live as exported string constants in `lib/<name>/<name>.go`, not in `cmd/` (except `cmd/root.go` and `cmd/plugin.go` subcommands, which define their text inline).
- `lib/plugin/` is the only lib package with several source files (`plugin.go`, `github.go` for download/checksum, `dispatch*.go` split by OS via build tags).
- `openssl` runtime dependency is documented in `CONTRIBUTING.md`, not the README install section.
- Plugin docs: `docs/plugins.md` (author guide) and the Plugins section of `docs/README.md` (user commands). `plugin-support-plan.md` is an untracked design note and is stale (it still references old `common.*` plugin helpers that now live in `lib/plugin/`).
- Plugin management commands take the plugin's short name (`ccrv`), not `contract-plugin-ccrv`; `plugin list` shows short names too.
