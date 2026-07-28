package app

import (
	"strings"
	"testing"

	lipgloss "charm.land/lipgloss/v2"

	"github.com/issafalcon/dotfiles-tui/internal/theme"
)

func TestFramePanelClipsOverflow(t *testing.T) {
	tall := strings.Repeat("module card line that should be clipped\n", 40)
	const width, height = 24, 12
	out := framePanel(tall, true, width, height)

	if h := lipgloss.Height(out); h != height {
		t.Fatalf("height=%d want %d\n%s", h, height, out)
	}
	if w := lipgloss.Width(out); w > width {
		t.Fatalf("width=%d exceeds %d", w, width)
	}
}

func TestClipTruncates(t *testing.T) {
	in := strings.Repeat("x\n", 20)
	out := theme.Clip(in, 10, 5)
	if h := lipgloss.Height(out); h != 5 {
		t.Fatalf("height=%d want 5", h)
	}
}
