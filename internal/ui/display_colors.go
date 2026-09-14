package ui

import "fmt"

// ANSI Color Codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	// Foreground Colors
	FgCyan    = "\033[36m"
	FgGreen   = "\033[32m"
	FgYellow  = "\033[33m"
	FgMagenta = "\033[35m"
	FgRed     = "\033[31m"
	FgBlue    = "\033[34m"
	FgWhite   = "\033[37m"
	FgGray    = "\033[90m"

	// Background Colors
	BgDarkGray = "\033[100m"
	BgBlue     = "\033[44m"
)

func Colorize(text, color string) string {
	return color + text + Reset
}

func HexToRGB(hex string) (int, int, int) {
	var r, g, b int
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return r, g, b
}

func TrueColor(text, hex string) string {
	r, g, b := HexToRGB(hex)
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, text)
}
