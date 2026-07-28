package module

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromDirEmpty(t *testing.T) {
	dir := t.TempDir()
	res := LoadFromDir(dir)
	if res.Loaded != 0 {
		t.Fatalf("loaded=%d", res.Loaded)
	}
	if len(DefaultRegistry.All()) != 0 {
		t.Fatal("registry should be empty")
	}
}

func TestLoadFromDirYAML(t *testing.T) {
	dir := t.TempDir()
	mod := filepath.Join(dir, "demo")
	if err := os.Mkdir(mod, 0o755); err != nil {
		t.Fatal(err)
	}
	yaml := "name: demo\ndescription: Demo tool\ncategory: Utility\nstow_enabled: false\ncheck_command: true\n"
	if err := os.WriteFile(filepath.Join(mod, "module.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	res := LoadFromDir(dir)
	if res.Loaded != 1 {
		t.Fatalf("loaded=%d errors=%v", res.Loaded, res.Errors)
	}
	m := DefaultRegistry.Get("demo")
	if m == nil || m.Description != "Demo tool" {
		t.Fatalf("got %#v", m)
	}
}
