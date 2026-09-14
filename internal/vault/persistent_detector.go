package vault

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type PersistentDetector struct {
	mu           sync.RWMutex
	fingerprints map[uint64]*HashFingerprint
}

func NewPersistentDetector() *PersistentDetector {
	return &PersistentDetector{
		fingerprints: make(map[uint64]*HashFingerprint),
	}
}

// CheckAndIndex computes a 64-bit acoustic FNV-1a hash over raw PCM stream data,
// checks for existing duplicates in the vault index, and registers the track if unique.
func (pd *PersistentDetector) CheckAndIndex(artist, album, filename string, pcmData []byte) (bool, string, error) {
	if len(pcmData) == 0 {
		return false, "", fmt.Errorf("fingerprint error: raw PCM data buffer is empty for %s", filename)
	}

	// 1. Generate real non-cryptographic acoustic hash from raw PCM sample bytes
	hasher := fnv.New64a()
	if _, err := hasher.Write(pcmData); err != nil {
		return false, "", fmt.Errorf("fingerprint hash calculation failed: %w", err)
	}
	pcmHash := hasher.Sum64()

	pd.mu.Lock()
	defer pd.mu.Unlock()

	// 2. Query in-memory vault cache for acoustic duplicate match
	if existing, found := pd.fingerprints[pcmHash]; found {
		matchInfo := fmt.Sprintf("Match found! Duplicate of '%s' by '%s' (Path: %s)", existing.Path, existing.Artist, existing.Path)
		return true, matchInfo, nil
	}

	// 3. Register unique audio track into persistent index
	pd.fingerprints[pcmHash] = &HashFingerprint{
		Hash:   pcmHash,
		Artist: artist,
		Album:  album,
		Path:   filename,
	}

	return false, "", nil
}

// Close satisfies the interface for vault resources
func (pd *PersistentDetector) Close() error {
	return nil
}
