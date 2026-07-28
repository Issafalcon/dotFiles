// Package config loads and saves user configuration for the DotFiles TUI.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	appConfigDirName  = "dotfiles-tui"
	appConfigFileName = "config.yaml"
	envModulesDir     = "DOTFILES_MODULES_DIR"
	envDotfilesDir    = "DOTFILES_DIR"
)

// Config is persisted under ~/.config/dotfiles-tui/config.yaml.
type Config struct {
	ModulesDir string `yaml:"modules_dir"`
}

// Dir returns ~/.config/dotfiles-tui (creating it if needed when write=true).
func Dir(create bool) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", appConfigDirName)
	if create {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return "", err
		}
	}
	return dir, nil
}

// Path returns the config file path.
func Path() (string, error) {
	dir, err := Dir(false)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, appConfigFileName), nil
}

// Load reads the config file. Missing file returns empty Config and nil error.
func Load() (Config, error) {
	path, err := Path()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, err
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}
	return c, nil
}

// Save writes the config file (creates the config directory).
func Save(c Config) error {
	dir, err := Dir(true)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, appConfigFileName), data, 0o644)
}

// ResolveModulesDir returns the modules directory using:
//  1. DOTFILES_MODULES_DIR
//  2. config modules_dir
//  3. DOTFILES_DIR/modules or auto-detected repo modules/
//  4. "" if none configured (caller should run first-run setup)
func ResolveModulesDir() (string, error) {
	if env := os.Getenv(envModulesDir); env != "" {
		return filepath.Clean(env), nil
	}

	cfg, err := Load()
	if err != nil {
		return "", err
	}
	if cfg.ModulesDir != "" {
		return filepath.Clean(cfg.ModulesDir), nil
	}

	if env := os.Getenv(envDotfilesDir); env != "" {
		return filepath.Join(filepath.Clean(env), "modules"), nil
	}

	if root := findRepoWithModules(); root != "" {
		return filepath.Join(root, "modules"), nil
	}

	return "", nil
}

// findRepoWithModules walks up from cwd / executable looking for modules/.
func findRepoWithModules() string {
	candidates := []string{}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, cwd)
	}
	if execPath, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Dir(execPath))
	}
	for _, start := range candidates {
		dir := start
		for i := 0; i < 8; i++ {
			mod := filepath.Join(dir, "modules")
			if st, err := os.Stat(mod); err == nil && st.IsDir() {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return ""
}

// SetModulesDir saves modules_dir to the config file.
func SetModulesDir(dir string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	cfg.ModulesDir = filepath.Clean(dir)
	return Save(cfg)
}
