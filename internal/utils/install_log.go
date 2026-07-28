package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	installLogPrefix     = "install-"
	installLogSuffix     = ".log"
	installLogDateLayout = "2006-01-02"
	// Keep today + yesterday; older dated log files are removed.
	installLogKeepDays = 2
)

// InstallLogDir returns ~/.local/state/dotfiles-tui/logs (or $XDG_STATE_HOME/…).
func InstallLogDir() (string, error) {
	state := os.Getenv("XDG_STATE_HOME")
	if state == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		state = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(state, "dotfiles-tui", "logs"), nil
}

// InstallLog appends install/uninstall script output to a dated log file.
type InstallLog struct {
	f *os.File
}

// OpenInstallLog creates/appends today's log, cleans older logs, and writes a
// session header. Callers should defer Close(). Returns nil, nil if logging
// cannot be set up (installs still proceed).
func OpenInstallLog(moduleName, action string) (*InstallLog, error) {
	dir, err := InstallLogDir()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	_ = CleanOldInstallLogs(dir, time.Now())

	path := filepath.Join(dir, installLogPrefix+time.Now().Format(installLogDateLayout)+installLogSuffix)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	_, _ = fmt.Fprintf(f, "\n===== %s %s %s =====\n",
		time.Now().Format(time.RFC3339), action, moduleName)
	return &InstallLog{f: f}, nil
}

// Line writes one output line. isStderr marks stderr with a leading "! ".
func (l *InstallLog) Line(line string, isStderr bool) {
	if l == nil || l.f == nil {
		return
	}
	prefix := "  "
	if isStderr {
		prefix = "! "
	}
	_, _ = fmt.Fprintf(l.f, "%s%s\n", prefix, line)
}

// Close flushes and closes the log file.
func (l *InstallLog) Close() {
	if l == nil || l.f == nil {
		return
	}
	_ = l.f.Close()
	l.f = nil
}

// CleanOldInstallLogs deletes install-YYYY-MM-DD.log files older than
// installLogKeepDays calendar days relative to now.
func CleanOldInstallLogs(dir string, now time.Time) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	// Keep the last installLogKeepDays calendar days (today inclusive).
	y, m, d := now.Date()
	oldestKeep := time.Date(y, m, d, 0, 0, 0, 0, now.Location()).AddDate(0, 0, -(installLogKeepDays - 1))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasPrefix(name, installLogPrefix) || !strings.HasSuffix(name, installLogSuffix) {
			continue
		}
		dateStr := strings.TrimSuffix(strings.TrimPrefix(name, installLogPrefix), installLogSuffix)
		t, err := time.ParseInLocation(installLogDateLayout, dateStr, now.Location())
		if err != nil {
			continue
		}
		if t.Before(oldestKeep) {
			_ = os.Remove(filepath.Join(dir, name))
		}
	}
	return nil
}
