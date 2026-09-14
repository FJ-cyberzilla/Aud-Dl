package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bogem/id3v2/v2"
)

type TrackInfo struct {
	Path     string
	Title    string
	Artist   string
	Duration int // Duration in seconds
}

type PlaylistGenerator struct {
	VaultBasePath string
	PlaylistsDir  string
}

func NewPlaylistGenerator(vaultBasePath string) (*PlaylistGenerator, error) {
	playlistsDir := filepath.Join(vaultBasePath, "Playlists")
	if err := os.MkdirAll(playlistsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create playlists folder: %w", err)
	}

	return &PlaylistGenerator{
		VaultBasePath: vaultBasePath,
		PlaylistsDir:  playlistsDir,
	}, nil
}

// RegenerateAllPlaylists scans the entire vault and rebuilds Genre, Artist, and Full Library playlists
func (pg *PlaylistGenerator) RegenerateAllPlaylists() error {
	tracksByArtist := make(map[string][]TrackInfo)
	tracksByGenre := make(map[string][]TrackInfo)
	var allTracks []TrackInfo

	// Walk through the vault directory
	err := filepath.Walk(pg.VaultBasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(strings.ToLower(info.Name()), ".mp3") {
			return nil
		}

		// Extract track metadata and relative path for playlist portability
		track := pg.extractTrackInfo(path)
		relPath, _ := filepath.Rel(pg.PlaylistsDir, path)
		track.Path = filepath.ToSlash(relPath) // Ensure Unix slashes for cross-platform playability

		allTracks = append(allTracks, track)

		if track.Artist != "" {
			tracksByArtist[track.Artist] = append(tracksByArtist[track.Artist], track)
		}
		if genre := pg.detectGenreFromPath(path); genre != "" {
			tracksByGenre[genre] = append(tracksByGenre[genre], track)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed walking vault: %w", err)
	}

	// 1. Generate Master All-Tracks Playlist
	if err := pg.writeM3U8("Master_Library.m3u8", allTracks); err != nil {
		return err
	}

	// 2. Generate Artist Playlists
	for artist, tracks := range tracksByArtist {
		filename := fmt.Sprintf("Artist_%s.m3u8", sanitizePathSegment(artist))
		if err := pg.writeM3U8(filename, tracks); err != nil {
			return err
		}
	}

	// 3. Generate Genre Playlists
	for genre, tracks := range tracksByGenre {
		filename := fmt.Sprintf("Genre_%s.m3u8", sanitizePathSegment(genre))
		if err := pg.writeM3U8(filename, tracks); err != nil {
			return err
		}
	}

	return nil
}

// writeM3U8 creates an extended UTF-8 formatted M3U playlist file
func (pg *PlaylistGenerator) writeM3U8(filename string, tracks []TrackInfo) error {
	playlistPath := filepath.Join(pg.PlaylistsDir, filename)
	file, err := os.Create(playlistPath)
	if err != nil {
		return fmt.Errorf("failed creating playlist file %s: %w", filename, err)
	}
	defer file.Close()

	// Write M3U8 Header
	file.WriteString("#EXTM3U\n")
	file.WriteString("#PLAYLIST:" + strings.TrimSuffix(filename, ".m3u8") + "\n\n")

	for _, track := range tracks {
		// Write Extended M3U Metadata: #EXTINF:<duration>,<artist> - <title>
		extInf := fmt.Sprintf("#EXTINF:%d,%s - %s\n", track.Duration, track.Artist, track.Title)
		file.WriteString(extInf)
		file.WriteString(track.Path + "\n\n")
	}

	return nil
}

func (pg *PlaylistGenerator) extractTrackInfo(fullPath string) TrackInfo {
	info := TrackInfo{
		Title:    strings.TrimSuffix(filepath.Base(fullPath), ".mp3"),
		Artist:   "Unknown Artist",
		Duration: 180, // Fallback estimated duration
	}

	// Read ID3 metadata if present
	tag, err := id3v2.Open(fullPath, id3v2.Options{Parse: true})
	if err == nil {
		defer tag.Close()
		if title := tag.Title(); title != "" {
			info.Title = title
		}
		if artist := tag.Artist(); artist != "" {
			info.Artist = artist
		}
	}

	return info
}

func (pg *PlaylistGenerator) detectGenreFromPath(fullPath string) string {
	parts := strings.Split(filepath.ToSlash(fullPath), "/")
	for i, part := range parts {
		if part == "Compilations_by_Genre" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// AppendTrack adds a single track to a specified playlist file
func (pg *PlaylistGenerator) AppendTrack(playlistPath, trackPath string) error {
	track := pg.extractTrackInfo(trackPath)
	relPath, _ := filepath.Rel(pg.PlaylistsDir, trackPath)
	track.Path = filepath.ToSlash(relPath)

	f, err := os.OpenFile(playlistPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed opening playlist %s: %w", playlistPath, err)
	}
	defer f.Close()

	// Ensure header exists if new file
	info, _ := f.Stat()
	if info.Size() == 0 {
		f.WriteString("#EXTM3U\n")
		f.WriteString("#PLAYLIST:" + strings.TrimSuffix(filepath.Base(playlistPath), ".m3u8") + "\n\n")
	}

	// Write Extended M3U Metadata: #EXTINF:<duration>,<artist> - <title>
	extInf := fmt.Sprintf("#EXTINF:%d,%s - %s\n", track.Duration, track.Artist, track.Title)
	f.WriteString(extInf)
	f.WriteString(track.Path + "\n\n")

	return nil
}
