package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

var (
	dimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	onStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	offStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

// boxChrome is the horizontal space borderStyle adds around its content:
// left margin (2) + two border columns (2) + padding on both sides (2).
const boxChrome = 6

// defaultWidth is used until bubbletea reports the real terminal size.
const defaultWidth = 100

// contentWidth is how many columns are left inside the bordered box for a terminal this wide.
func contentWidth(termWidth int) int {
	if termWidth <= 0 {
		termWidth = defaultWidth
	}
	if w := termWidth - boxChrome; w > 20 {
		return w
	}
	return 20
}

// pad pads s with spaces to width display columns (no-op if already wider).
// Pad before styling: fmt's %-*s counts ANSI escape bytes and wide runes wrong.
func pad(s string, width int) string {
	if gap := width - lipgloss.Width(s); gap > 0 {
		return s + strings.Repeat(" ", gap)
	}
	return s
}

// truncate shortens s to at most width display columns, ending in "…" when cut.
func truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	return ansi.Truncate(s, width, "…")
}

// labelWidth is the widest label, for aligning "Label: value" forms.
func labelWidth(labels []string) int {
	w := 0
	for _, l := range labels {
		if lw := lipgloss.Width(l); lw > w {
			w = lw
		}
	}
	return w
}

// visibleWindow returns the [start, end) slice of n rows to show in height rows,
// keeping cursor in view and the window as stable as possible around offset.
func visibleWindow(n, cursor, offset, height int) (start, end int) {
	if height <= 0 || n <= height {
		return 0, n
	}
	start = offset
	if cursor < start {
		start = cursor
	}
	if cursor >= start+height {
		start = cursor - height + 1
	}
	if start > n-height {
		start = n - height
	}
	if start < 0 {
		start = 0
	}
	return start, start + height
}
