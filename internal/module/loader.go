package module

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadResult summarizes a modules-directory scan.
type LoadResult struct {
	Loaded int
	Skipped []string // dirs without module.yaml
	Errors  []string // parse / read failures
}

// LoadFromDir clears the registry and loads every modules/<name>/module.yaml.
// Directory name wins over yaml "name" when they disagree.
func LoadFromDir(modulesDir string) LoadResult {
	DefaultRegistry = NewRegistry()
	var result LoadResult

	if modulesDir == "" {
		return result
	}

	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("read modules dir: %v", err))
		return result
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		yamlPath := filepath.Join(modulesDir, name, "module.yaml")
		data, err := os.ReadFile(yamlPath)
		if err != nil {
			if os.IsNotExist(err) {
				result.Skipped = append(result.Skipped, name)
				continue
			}
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", name, err))
			continue
		}

		var m Module
		if err := yaml.Unmarshal(data, &m); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s/module.yaml: %v", name, err))
			continue
		}
		m.Name = name // directory is the source of truth
		if m.Category == "" {
			m.Category = "Utility"
		}
		DefaultRegistry.Register(&m)
		result.Loaded++
	}
	return result
}
