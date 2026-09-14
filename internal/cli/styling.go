package cli

import (
	"audio-command-center/internal/ui"
	"fmt"
)

type ThemeMode int

const (
	ThemeDark ThemeMode = iota
	ThemeNeon
)

type StyleDispatcher struct {
	theme ui.Theme
}

func NewStyleDispatcher(mode ThemeMode) *StyleDispatcher {
	var theme ui.Theme
	switch mode {
	case ThemeNeon:
		theme = ui.MidnightNeonTheme
	default:
		theme = ui.CyberDarkTheme
	}
	return &StyleDispatcher{theme: theme}
}

// Colorize applies ANSI color codes based on theme colors.
func (sd *StyleDispatcher) Colorize(text string, style string) string {
	var colorCode string
	switch style {
	case "primary":
		colorCode = "\033[36m" // Cyan-ish
	case "success":
		colorCode = "\033[32m" // Green
	case "warning":
		colorCode = "\033[33m" // Yellow
	case "danger":
		colorCode = "\033[31m" // Red
	case "bold":
		colorCode = "\033[1m"
	default:
		colorCode = "\033[0m"
	}
	return fmt.Sprintf("%s%s\033[0m", colorCode, text)
}

func (sd *StyleDispatcher) RenderBadge(text string, active bool) string {
	if active {
		return sd.Colorize(fmt.Sprintf("[%s]", text), "success")
	}
	return sd.Colorize(fmt.Sprintf("[%s]", text), "warning")
}

type JunctionTemplate struct {
}

func NewJunctionTemplate() *JunctionTemplate {
	return &JunctionTemplate{}
}
