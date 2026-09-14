package vault

import (
	"testing"
)

func TestSanitizeInternationalTitle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "Hello World"},
		{"Hello/World", "Hello_World"},
		{"你好世界", "你好世界"}, // Chinese
		{"привет", "привет"},   // Cyrillic
		{"invalid/chars: *?", "invalid_chars_ _"},
		{"  trimmed  ", "trimmed"},
		{"multiple__underscores", "multiple_underscores"},
		{"", "unknown_track"},
	}

	for _, test := range tests {
		actual := SanitizeInternationalTitle(test.input)
		// Debugging: let's see why it failed
		if actual != test.expected {
			t.Errorf("Input '%s': expected '%s', got '%s'", test.input, test.expected, actual)
		}
	}
}
