package gate

import (
	"fmt"
	"net/url"
	"strings"
)

type HubNormalizer struct{}

func NewHubNormalizer() *HubNormalizer {
	return &HubNormalizer{}
}

// NormalizeURL strips tracking parameters, canonicalizes protocols, and validates host syntax
func (hn *HubNormalizer) NormalizeURL(rawURL string) (string, error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return "", fmt.Errorf("empty URL provided")
	}

	// Default to https if protocol is omitted
	if !strings.HasPrefix(trimmed, "http://") && !strings.HasPrefix(trimmed, "https://") {
		trimmed = "https://" + trimmed
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	parsed.Host = strings.ToLower(parsed.Host)

	// Clean tracking parameters (UTM, referral codes, affiliate trackers)
	q := parsed.Query()
	for param := range q {
		if strings.HasPrefix(param, "utm_") || param == "fbclid" || param == "gclid" || param == "ref" {
			q.Del(param)
		}
	}
	parsed.RawQuery = q.Encode()

	return parsed.String(), nil
}

// SanitizeHeaders strips leaking internal environment headers
func (hn *HubNormalizer) SanitizeHeaders(headers map[string]string) map[string]string {
	cleaned := make(map[string]string)
	for k, v := range headers {
		lower := strings.ToLower(k)
		if lower == "x-forwarded-for" || lower == "via" || lower == "proxy-connection" {
			continue
		}
		cleaned[k] = v
	}
	return cleaned
}
