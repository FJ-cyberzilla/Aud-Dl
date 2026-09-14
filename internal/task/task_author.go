package task

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type JobPriority int

const (
	PriorityNormal JobPriority = iota
	PriorityHigh
	PriorityBatch
)

type AudioFormat string

const (
	FormatMP3V0  AudioFormat = "V0"
	FormatMP3320K AudioFormat = "320k"
	FormatOpus    AudioFormat = "160k"
)

type TaskManifest struct {
	ID         string      `json:"id"`
	SourceURL  string      `json:"source_url"`
	Artist     string      `json:"artist"`
	Title      string      `json:"title"`
	Format     AudioFormat `json:"format"`
	Priority   JobPriority `json:"priority"`
	CreatedAt  time.Time   `json:"created_at"`
	MaxRetries int         `json:"max_retries"`
}

type TaskAuthor struct{}

func NewTaskAuthor() *TaskAuthor {
	return &TaskAuthor{}
}

// AuthorManifest builds a normalized, low-complexity execution manifest for pipeline processing
func (ta *TaskAuthor) AuthorManifest(rawURL, artist, title string, format AudioFormat) (*TaskManifest, error) {
	if strings.TrimSpace(rawURL) == "" {
		return nil, fmt.Errorf("cannot author task: empty source URL")
	}
	if format == "" {
		format = FormatMP3V0 // Default to VBR MP3
	}
	id := generateTaskID()
	return &TaskManifest{
		ID:         id,
		SourceURL:  strings.TrimSpace(rawURL),
		Artist:     strings.TrimSpace(artist),
		Title:      strings.TrimSpace(title),
		Format:     format,
		Priority:   PriorityNormal,
		CreatedAt:  time.Now(),
		MaxRetries: 3,
	}, nil
}

func generateTaskID() string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
