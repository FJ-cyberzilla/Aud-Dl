package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStoreTrack(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vault_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	v, err := NewMusicVault(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create dummy file to move
	srcFile := filepath.Join(tmpDir, "src.mp3")
	if err := os.WriteFile(srcFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	meta := Metadata{
		Title:  "Track",
		Artist: "Artist",
		Album:  "Album",
	}

	dest, err := v.StoreTrack(srcFile, meta)
	if err != nil {
		t.Fatalf("StoreTrack failed: %v", err)
	}

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Errorf("File was not moved to %s", dest)
	}
}

func TestResolveTargetPath(t *testing.T) {
	v := &MusicVault{BasePath: "/tmp/vault"}

	tests := []struct {
		name     string
		meta     Metadata
		expected string
	}{
		{
			name: "Standard Routing",
			meta: Metadata{Artist: "Artist", Album: "Album"},
			expected: "/tmp/vault/Artists/Artist/Album",
		},
		{
			name: "Various Artists",
			meta: Metadata{IsVarious: true, Genre: "Rock"},
			expected: "/tmp/vault/Compilations_by_Genre/Rock",
		},
		{
			name: "Unknown Artist",
			meta: Metadata{Artist: "", Album: ""},
			expected: "/tmp/vault/Artists/Unknown Artist/Singles & EPs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := v.resolveTargetPath(tt.meta)
			if !strings.HasSuffix(got, tt.expected) {
				t.Errorf("resolveTargetPath() = %v, want suffix %v", got, tt.expected)
			}
		})
	}
}
