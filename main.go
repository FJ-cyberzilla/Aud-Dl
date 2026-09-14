package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"unicode"

	"audio-command-center/internal/cache"
	"audio-command-center/internal/cli"
	"audio-command-center/internal/ui"
	"audio-command-center/internal/vault"

	"golang.org/x/text/unicode/norm"
)

func SanitizeInternationalTitle(rawTitle string) string {
	normalized := norm.NFC.String(rawTitle)
	var sb strings.Builder
	for _, r := range normalized {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			sb.WriteRune(r)
		} else {
			switch r {
			case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
				sb.WriteRune('_')
			}
		}
	}
	cleaned := strings.TrimSpace(sb.String())
	for strings.Contains(cleaned, "__") {
		cleaned = strings.ReplaceAll(cleaned, "__", "_")
	}
	if cleaned == "" {
		return "unknown_track"
	}
	return cleaned
}

func handleSelection(selected int) {
	fmt.Printf("Selected item: %d\n", selected)
}

func main() {
	// Initialize subsystems
	cacheMgr := cache.NewCacheManager()
	vaultMgr := vault.NewVault()
	style := cli.NewInteractiveStyle()

	// Create advanced menu system
	menu := cli.NewTelemetryMenu(style, cacheMgr, log.Default())
	smartMenu := cli.NewSmartMenu(menu.MenuTemplate, vaultMgr)

	// Render with all features
	ctx := context.Background()
	config := cli.MenuConfig{
		Title:      "🎵 AUDIO COMMAND CENTER",
		Subtitle:   "v3.0.0 - Enterprise Edition",
		ShowHelp:   true,
		Paginate:   true,
		PageSize:   10,
		WrapAround: true,
	}

	// Main loop
	selected := 0
	for {
		output := smartMenu.RenderSmart(ctx, cli.MainMenuItems, selected, config)
		fmt.Print(output)

		// Handle input
		key := ui.ReadKey()
		switch key {
		case ui.KeyUp:
			selected--
			if selected < 0 {
				selected = len(cli.MainMenuItems) - 1
			}
		case ui.KeyDown:
			selected++
			if selected >= len(cli.MainMenuItems) {
				selected = 0
			}
		case ui.KeyEnter:
			// Execute selection
			handleSelection(selected)
		}
	}
}
