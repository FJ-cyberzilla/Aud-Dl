package task

import (
	"context"
	"fmt"
	"audio-command-center/internal/cache"
	"audio-command-center/internal/cli"
	"audio-command-center/internal/ui"
)

type TaskController struct {
	author   *TaskAuthor
	bindings *SystemBindings
	division *TaskDivision
	theme    ui.Theme
	cache    *cache.CacheManager
}

func NewTaskController(bindings *SystemBindings, division *TaskDivision, theme ui.Theme, cache *cache.CacheManager) *TaskController {
	return &TaskController{
		author:   NewTaskAuthor(),
		bindings: bindings,
		division: division,
		theme:    theme,
		cache:    cache,
	}
}

// DispatchDownload authors a task manifest and queues it for execution
func (tc *TaskController) DispatchDownload(ctx context.Context, rawURL, artist, title string) (string, error) {
	manifest, err := tc.author.AuthorManifest(rawURL, artist, title, FormatMP3V0)
	if err != nil {
		return "", fmt.Errorf("task controller dispatch failed: %w", err)
	}

	if err := tc.division.Enqueue(ctx, manifest); err != nil {
		return "", fmt.Errorf("failed to enqueue manifest [%s]: %w", manifest.ID, err)
	}

	return manifest.ID, nil
}

func (tc *TaskController) GetDivision() *TaskDivision {
	return tc.division
}

// RunInteractiveSession starts the interactive CLI loop.
func (tc *TaskController) RunInteractiveSession(ctx context.Context) {
	menu := cli.NewInteractiveMenu(tc.theme, tc.cache)
	for {
		_, err := menu.PromptMainMenu(ctx)
		if err != nil {
			fmt.Printf("Menu error: %v\n", err)
			return
		}
	}
}
