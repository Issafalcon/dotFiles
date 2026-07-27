package app

import "testing"

func TestInstallProgress(t *testing.T) {
	plan := []string{"homebrew", "node", "nvim"}
	// After starting homebrew, remaining is [node, nvim]
	step, total := installProgress(plan, []string{"node", "nvim"}, "homebrew")
	if step != 1 || total != 3 {
		t.Fatalf("first: got %d/%d want 1/3", step, total)
	}
	// After homebrew done, starting node, remaining [nvim]
	step, total = installProgress(plan, []string{"nvim"}, "node")
	if step != 2 || total != 3 {
		t.Fatalf("second: got %d/%d want 2/3", step, total)
	}
	step, total = installProgress(plan, nil, "nvim")
	if step != 3 || total != 3 {
		t.Fatalf("last: got %d/%d want 3/3", step, total)
	}
}
