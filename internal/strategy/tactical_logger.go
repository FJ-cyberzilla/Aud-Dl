package strategy

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type AntibotTacticalLogger struct {
	mu     sync.Mutex
	writer io.Writer
}

func NewAntibotTacticalLogger(out io.Writer) *AntibotTacticalLogger {
	if out == nil {
		out = os.Stdout
	}
	return &AntibotTacticalLogger{
		writer: out,
	}
}

// LogEvent records a strategy event with structured timestamp and level output
func (atl *AntibotTacticalLogger) LogEvent(level LogLevel, component, message string) error {
	atl.mu.Lock()
	defer atl.mu.Unlock()

	event := fmt.Sprintf("[%s] [%s] [%s] %s\n",
		time.Now().Format("15:04:05.000"),
		level.String(),
		component,
		message,
	)

	_, err := atl.writer.Write([]byte(event))
	if err != nil {
		return fmt.Errorf("failed to write log: %w", err)
	}
	return nil
}
