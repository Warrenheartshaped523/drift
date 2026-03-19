package main

import (
	"os"

	"github.com/phlx0/drift/cmd/drift"
)

// Build information injected by goreleaser via -ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	drift.SetBuildInfo(version, commit, date)
	if err := drift.Execute(); err != nil {
		os.Exit(1)
	}
}
