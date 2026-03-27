package main

import (
	"os"
	"runtime/debug"

	"github.com/phlx0/drift/cmd/drift"
)

// Build information injected by goreleaser via -ldflags.
// Falls back to embedded module info when installed via go install.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok {
			if v := info.Main.Version; v != "" && v != "(devel)" {
				version = v
			}
			for _, s := range info.Settings {
				switch s.Key {
				case "vcs.revision":
					if len(s.Value) >= 7 {
						commit = s.Value[:7]
					}
				case "vcs.time":
					date = s.Value
				}
			}
		}
	}
	drift.SetBuildInfo(version, commit, date)
	if err := drift.Execute(); err != nil {
		os.Exit(1)
	}
}
