package utils

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

// ModuleSatisfied reports whether a module is already available and does not
// need to be installed again.
//
// Order of checks:
//  1. Name is listed in ~/.dotFileModules (explicitly installed via the TUI)
//  2. check_command succeeds when run under sh -c (binary/path already present)
//
// Literal check commands "true" / "false" are ignored for (2) — those modules
// are tracking-file only (e.g. skill packs with no single binary).
func ModuleSatisfied(name, checkCommand string) bool {
	installed, err := GetInstalledModules()
	if err == nil {
		for _, m := range installed {
			if m == name {
				return true
			}
		}
	}

	cmd := strings.TrimSpace(checkCommand)
	if cmd == "" || cmd == "true" || cmd == "false" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "sh", "-c", cmd).Run() == nil
}
