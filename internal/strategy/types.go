package strategy

import (
	"context"
	"time"
)

type LogLevel int

const (
	LogInfo LogLevel = iota
	LogWarn
	LogError
)

func (l LogLevel) String() string {
	switch l {
	case LogWarn:
		return "WARN"
	case LogError:
		return "ERR "
	default:
		return "INFO"
	}
}

func ParseLogLevel(s string) LogLevel {
	switch s {
	case "WARN":
		return LogWarn
	case "ERR ":
		return LogError
	default:
		return LogInfo
	}
}

// TacticalEvent represents a shared log entry.
type TacticalEvent struct {
	Timestamp time.Time
	Level     LogLevel
	Component string
	Message   string
}

type ChallengeType int

const (
	ChallengeNone ChallengeType = iota
	ChallengeRateLimited
	ChallengeCloudflareTurnstile
	ChallengeAkamaiBot
	ChallengeIPBlock
)

type ResponseAnalysis struct {
	Type       ChallengeType
	StatusCode int
}

type BehaviorShifter struct{}

func NewBehaviorShifter(min, max int64) *BehaviorShifter { return &BehaviorShifter{} }
func (s *BehaviorShifter) JitterWait(ctx context.Context) error { return nil }
