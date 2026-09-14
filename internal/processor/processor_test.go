package processor

import (
	"errors"
	"strings"
	"testing"
)

func TestPipelineError(t *testing.T) {
	rootErr := errors.New("connection reset")
	pe := &PipelineError{
		Category:    CatNetwork,
		Subsystem:   "gate",
		Message:     "failed to connect",
		Remediation: "check internet",
		Err:         rootErr,
	}

	errMsg := pe.Error()
	if !strings.Contains(errMsg, "[NETWORK_FAILURE]") {
		t.Errorf("Error string missing category, got: %s", errMsg)
	}
	
	if !errors.Is(pe, rootErr) {
		t.Error("Unwrap should return root error")
	}
}
