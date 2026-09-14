package gate

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

var (
	// Regex to target characters unsafe for cross-platform filesystems (Windows/Linux/macOS)
	illegalFileNameChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	multiSpaceRegex      = regexp.MustCompile(`\s+`)
)

type HubSanitizer struct {
	maxFileNameLength int
}

func NewHubSanitizer() *HubSanitizer {
	return &HubSanitizer{
		maxFileNameLength: 180, // Safe threshold under standard OS path limits (255 chars)
	}
}

// SanitizeFileName cleans artist/title metadata into a deterministic, filesystem-safe string
func (hs *HubSanitizer) SanitizeFileName(raw string) string {
	cleaned := strings.TrimSpace(raw)
	if cleaned == "" {
		return "unnamed_track"
	}
	// 1. Remove non-printable control characters
	cleaned = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, cleaned)
	// 2. Neutralize reserved characters
	cleaned = illegalFileNameChars.ReplaceAllString(cleaned, "_")
	// 3. Collapse multiple whitespace blocks
	cleaned = multiSpaceRegex.ReplaceAllString(cleaned, " ")
	// 4. Truncate long strings to avoid OS path overflow
	if len(cleaned) > hs.maxFileNameLength {
		cleaned = cleaned[:hs.maxFileNameLength]
	}
	return strings.Trim(cleaned, ". ")
}

// SecureDestinationPath ensures output target directories remain contained within the base vault path
func (hs *HubSanitizer) SecureDestinationPath(baseDir, relativePath string) (string, error) {
	cleanBase := filepath.Clean(baseDir)
	targetPath := filepath.Clean(filepath.Join(cleanBase, relativePath))
	// Enforce strict boundary check to neutralize directory traversal attacks
	if !strings.HasPrefix(targetPath, cleanBase+string(filepath.Separator)) && targetPath != cleanBase {
		return "", fmt.Errorf("path security escape attempt blocked: %s", relativePath)
	}
	return targetPath, nil
}
