package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPlaylistGenerator(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vault_test_playlist")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	pg, err := NewPlaylistGenerator(tmpDir)
	if err != nil {
		t.Fatalf("NewPlaylistGenerator failed: %v", err)
	}

	// Create dummy mp3
	mp3File := filepath.Join(tmpDir, "test_track.mp3")
	if err := os.WriteFile(mp3File, []byte("fake_mp3_data"), 0644); err != nil {
		t.Fatal(err)
	}

	// Test AppendTrack
	playlistFile := filepath.Join(pg.PlaylistsDir, "test.m3u8")
	if err := pg.AppendTrack(playlistFile, mp3File); err != nil {
		t.Fatalf("AppendTrack failed: %v", err)
	}

	// Verify content
	content, err := os.ReadFile(playlistFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if !strings.Contains(string(content), "#EXTM3U") {
		t.Errorf("Playlist missing #EXTM3U header")
	}
	if !strings.Contains(string(content), "test_track.mp3") {
		t.Errorf("Playlist missing track path")
	}
}
