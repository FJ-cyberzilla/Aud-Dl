package processor

import (
	"audio-command-center/internal/gate"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStreamResumer_SaveAndResume(t *testing.T) {
	tempDir := t.TempDir()
	tlm := gate.NewTmpLogManager(tempDir)
	sr := NewStreamResumer(tlm)

	streamID := "test-stream"
	session := ResumeSession{
		StreamID:   streamID,
		BytesRead:  100,
		TotalBytes: 1000,
		URL:        "http://example.com/file",
	}

	if err := sr.SaveCheckpoint(session); err != nil {
		t.Fatalf("Failed to save checkpoint: %v", err)
	}

	// Mock server to verify Range header
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rangeHeader := r.Header.Get("Range")
		if rangeHeader != "bytes=100-" {
			t.Errorf("Expected Range bytes=100-, got %s", rangeHeader)
		}
		w.WriteHeader(http.StatusPartialContent)
	}))
	defer ts.Close()

	session.URL = ts.URL
	// Update with new URL to point to mock server
	sr.SaveCheckpoint(session)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, read, err := sr.ResumeDownload(ctx, streamID)
	if err != nil {
		t.Fatalf("Resume failed: %v", err)
	}
	defer body.Close()

	if read != 100 {
		t.Errorf("Expected read 100, got %d", read)
	}
}
