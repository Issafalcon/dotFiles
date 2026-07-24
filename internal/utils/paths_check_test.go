package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestModulesLayout(t *testing.T) {
	// Prefer cwd-based detection: run from tui/ in this repo.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// If we're in .../tui/internal/utils, walk up for modules/
	root := wd
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(filepath.Join(root, "modules")); err == nil {
			break
		}
		root = filepath.Dir(root)
	}
	os.Setenv("DOTFILES_DIR", root)
	defer os.Unsetenv("DOTFILES_DIR")

	if GetModulesDir() != filepath.Join(root, "modules") {
		t.Fatalf("GetModulesDir=%s", GetModulesDir())
	}
	if !ModuleScriptExists("nvim", "install.sh") {
		t.Fatal("expected modules/nvim/install.sh")
	}
	if !ModuleScriptExists("zsh", "uninstall.sh") {
		t.Fatal("expected modules/zsh/uninstall.sh")
	}
}
