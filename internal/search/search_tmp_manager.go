package search

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type TempArtifact struct {
	ID        string
	Path      string
	CreatedAt time.Time
}

type SearchTmpManager struct {
	mu       sync.Mutex
	tmpDir   string
	ttl      time.Duration
	registry map[string]*TempArtifact
}

func NewSearchTmpManager(baseDir string, ttl time.Duration) (*SearchTmpManager, error) {
	tmpPath := filepath.Join(baseDir, ".tmp_search")
	if err := os.MkdirAll(tmpPath, 0755); err != nil {
		return nil, err
	}
	return &SearchTmpManager{
		tmpDir:   tmpPath,
		ttl:      ttl,
		registry: make(map[string]*TempArtifact),
	}, nil
}

// WriteScratchpad dumps temporary search metadata to disk for low-memory processing
func (tm *SearchTmpManager) WriteScratchpad(id string, data []byte) (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	filePath := filepath.Join(tm.tmpDir, id+".tmp")
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return "", err
	}
	tm.registry[id] = &TempArtifact{
		ID:        id,
		Path:      filePath,
		CreatedAt: time.Now(),
	}
	return filePath, nil
}

// PurgeExpired removes stale search artifacts older than the TTL
func (tm *SearchTmpManager) PurgeExpired() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	now := time.Now()
	for id, artifact := range tm.registry {
		if now.Sub(artifact.CreatedAt) > tm.ttl {
			_ = os.Remove(artifact.Path)
			delete(tm.registry, id)
		}
	}
}
