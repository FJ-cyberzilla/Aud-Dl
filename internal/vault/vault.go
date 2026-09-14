package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Metadata struct {
	Title     string
	Artist    string
	Album     string
	Genre     string
	Year      string
	IsVarious bool // True if track is part of a compilation/various artists album
}

type MusicVault struct {
	BasePath string
}

func NewMusicVault(basePath string) (*MusicVault, error) {
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("invalid vault base path: %w", err)
	}

	// Ensure root vault directory exists
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to initialize vault directory: %w", err)
	}

	return &MusicVault{BasePath: absPath}, nil
}

// StoreTrack moves a processed MP3 into the vault structure based on metadata rules
func (v *MusicVault) StoreTrack(tempFilePath string, meta Metadata) (string, error) {
	targetFolder := v.resolveTargetPath(meta)

	// Ensure destination directory exists
	if err := os.MkdirAll(targetFolder, 0755); err != nil {
		return "", fmt.Errorf("failed to create vault folder structure: %w", err)
	}

	// Format filename: "01 - Track Title.mp3" or "Artist - Track Title.mp3"
	fileName := fmt.Sprintf("%s - %s.mp3", sanitizePathSegment(meta.Artist), sanitizePathSegment(meta.Title))
	if meta.IsVarious {
		fileName = fmt.Sprintf("Various Artists - %s.mp3", sanitizePathSegment(meta.Title))
	}

	finalDestination := filepath.Join(targetFolder, fileName)

	// Atomic Move (Rename) from temp location to permanent Vault
	err := os.Rename(tempFilePath, finalDestination)
	if err != nil {
		// Fallback to Copy & Remove if moving across different mount points/disks
		err = copyAndRemove(tempFilePath, finalDestination)
		if err != nil {
			return "", fmt.Errorf("failed to transfer file to vault: %w", err)
		}
	}

	return finalDestination, nil
}

// Smart Routing Logic: Organizes by Artist/Album OR Genre/Various
func (v *MusicVault) resolveTargetPath(meta Metadata) string {
	if meta.IsVarious || strings.EqualFold(meta.Artist, "Various Artists") {
		genre := sanitizePathSegment(meta.Genre)
		if genre == "" {
			genre = "Unsorted_Genre"
		}
		return filepath.Join(v.BasePath, "Compilations_by_Genre", genre)
	}

	// Standard Artist Routing: Artists -> [Artist Name] -> [Album Name (Year)]
	artist := sanitizePathSegment(meta.Artist)
	if artist == "" {
		artist = "Unknown Artist"
	}

	album := sanitizePathSegment(meta.Album)
	if album == "" {
		album = "Singles & EPs"
	} else if meta.Year != "" {
		album = fmt.Sprintf("%s (%s)", album, sanitizePathSegment(meta.Year))
	}

	return filepath.Join(v.BasePath, "Artists", artist, album)
}

// Sanitize string inputs to prevent invalid OS directory/file paths
var invalidPathChars = regexp.MustCompile(`[\\/:\*\?"<>\|]`)

func sanitizePathSegment(segment string) string {
	trimmed := strings.TrimSpace(segment)
	clean := invalidPathChars.ReplaceAllString(trimmed, "_")
	return strings.Join(strings.Fields(clean), " ")
}

func copyAndRemove(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.WriteFile(dst, input, 0644); err != nil {
		return err
	}
	return os.Remove(src)
}
