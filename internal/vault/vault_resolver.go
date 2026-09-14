package vault

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

type LayoutStrategy int

const (
	StrategyArtistAlbum LayoutStrategy = iota
	StrategyCompilationGenre
)

type VaultResolver struct {
	BaseDir string
}

func NewVaultResolver(baseDir string) *VaultResolver {
	return &VaultResolver{BaseDir: baseDir}
}

// ResolveDestination calculates target paths based on track metadata and chosen strategy
func (vr *VaultResolver) ResolveDestination(artist, album, genre, filename string, year int, strategy LayoutStrategy) (string, error) {
	var relativeDir string

	switch strategy {
	case StrategyArtistAlbum:
		artistDir := SanitizeInternationalTitle(artist)
		if artistDir == "unknown_track" {
			artistDir = "Unknown_Artist"
		}
		albumDir := SanitizeInternationalTitle(album)
		if albumDir == "unknown_track" {
			albumDir = "Unknown_Album"
		}
		if year > 0 {
			albumDir = fmt.Sprintf("%s (%d)", albumDir, year)
		}
		relativeDir = filepath.Join("Artists", artistDir, albumDir)

	case StrategyCompilationGenre:
		genreDir := SanitizeInternationalTitle(genre)
		if genreDir == "unknown_track" {
			genreDir = "Unsorted"
		}
		relativeDir = filepath.Join("Compilations_by_Genre", genreDir)

	default:
		return "", fmt.Errorf("vault_resolver: unknown layout strategy %d", strategy)
	}

	cleanFile := filepath.Base(filepath.Clean(filename))
	return filepath.Join(vr.BaseDir, relativeDir, cleanFile), nil
}

// ResolvePlaylistPath returns the absolute target path for dynamic M3U8 playlists
func (vr *VaultResolver) ResolvePlaylistPath(playlistName string) string {
	cleanName := SanitizeInternationalTitle(playlistName)
	if !strings.HasSuffix(cleanName, ".m3u8") {
		cleanName += ".m3u8"
	}
	return filepath.Join(vr.BaseDir, "Playlists", cleanName)
}

func SanitizeInternationalTitle(rawTitle string) string {
	normalized := norm.NFC.String(rawTitle)
	var sb strings.Builder
	for _, r := range normalized {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			sb.WriteRune(r)
		} else {
			switch r {
			case '/', '\\', ':', '*', '?', '"', '<', '>', '|', '_':
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
