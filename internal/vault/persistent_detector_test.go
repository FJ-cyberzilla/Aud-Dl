package vault

import (
	"testing"
)

func TestPersistentDetector_CheckAndIndex(t *testing.T) {
	pd := NewPersistentDetector()

	// Test 1: Unique track
	isDup, _, err := pd.CheckAndIndex("Artist1", "Album1", "path1.mp3", []byte("audio_data_1"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if isDup {
		t.Errorf("expected unique track, got duplicate")
	}

	// Test 2: Duplicate track
	isDup, _, err = pd.CheckAndIndex("Artist1", "Album1", "path2.mp3", []byte("audio_data_1"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isDup {
		t.Errorf("expected duplicate track, got unique")
	}

	// Test 3: Unique track (different data)
	isDup, _, err = pd.CheckAndIndex("Artist2", "Album2", "path3.mp3", []byte("audio_data_2"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if isDup {
		t.Errorf("expected unique track, got duplicate")
	}

	// Test 4: Empty data error
	_, _, err = pd.CheckAndIndex("Artist", "Album", "path4.mp3", []byte(""))
	if err == nil {
		t.Errorf("expected error for empty data, got nil")
	}
}
