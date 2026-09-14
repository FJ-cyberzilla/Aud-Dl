package ui

import (
	"fmt"
	"strings"
	"time"
)

type PlatformDisplay struct {
	spinnerFrames []string
	frameIndex    int
}

func NewPlatformDisplay() *PlatformDisplay {
	return &PlatformDisplay{
		spinnerFrames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		frameIndex:    0,
	}
}

type TermInfo struct {
	IsTTY  bool
	Width  int
	Height int
}

func DetectTerminal() TermInfo {
	return TermInfo{IsTTY: true, Width: 80, Height: 24}
}

// RenderProgressBar draws a dynamic terminal progress bar with calculated ETA
func (pd *PlatformDisplay) RenderProgressBar(current, total int64, startTime time.Time) string {
	if total <= 0 {
		total = 100
	}
	if current > total {
		current = total
	}
	percentage := float64(current) / float64(total) * 100.0
	barWidth := 30
	filledWidth := int(float64(barWidth) * (float64(current) / float64(total)))
	bar := strings.Repeat("█", filledWidth) + strings.Repeat("░", barWidth-filledWidth)

	// Calculate ETA
	elapsed := time.Since(startTime)
	var eta string
	if current > 0 && elapsed > 0 {
		estimatedTotal := time.Duration(float64(elapsed) / (float64(current) / float64(total)))
		remaining := estimatedTotal - elapsed
		if remaining < 0 {
			remaining = 0
		}
		eta = remaining.Round(time.Second).String()
	} else {
		eta = "calculating..."
	}

	frame := pd.spinnerFrames[pd.frameIndex%len(pd.spinnerFrames)]
	pd.frameIndex++

	return fmt.Sprintf("\r%s [%s] %.1f%% (%d/%d) ETA: %s", frame, bar, percentage, current, total, eta)
}
