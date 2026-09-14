package gate

import (
	"encoding/json"
	"os"
	"testing"
)

func TestTmpLogManager_HealthCheck(t *testing.T) {
	tempDir := t.TempDir()
	tlm := NewTmpLogManager(tempDir)

	if err := tlm.RecordHealthQuestCheck(); err != nil {
		t.Fatalf("Failed to record health check: %v", err)
	}

	if _, err := os.Stat(tlm.LogPath); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
}

func TestTmpLogManager_ResumeState(t *testing.T) {
	tempDir := t.TempDir()
	tlm := NewTmpLogManager(tempDir)
	streamID := "test-stream"
	state := map[string]string{"position": "100"}

	if err := tlm.SaveResumeState(streamID, state); err != nil {
		t.Fatalf("Failed to save state: %v", err)
	}

	raw, exists := tlm.GetResumeState(streamID)
	if !exists {
		t.Fatal("State should exist")
	}

	var retrieved map[string]string
	if err := json.Unmarshal(raw, &retrieved); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if retrieved["position"] != "100" {
		t.Errorf("Expected position 100, got %s", retrieved["position"])
	}

	if err := tlm.ClearResumeState(streamID); err != nil {
		t.Fatalf("Failed to clear state: %v", err)
	}

	_, exists = tlm.GetResumeState(streamID)
	if exists {
		t.Error("State should have been cleared")
	}
}
