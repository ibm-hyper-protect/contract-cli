package main

import "github.com/ibm-hyper-protect/contract-cli/cmd"

var (
	version   = "dev"
	osName    = "unknown"
	osArch    = "unknown"
	buildDate = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, osName, osArch, buildDate)
	cmd.Execute()
}
