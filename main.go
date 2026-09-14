package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
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
	cacheMgr := cache.NewCacheManager(1 * time.Minute)
	_, _ = vault.NewMusicVault("./music_vault")
	style := cli.NewInteractiveStyle(ui.CyberDarkTheme)

	// Create advanced menu system
	_ = cli.NewTelemetryMenu(style, cacheMgr, slog.Default())
	// smartMenu := cli.NewSmartMenu(menu.MenuTemplate, vaultMgr) // TODO: Implement SmartMenu

	// Render with all features
	_ = context.Background()
	_ = cli.MenuConfig{
		Title:      "🎵 AUDIO COMMAND CENTER",
		Subtitle:   "v3.0.0 - Enterprise Edition",
		ShowHelp:   true,
		Paginate:   true,
		PageSize:   10,
		WrapAround: true,
	}

	// Main loop stub
	fmt.Println("Audio Command Center initialized. (Main loop input pending implementation)")
}
