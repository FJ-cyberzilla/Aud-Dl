package task

import (
	"context"
	"testing"
	"time"
)

func TestTaskDivision_Queue(t *testing.T) {
	td := NewTaskDivision(&SystemBindings{}, 1)
	ctx := context.Background()
	manifest := &TaskManifest{Artist: "Test", Title: "Test"}

	if err := td.Enqueue(ctx, manifest); err != nil {
		t.Fatalf("Failed to enqueue: %v", err)
	}

	// Queue should be full now
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()
	if err := td.Enqueue(ctx2, manifest); err == nil {
		t.Error("Expected error for full queue")
	}
}

func TestTaskDivision_ActiveWorkers(t *testing.T) {
	td := NewTaskDivision(&SystemBindings{}, 1)
	if workers := td.GetActiveWorkers(); workers != 0 {
		t.Errorf("Expected 0 workers, got %d", workers)
	}
}
