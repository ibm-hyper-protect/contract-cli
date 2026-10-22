# AGENTS.md

This file provides guidance to agents when working with code in this repository.

## Non-Obvious Coding Rules

- `common.SetMandatoryFlagError(cmd, err)` prints error + usage then `os.Exit(1)` — it does NOT return; don't write code after it expecting to continue. Missing mandatory flags in `ValidateInput` use it rather than `return err`.
- `common.ValidateStdinInput` is side-effectful (exits) — call it only after reading the `--in` flag.
- `common.GetPrivateKey("")` shells out to `openssl genrsa` (RSA-4096) — `--priv` is optional on most commands.
- `common.WriteDataToFile` writes `0600`; don't use `os.WriteFile`/`os.Create` for outputs.
- `ReadDataFromFile`/`GetDataFromFile` return `(string, error)`.
- Cobra `Run` in `cmd/` uses `log.Fatal(err)`.
- Flag names come from exported constants in `lib/<name>/` — import them in `cmd/` and tests, never literal strings.
- Use `sigs.k8s.io/yaml` for YAML marshaling.
- Don't call `common.SetCustomHelpTemplate` on parent/group commands (e.g. `pluginCmd`) — it suppresses the subcommand listing. Leaf commands without named flags pass `map[string]bool{}`.
- Plugin tests must isolate `HOME` via `t.Setenv` (see `setPluginDir` in `cmd/plugin_test.go`); plugin dispatch lives in `rootCmd.Args` (`cmd/root.go`), and built-ins take priority over plugins.
- Test files write to `../../build/` from `lib/<pkg>/`; `build/` must exist.
