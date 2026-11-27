# AGENTS.md

This file provides guidance to agents when working with code in this repository.

## Non-Obvious Coding Rules

Root `AGENTS.md` covers the `common/` helpers and the plugin layout; these are the extra traps when editing.

- Code after `common.SetMandatoryFlagError(...)` is dead — it calls `os.Exit(1)`.
- Plugin leaf commands print with `fmt` to the process `os.Stdout`, so `cmd` tests capture stdout via `os.Pipe` (`runPluginCmd` in `cmd/plugin_test.go`); help output goes through `cmd.OutOrStdout()`, so help tests need `cmd.SetOut(...)` instead (and must reset the `help` flag — Cobra command state is shared across tests).
- Plugin install/update functions take an `insecure bool`; wire new call sites to `plugin.InsecureFlag`, never a literal `"insecure"`.
- Plugin `--version` output is parsed by `plugin.parseVersion` (token after `version`); keep `PluginVersion` returning the raw first line.
