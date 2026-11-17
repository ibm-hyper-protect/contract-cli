// Copyright (c) 2026 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// HandshakeEnvVar is the environment variable injected by contract-cli on every plugin dispatch.
	// Plugin authors must check for this variable at startup using CheckHandshake().
	HandshakeEnvVar = "CONTRACT_CLI_PLUGIN"

	// HandshakeEnvVal is the expected value of HandshakeEnvVar.
	HandshakeEnvVal = "1"
)

// CheckHandshake returns an error if the current process was NOT started by contract-cli.
// Plugin authors must call this at the very top of main() and exit if it returns an error:
//
//	if err := plugin.CheckHandshake(); err != nil {
//	    fmt.Fprintln(os.Stderr, "Error:", err)
//	    os.Exit(1)
//	}
func CheckHandshake() error {
	if os.Getenv(HandshakeEnvVar) != HandshakeEnvVal {
		return fmt.Errorf("this plugin must be run through contract-cli")
	}
	return nil
}

// FindPlugin looks up a plugin binary for the given sub-command name.
// It returns the full path to the binary and true when found and executable,
// or ("", false) when no matching plugin exists.
func FindPlugin(name, pluginDir string) (string, bool) {
	binaryName := PluginPrefix + name
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	path := filepath.Join(pluginDir, binaryName)
	info, err := os.Stat(path)
	if err != nil {
		return "", false
	}
	// On Unix verify the executable bit; on Windows any regular file is runnable.
	if runtime.GOOS != "windows" && (info.Mode()&0111 == 0) {
		return "", false
	}
	return path, true
}

// Exec replaces the current process (Linux/macOS) or spawns a child process (Windows)
// with the plugin binary, passing args and inheriting all standard I/O.
// CONTRACT_CLI_PLUGIN=1 is always injected into the child environment.
//
// On Linux/macOS this function never returns on success (syscall.Exec replaces the process).
// On Windows it returns after the child exits, calling os.Exit with the child's exit code.
func Exec(binaryPath string, args []string) error {
	env := append(os.Environ(), HandshakeEnvVar+"="+HandshakeEnvVal)

	if runtime.GOOS == "windows" {
		return execWindows(binaryPath, args, env)
	}
	return execUnix(binaryPath, args, env)
}

// execUnix uses syscall.Exec to replace the current process on Linux/macOS.
// This provides clean exit-code propagation.
func execUnix(binaryPath string, args []string, env []string) error {
	argv := append([]string{binaryPath}, args...)
	return execSyscall(binaryPath, argv, env)
}

// execWindows runs the plugin as a child process and forwards its exit code.
func execWindows(binaryPath string, args []string, env []string) error {
	cmd := exec.Command(binaryPath, args...) // #nosec G204
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return err
	}
	os.Exit(0)
	return nil // unreachable
}

// pluginNameFromArg extracts the plugin short-name from an argument string.
// e.g. "ccrv" → "ccrv" (unchanged).
func pluginNameFromArg(arg string) string {
	return strings.TrimPrefix(arg, PluginPrefix)
}
