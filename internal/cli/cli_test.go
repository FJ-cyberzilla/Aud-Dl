package cli

import (
	"testing"
)

func TestKeywordBroker_ParseInput(t *testing.T) {
	kb := NewKeywordBroker()
	
	tests := []struct {
		input    string
		expected CommandType
	}{
		{"q", CmdQuit},
		{"QUIT", CmdQuit},
		{"s", CmdStop},
		{"n", CmdNext},
		{"p", CmdPrev},
		{"unknown", CmdUnknown},
	}
	
	for _, test := range tests {
		actual := kb.ParseInput(test.input)
		if actual != test.expected {
			t.Errorf("For input %s, expected %s, got %s", test.input, test.expected, actual)
		}
	}
}
