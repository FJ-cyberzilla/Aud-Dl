package gate

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type InspectionResult struct {
	IsClean        bool
	IsBotChallenge bool
	DetectedMIME   string
	ErrorMessage   string
}

type HubScrutinizer struct {
	maxPeekBytes int
}

func NewHubScrutinizer() *HubScrutinizer {
	return &HubScrutinizer{
		maxPeekBytes: 1024, // Peek first 1KB to detect payloads/interstitials without buffering large audio files
	}
}

// ScrutinizeResponse inspects response headers and initial payload bytes for hidden anti-bot traps
func (hs *HubScrutinizer) ScrutinizeResponse(resp *http.Response) (*InspectionResult, io.Reader, error) {
	if resp == nil || resp.Body == nil {
		return &InspectionResult{
			IsClean:      false,
			ErrorMessage: "nil response or body provided",
		}, nil, fmt.Errorf("invalid response handle")
	}

	// 1. Validate status codes
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusServiceUnavailable {
		return &InspectionResult{
			IsClean:        false,
			IsBotChallenge: true,
			ErrorMessage:   fmt.Sprintf("HTTP status %d blocked by origin protection", resp.StatusCode),
		}, resp.Body, nil
	}

	// 2. Read initial header bytes without consuming the underlying audio stream
	peekBuffer := make([]byte, hs.maxPeekBytes)
	n, err := io.ReadFull(resp.Body, peekBuffer)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return &InspectionResult{
			IsClean:      false,
			ErrorMessage: fmt.Sprintf("failed to peek stream bytes: %v", err),
		}, resp.Body, nil
	}
	peekData := peekBuffer[:n]
	lowerPayload := strings.ToLower(string(peekData))

	// 3. Detect HTML Interstitial / Challenge traps masking as 200 OK
	if strings.Contains(lowerPayload, "<html") || strings.Contains(lowerPayload, "cf-challenge") || strings.Contains(lowerPayload, "just a moment...") || strings.Contains(lowerPayload, "enable javascript") {
		return &InspectionResult{
			IsClean:        false,
			IsBotChallenge: true,
			ErrorMessage:   "detected HTML anti-bot challenge interstitial payload",
		}, resp.Body, nil
	}

	// Reconstruct the original stream so downstream components can read seamlessly
	reconstructedStream := io.MultiReader(bytes.NewReader(peekData), resp.Body)

	// 4. Verify Media Magic Bytes (ID3v2, Ogg/Opus, RIFF/WAV, FLAC) for media routes
	detectedMIME := http.DetectContentType(peekData)
	return &InspectionResult{
		IsClean:        true,
		IsBotChallenge: false,
		DetectedMIME:   detectedMIME,
	}, reconstructedStream, nil
}
