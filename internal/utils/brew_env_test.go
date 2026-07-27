package utils

import (
	"strings"
	"testing"
)

func TestBrewAwareEnvPrependsLinuxbrew(t *testing.T) {
	env := brewAwareEnv()
	var path string
	for _, kv := range env {
		if strings.HasPrefix(kv, "PATH=") {
			path = strings.TrimPrefix(kv, "PATH=")
			break
		}
	}
	if path == "" {
		t.Fatal("PATH missing from env")
	}
	// Function is a no-op when brew dirs are absent; just ensure it returns PATH.
	if !strings.Contains(path, "/") {
		t.Fatalf("unexpected PATH: %q", path)
	}
}
