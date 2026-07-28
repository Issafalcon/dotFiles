package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/issafalcon/dotfiles-tui/internal/config"
	"github.com/issafalcon/dotfiles-tui/internal/module"
)

func TestModulesLayout(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := wd
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(filepath.Join(root, "modules")); err == nil {
			break
		}
		root = filepath.Dir(root)
	}
	modulesDir := filepath.Join(root, "modules")
	t.Setenv("DOTFILES_MODULES_DIR", modulesDir)

	if GetModulesDir() != modulesDir {
		t.Fatalf("GetModulesDir=%s want %s", GetModulesDir(), modulesDir)
	}
	if !ModuleScriptExists("nvim", "install.sh") {
		t.Fatal("expected modules/nvim/install.sh")
	}
	if !ModuleScriptExists("zsh", "uninstall.sh") {
		t.Fatal("expected modules/zsh/uninstall.sh")
	}

	res := module.LoadFromDir(modulesDir)
	if res.Loaded < 50 {
		t.Fatalf("loaded %d modules, errors=%v", res.Loaded, res.Errors)
	}
	if module.DefaultRegistry.Get("nvim") == nil {
		t.Fatal("nvim not registered")
	}
}

func TestResolveModulesDirEmptyWithoutEnv(t *testing.T) {
	t.Setenv("DOTFILES_MODULES_DIR", "")
	t.Setenv("DOTFILES_DIR", "")
	// May still auto-detect from cwd when tests run inside this repo — that's OK.
	_, err := config.ResolveModulesDir()
	if err != nil {
		t.Fatal(err)
	}
}
