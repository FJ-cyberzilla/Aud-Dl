package processor

import "fmt"

type ErrorCategory string

const (
	CatSecurity ErrorCategory = "SECURITY_BLOCK"
	CatCorrupt  ErrorCategory = "STREAM_CORRUPTION"
	CatNetwork  ErrorCategory = "NETWORK_FAILURE"
	CatUnknown  ErrorCategory = "UNKNOWN_ANOMALY"
)

type PipelineError struct {
	Category    ErrorCategory
	Subsystem   string
	Message     string
	Remediation string
	Err         error
}

func (pe *PipelineError) Error() string {
	if pe.Err != nil {
		return fmt.Sprintf("[%s] %s (%s) -> caused by: %v", pe.Category, pe.Message, pe.Subsystem, pe.Err)
	}
	return fmt.Sprintf("[%s] %s (%s)", pe.Category, pe.Message, pe.Subsystem)
}

func (pe *PipelineError) Unwrap() error {
	return pe.Err
}
