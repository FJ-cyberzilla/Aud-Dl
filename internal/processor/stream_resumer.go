package processor

import (
	"audio-command-center/internal/gate"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ResumeSession struct {
	StreamID   string `json:"stream_id"`
	FilePath   string `json:"file_path"`
	BytesRead  int64  `json:"bytes_read"`
	TotalBytes int64  `json:"total_bytes"`
	URL        string `json:"url"`
	UpdatedAt  int64  `json:"updated_at"`
}

type StreamResumer struct {
	// Assumes TmpLogManager exists in internal/gate
	tmpLog *gate.TmpLogManager
}

func NewStreamResumer(tmpLog *gate.TmpLogManager) *StreamResumer {
	return &StreamResumer{tmpLog: tmpLog}
}

// SaveCheckpoint writes download progress to runtime log
func (sr *StreamResumer) SaveCheckpoint(session ResumeSession) error {
	session.UpdatedAt = time.Now().Unix()
	return sr.tmpLog.SaveResumeState(session.StreamID, session)
}

// ResumeDownload fetches remaining byte range if checkpoint exists and is under 24h old
func (sr *StreamResumer) ResumeDownload(ctx context.Context, streamID string) (io.ReadCloser, int64, error) {
	raw, found := sr.tmpLog.GetResumeState(streamID) // Automatic 24-Hour Expiration Check
	if !found {
		return nil, 0, fmt.Errorf("resumer: session not found")
	}

	var session ResumeSession
	if err := json.Unmarshal(raw, &session); err != nil {
		return nil, 0, fmt.Errorf("resumer: failed to unmarshal session: %w", err)
	}

	if time.Now().Unix()-session.UpdatedAt > 86400 {
		_ = sr.tmpLog.ClearResumeState(streamID) // Auto-cleanup stale session
		return nil, 0, fmt.Errorf("resumer: session expired")
	}
	req, err := http.NewRequestWithContext(ctx, "GET", session.URL, nil)
	if err != nil {
		return nil, 0, err
	}
	// Request byte range offset
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", session.BytesRead))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	return resp.Body, session.BytesRead, nil
}
