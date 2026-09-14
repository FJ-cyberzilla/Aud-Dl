package ui

import (
	"errors"
	"fmt"

	"audio-command-center/internal/processor"
)

// RenderPipelineError inspects the error and prints a structured diagnostic panel if it's a PipelineError.
func RenderPipelineError(target string, err error) {
	if err == nil {
		fmt.Println("\033[32m[PASS] Pipeline integrity check and download completed successfully.\033[0m")
		return
	}

	var pipeErr *processor.PipelineError

	// Idiomatic Go check using errors.As
	if errors.As(err, &pipeErr) {
		fmt.Println("\n┌────────────────────────────────────────────────────────┐")
		fmt.Printf("│ \033[31m[PIPELINE REJECTION REPORT]\033[0m                           │\n")
		fmt.Println("├────────────────────────────────────────────────────────┤")
		fmt.Printf("│ Target:      %-41s │\n", target)
		fmt.Printf("│ Category:    %-41s │\n", pipeErr.Category)
		fmt.Printf("│ Subsystem:   %-41s │\n", pipeErr.Subsystem)
		fmt.Println("├────────────────────────────────────────────────────────┤")
		fmt.Printf("│ Root Cause:  %-41s │\n", pipeErr.Message)
		fmt.Printf("│ Action:      %-41s │\n", pipeErr.Remediation)
		fmt.Println("└────────────────────────────────────────────────────────┘")
	} else {
		// Fallback for native unhandled standard library errors
		fmt.Printf("\n\033[31m[UNKNOWN ERROR] Target: %s | Details: %v\033[0m\n", target, err)
	}
}
