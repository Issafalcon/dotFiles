package utils

import "testing"

func TestModuleSatisfiedTrueIsTrackingOnly(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if ModuleSatisfied("ai", "true") {
		t.Fatal("check_command true must not imply satisfied without tracking")
	}
}

func TestModuleSatisfiedCheckCommand(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if !ModuleSatisfied("shell", "command -v sh") {
		t.Fatal("expected sh to satisfy command -v sh")
	}
	if ModuleSatisfied("missing", "command -v definitely-not-a-real-bin-xyz") {
		t.Fatal("missing binary should not satisfy")
	}
}
