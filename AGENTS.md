# AGENTS.md

This file provides guidance to agents when working with code in this repository.

Go CLI (`contract-cli`, Cobra) for IBM Confidential Computing contracts. Crypto is done by `contract-go/v2`; some paths shell out to **`openssl`** (must be in `PATH`) via `common.ExecCommand`. `CGO_ENABLED=0` — never add C dependencies.

## Commands

```bash
make test                # go test ./... -v   (no linter configured; use `make fmt`)
make build               # binary at build/contract-cli
go test ./cmd -v -run TestEncryptCmdText                                 # single cmd test
go test ./lib/encrypt -v -run TestGenerateSignedEncryptContract_Success  # single lib test
```

`build/` must exist (tracked via `.gitkeep`, contents gitignored). Tests write outputs there from `lib/<pkg>/` via `../../build/`.

## Architecture: strict two-layer pattern

- `cmd/<name>.go` — Cobra wiring only: flags, `ValidateInput()`, `Output()`; `Run` uses `log.Fatal(err)` (not `fmt` + `os.Exit`).
- `lib/<name>/<name>.go` — all logic plus **exported constants** for flag names, descriptions and `ParameterName`. Never hardcode flag strings in `cmd/` (or tests) — import the constants.
- Shared code goes only in `common/`.

## common/ gotchas

- `SetMandatoryFlagError(cmd, err)` prints usage and calls `os.Exit(1)` — it never returns. Used for missing mandatory flags in `ValidateInput` instead of `return err`.
- `ValidateStdinInput(cmd, in)` is side-effectful (exits): `--in -` without piped stdin, or piped stdin with `--in <file>`, are both fatal.
- `GetPrivateKey("")` generates a fresh RSA-4096 key through `openssl genrsa` — `--priv` is optional, but omitting it adds an openssl dependency.
- `WriteDataToFile` writes `0600`; use it instead of `os.WriteFile`/`os.Create` for outputs.
- `ReadDataFromFile`/`GetDataFromFile` return `(string, error)`; `GetDataFromFile("")` returns `""` for optional file flags.
- YAML: use `sigs.k8s.io/yaml` (Kubernetes-compatible), not `gopkg.in/yaml.v3` alone.
- `--os` platforms: `hpvs` → `hyper-protect-basic` prefix; `ccrt`/`ccrv`/`ccco`/empty → `contract-basic`.

## New command checklist

1. `lib/<name>/<name>.go`, then `cmd/<name>.go`; in `init()` call both `common.SetCustomHelpTemplate(cmd, map[string]bool{"flag": true})` (required flags) and `common.SetCustomErrorTemplate(cmd)`.
2. Tests: `cmd/<name>_test.go` (via `rootCmd.SetArgs` + `Execute`) and `lib/<name>/<name>_test.go`; names `TestFunc_Scenario_Expected`.
3. Document in `docs/README.md`.

## Plugins (`cmd/plugin.go`, `lib/plugin/`, `docs/plugins.md`)

- Binaries `contract-plugin-<name>` live in `~/.contract-cli/plugins` (`plugin.PluginDir()` in `lib/plugin`, auto-created); sources tracked in `.sources.json`.
- Dispatch happens in `rootCmd.Args` in `cmd/root.go`: unknown first arg → `plugin.FindPlugin` → `plugin.Exec` (`syscall.Exec` on Unix, child process on Windows; sets `CONTRACT_CLI_PLUGIN=1`). Built-in commands always win over a same-named plugin.
- Group command `pluginCmd` must NOT get `SetCustomHelpTemplate` (it would hide the subcommand list); leaf plugin commands pass an empty required-flags map.
- Tests isolate by `t.Setenv("HOME", tmp)` (`setPluginDir` in `cmd/plugin_test.go`); `plugin.Exec` can't be tested in-process.

## CI-enforced conventions

- Branch: `<Feature|Fix|Docs|Refactor|Performance|Test|Chore|CI>/<name>`; no force pushes on PRs.
- Conventional Commits (`feat:`, `fix:`, …); all commits must be GPG/SSH signed.
