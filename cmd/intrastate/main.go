// Command intrastate is the CLI entry point. It delegates immediately to
// internal/cli so the binary stays a thin shim and all wiring (flags,
// subcommands, exit-code mapping) lives in one testable package.
package main

import "github.com/newcoinc/intrastate/internal/cli"

func main() {
	cli.Execute()
}
