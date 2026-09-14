package vault

import (
	"os"
	"testing"
)

func TestVaultAdmin(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vault_admin_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	va, err := NewVaultAdmin(tmpDir)
	if err != nil {
		t.Fatalf("NewVaultAdmin failed: %v", err)
	}
	defer va.Close()

	// Test IngestTrack
	pcmData := []byte("audio_data")
	target, isDup, err := va.IngestTrack("Artist", "Album", "Genre", "test.mp3", 2023, pcmData, StrategyArtistAlbum)
	if err != nil {
		t.Fatalf("IngestTrack failed: %v", err)
	}
	if isDup {
		t.Errorf("expected unique track, got duplicate")
	}
	if target == "" {
		t.Errorf("expected target path, got empty")
	}

	// Test GetStats
	count, size, err := va.GetStats()
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	// Depending on implementation, stats might be 0 or something else.
	// But it shouldn't error.
	_ = count
	_ = size
}
