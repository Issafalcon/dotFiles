package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBrewAwareEnvSetsSystemCABundle(t *testing.T) {
	if _, err := os.Stat(systemCABundle); err != nil {
		t.Skip("system CA bundle not present")
	}
	t.Setenv("SSL_CERT_FILE", "")
	t.Setenv("CURL_CA_BUNDLE", "")
	// Clear via environ rebuild: Setenv to empty still leaves key present.
	// Unset so ensureSystemCABundle can fill them.
	_ = os.Unsetenv("SSL_CERT_FILE")
	_ = os.Unsetenv("CURL_CA_BUNDLE")

	env := brewAwareEnv()
	got := map[string]string{}
	for _, kv := range env {
		if i := strings.IndexByte(kv, '='); i > 0 {
			got[kv[:i]] = kv[i+1:]
		}
	}
	if got["SSL_CERT_FILE"] != systemCABundle {
		t.Fatalf("SSL_CERT_FILE=%q, want %q", got["SSL_CERT_FILE"], systemCABundle)
	}
	if got["CURL_CA_BUNDLE"] != systemCABundle {
		t.Fatalf("CURL_CA_BUNDLE=%q, want %q", got["CURL_CA_BUNDLE"], systemCABundle)
	}
}

func TestBrewAwareEnvPreservesExistingCABundle(t *testing.T) {
	custom := "/tmp/custom-ca.pem"
	t.Setenv("SSL_CERT_FILE", custom)
	t.Setenv("CURL_CA_BUNDLE", custom)

	env := brewAwareEnv()
	got := map[string]string{}
	for _, kv := range env {
		if i := strings.IndexByte(kv, '='); i > 0 {
			got[kv[:i]] = kv[i+1:]
		}
	}
	if got["SSL_CERT_FILE"] != custom {
		t.Fatalf("SSL_CERT_FILE overridden: %q", got["SSL_CERT_FILE"])
	}
	if got["CURL_CA_BUNDLE"] != custom {
		t.Fatalf("CURL_CA_BUNDLE overridden: %q", got["CURL_CA_BUNDLE"])
	}
}

func TestCleanOldInstallLogsKeepsTwoDays(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 7, 28, 15, 0, 0, 0, time.Local)
	for _, name := range []string{
		"install-2026-07-28.log", // today — keep
		"install-2026-07-27.log", // yesterday — keep
		"install-2026-07-26.log", // 2 days ago — remove
		"install-2026-07-20.log", // older — remove
		"other.txt",              // ignore
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := CleanOldInstallLogs(dir, now); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"install-2026-07-28.log", "install-2026-07-27.log", "other.txt"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("expected keep %s: %v", name, err)
		}
	}
	for _, name := range []string{"install-2026-07-26.log", "install-2026-07-20.log"} {
		if _, err := os.Stat(filepath.Join(dir, name)); !os.IsNotExist(err) {
			t.Fatalf("expected remove %s", name)
		}
	}
}

func TestOpenInstallLogAppends(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	log, err := OpenInstallLog("cursor", "install")
	if err != nil {
		t.Fatal(err)
	}
	log.Line("hello", false)
	log.Line("fail", true)
	log.Close()

	dir, err := InstallLogDir()
	if err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 1 {
		t.Fatalf("expected 1 log file, got %v err=%v", entries, err)
	}
	data, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if !strings.Contains(s, "cursor") || !strings.Contains(s, "  hello") || !strings.Contains(s, "! fail") {
		t.Fatalf("unexpected log contents:\n%s", s)
	}
}
