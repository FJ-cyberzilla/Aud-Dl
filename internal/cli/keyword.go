package cli

import "strings"

type CommandType string

const (
	CmdUnknown CommandType = "UNKNOWN"
	CmdQuit    CommandType = "QUIT"
	CmdCancel  CommandType = "CANCEL"
	CmdStop    CommandType = "STOP"
	CmdResume  CommandType = "RESUME"
	CmdNext    CommandType = "NEXT"
	CmdPrev    CommandType = "PREV"
)

type KeywordBroker struct{}

func NewKeywordBroker() *KeywordBroker { return &KeywordBroker{} }

func (kb *KeywordBroker) ParseInput(input string) CommandType {
	switch strings.TrimSpace(strings.ToLower(input)) {
	case "q", "quit", "exit": return CmdQuit
	case "esc", "c", "cancel": return CmdCancel
	case "s", "stop", "pause": return CmdStop
	case "r", "resume": return CmdResume
	case "n", "next": return CmdNext
	case "p", "prev": return CmdPrev
	default: return CmdUnknown
	}
}
