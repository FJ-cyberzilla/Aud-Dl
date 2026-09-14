package task

import (
	"audio-command-center/internal/processor"
	"audio-command-center/internal/vault"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type TaskDivision struct {
	mu            sync.Mutex
	bindings      *SystemBindings
	queue         chan *TaskManifest
	activeWorkers int
}

func NewTaskDivision(bindings *SystemBindings, queueSize int) *TaskDivision {
	return &TaskDivision{
		bindings: bindings,
		queue:    make(chan *TaskManifest, queueSize),
	}
}

// GetActiveWorkers returns the number of currently active task workers
func (td *TaskDivision) GetActiveWorkers() int {
	td.mu.Lock()
	defer td.mu.Unlock()
	return td.activeWorkers
}

func (td *TaskDivision) Enqueue(ctx context.Context, manifest *TaskManifest) error {
	select {
	case td.queue <- manifest:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("task division queue is full")
	}
}

// StartWorkerPool spins up worker goroutines to process incoming task manifests
func (td *TaskDivision) StartWorkerPool(ctx context.Context, concurrency int) {
	for i := 0; i < concurrency; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case manifest, ok := <-td.queue:
					if !ok {
						return
					}
					td.processManifest(ctx, manifest)
				}
			}
		}()
	}
}

func (td *TaskDivision) processManifest(ctx context.Context, manifest *TaskManifest) {
	td.mu.Lock()
	td.activeWorkers++
	td.mu.Unlock()

	defer func() {
		td.mu.Lock()
		td.activeWorkers--
		td.mu.Unlock()
	}()

	// 1. Fetch via GateHub
	filename := fmt.Sprintf("%s - %s", manifest.Artist, manifest.Title)
	stream, err := td.bindings.GateHub.IngestAndSanitizeStream(ctx, manifest.SourceURL, filename)
	if err != nil {
		log.Printf("TaskDivision: Fetch failed for %s: %v", manifest.SourceURL, err)
		return
	}

	// 2. Transcode to Temp File
	tmpFile, err := os.CreateTemp("", "aud-dl-*.mp3")
	if err != nil {
		log.Printf("TaskDivision: Failed to create temp file: %v", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	transcodedStream, err := td.bindings.Transcoder.StreamTranscode(ctx, stream, processor.TranscodeProfile(manifest.Format))
	if err != nil {
		tmpFile.Close()
		log.Printf("TaskDivision: Transcoding failed: %v", err)
		return
	}
	_, err = io.Copy(tmpFile, transcodedStream)
	transcodedStream.Close()
	tmpFile.Close()
	if err != nil {
		log.Printf("TaskDivision: Failed to copy stream to temp file: %v", err)
		return
	}

	// 3. Tag
	payload := processor.TagPayload{
		Title:   manifest.Title,
		Artist:  manifest.Artist,
		AppName: "Aud-Dl",
	}
	if err := td.bindings.Tagger.InjectMetadataAndBrand(tmpFile.Name(), payload); err != nil {
		log.Printf("TaskDivision: Tagging failed: %v", err)
		return
	}

	// 4. Ingest (VaultAdmin)
	pcmData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		log.Printf("TaskDivision: Failed to read temp file for ingestion: %v", err)
		return
	}

	strategy := vault.StrategyArtistAlbum
	targetPath, isDuplicate, err := td.bindings.VaultAdmin.IngestTrack(
		manifest.Artist,
		"", // Album
		"", // Genre
		filepath.Base(tmpFile.Name()),
		0,
		pcmData,
		strategy,
	)
	if err != nil {
		log.Printf("TaskDivision: Ingestion failed: %v", err)
		return
	}
	if isDuplicate {
		log.Printf("TaskDivision: Track %s is a duplicate, skipping.", filename)
		return
	}

	// Move file to target
	if err := os.Rename(tmpFile.Name(), targetPath); err != nil {
		// Fallback: Copy if rename fails (e.g. cross-device)
		data, err := os.ReadFile(tmpFile.Name())
		if err == nil {
			err = os.WriteFile(targetPath, data, 0644)
		}
		if err != nil {
			log.Printf("TaskDivision: Failed to move file to vault: %v", err)
		}
	}
}
