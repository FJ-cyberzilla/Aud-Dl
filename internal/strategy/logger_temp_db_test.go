package strategy

import (
	"context"
	"os"
	"testing"
)

func TestLoggerTempDB_Lifecycle(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	maxRows := 5
	logger, err := NewLoggerTempDB(tmpDir, maxRows)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Test Push
	for i := 0; i < 10; i++ {
		err := logger.PushLog(ctx, "INFO", "test", "msg")
		if err != nil {
			t.Errorf("PushLog failed: %v", err)
		}
	}

	// Test Fetch
	events, err := logger.FetchRecent(ctx, 10)
	if err != nil {
		t.Fatalf("FetchRecent failed: %v", err)
	}

	if len(events) != maxRows {
		t.Errorf("Expected %d rows, got %d", maxRows, len(events))
	}

	// Test Close
	if err := logger.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Verify file removal
	if _, err := os.Stat(logger.dbPath); err == nil {
		t.Error("Database file still exists after Close")
	}
}
