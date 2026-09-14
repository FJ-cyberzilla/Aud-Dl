package ui

import (
	"strings"

	"audio-command-center/internal/processor"
)

type DiagnosticMatch struct {
	Category    processor.ErrorCategory
	RootCause   string
	Remediation string
}

type DiagnosticRule func(errStr string) (DiagnosticMatch, bool)

// Registry of diagnostic rules making it trivial to add new error patterns later
var DiagnosticRegistry = []DiagnosticRule{
	func(errStr string) (DiagnosticMatch, bool) {
		if strings.Contains(errStr, "pipeline security error") || strings.Contains(errStr, "text/html") {
			return DiagnosticMatch{
				Category:    processor.CatSecurity,
				RootCause:   "Ingress stream returned HTML/JSON block or anti-bot challenge instead of audio.",
				Remediation: "Rotate proxy IP or update browser headers in strategy/.",
			}, true
		}
		return DiagnosticMatch{}, false
	},
	func(errStr string) (DiagnosticMatch, bool) {
		if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "connection refused") {
			return DiagnosticMatch{
				Category:    processor.CatNetwork,
				RootCause:   "Target provider dropped connection or network route timed out.",
				Remediation: "Increase pacing interval or switch fallback provider.",
			}, true
		}
		return DiagnosticMatch{}, false
	},
}

func LookupDiagnostics(err error) DiagnosticMatch {
	if err == nil {
		return DiagnosticMatch{Category: "SUCCESS", RootCause: "None", Remediation: "None"}
	}

	errStr := err.Error()
	for _, rule := range DiagnosticRegistry {
		if match, ok := rule(errStr); ok {
			return match
		}
	}

	// Fallback for unknown anomalies
	return DiagnosticMatch{
		Category:    processor.CatUnknown,
		RootCause:   errStr,
		Remediation: "Inspect underlying stack trace or logs.",
	}
}
