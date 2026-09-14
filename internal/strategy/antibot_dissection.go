package strategy

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

type DissectionResult struct {
	Type         ChallengeType
	StatusCode   int
	RequiresWait bool
}

type AntibotDissection struct{}

func NewAntibotDissection() *AntibotDissection {
	return &AntibotDissection{}
}

// InspectResponse reads response headers and body snippets to classify anti-bot blocks
func (ad *AntibotDissection) InspectResponse(resp *http.Response) DissectionResult {
	if resp == nil {
		return DissectionResult{Type: ChallengeNone}
	}
	result := DissectionResult{
		StatusCode: resp.StatusCode,
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		result.Type = ChallengeRateLimited
		result.RequiresWait = true
		return result
	}

	// Read first 2KB of response body for fingerprint inspection
	buf := make([]byte, 2048)
	n, _ := io.ReadFull(resp.Body, buf)
	bodySnippet := strings.ToLower(string(buf[:n]))

	// Restore body reader for downstream consumers
	resp.Body = struct {
		io.Reader
		io.Closer
	}{
		Reader: io.MultiReader(bytes.NewReader(buf[:n]), resp.Body),
		Closer: resp.Body,
	}

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusServiceUnavailable {
		if strings.Contains(bodySnippet, "just a moment...") || strings.Contains(bodySnippet, "cf-mitigated") {
			result.Type = ChallengeCloudflareTurnstile
			result.RequiresWait = true
			return result
		}
		if strings.Contains(bodySnippet, "ak_bmsc") || strings.Contains(bodySnippet, "access denied") {
			result.Type = ChallengeAkamaiBot
			result.RequiresWait = true
			return result
		}
		result.Type = ChallengeIPBlock
	}
	return result
}
