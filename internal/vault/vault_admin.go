package vault

import (
	"fmt"
	"os"
	"path/filepath"
)


// SmartVaultRouter directs audio payloads to optimal vault partitions.
type SmartVaultRouter struct {
	baseDir string
}

func NewSmartVaultRouter(baseDir string) *SmartVaultRouter {
	return &SmartVaultRouter{baseDir: baseDir}
}

type VaultAdmin struct {
	baseDir  string
	resolver *VaultResolver
	detector *PersistentDetector
	playlist *PlaylistGenerator
}

func NewVaultAdmin(vaultDir string) (*VaultAdmin, error) {
	if err := os.MkdirAll(vaultDir, 0755); err != nil {
		return nil, fmt.Errorf("vault_admin: failed to create base dir: %w", err)
	}

	resolver := NewVaultResolver(vaultDir)
	detector := NewPersistentDetector()
	playlist, err := NewPlaylistGenerator(vaultDir)
	if err != nil {
		return nil, fmt.Errorf("vault_admin: failed to create playlist generator: %w", err)
	}

	return &VaultAdmin{
		baseDir:  vaultDir,
		resolver: resolver,
		detector: detector,
		playlist: playlist,
	}, nil
}

// GetStats calculates vault statistics by walking the base directory
func (va *VaultAdmin) GetStats() (int, float64, error) {
	trackCount := 0
	var totalSize int64

	err := filepath.Walk(va.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".mp3" {
			trackCount++
			totalSize += info.Size()
		}
		return nil
	})

	return trackCount, float64(totalSize) / 1024 / 1024, err
}

// IngestTrack checks for duplicates, computes destination path, ensures directories exist, and indexes fingerprint
func (va *VaultAdmin) IngestTrack(artist, album, genre, filename string, year int, pcmData []byte, strategy LayoutStrategy) (string, bool, error) {
	// 1. Acoustic Fingerprint Duplication Check
	isDuplicate, existingPath, err := va.detector.CheckAndIndex(artist, album, filename, pcmData)
	if err != nil {
		return "", false, fmt.Errorf("vault_admin: fingerprint verification failed: %w", err)
	}
	if isDuplicate {
		return existingPath, true, nil
	}

	// 2. Resolve target storage path
	targetPath, err := va.resolver.ResolveDestination(artist, album, genre, filename, year, strategy)
	if err != nil {
		return "", false, err
	}

	// 3. Ensure target directory structure exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return "", false, fmt.Errorf("vault_admin: failed to create directory tree: %w", err)
	}

	return targetPath, false, nil
}

// RegisterToPlaylist appends processed audio file into dynamic UTF-8 M3U8 target
func (va *VaultAdmin) RegisterToPlaylist(playlistName, trackPath string) error {
	playlistTarget := va.resolver.ResolvePlaylistPath(playlistName)
	return va.playlist.AppendTrack(playlistTarget, trackPath)
}

func (va *VaultAdmin) Close() error {
	return va.detector.Close()
}
