package internal

import (
	"testing"
)

func TestNormalizeWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "regular spaces preserved",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "tabs preserved",
			input:    "hello\tworld",
			expected: "hello\tworld",
		},
		{
			name:     "newlines preserved",
			input:    "hello\nworld",
			expected: "hello\nworld",
		},
		{
			name:     "carriage returns removed",
			input:    "hello\rworld",
			expected: "helloworld",
		},
		{
			name:     "CRLF converted to LF",
			input:    "hello\r\nworld",
			expected: "hello\nworld",
		},
		{
			name:     "non-breaking space converted",
			input:    "hello\u00A0world",
			expected: "hello world",
		},
		{
			name:     "em space converted",
			input:    "hello\u2003world",
			expected: "hello world",
		},
		{
			name:     "en space converted",
			input:    "hello\u2002world",
			expected: "hello world",
		},
		{
			name:     "thin space converted",
			input:    "hello\u2009world",
			expected: "hello world",
		},
		{
			name:     "zero-width space converted",
			input:    "hello\u200Bworld",
			expected: "hello world",
		},
		{
			name:     "ideographic space converted",
			input:    "hello\u3000world",
			expected: "hello world",
		},
		{
			name:     "mixed whitespace",
			input:    "normal space\u00A0nbsp\u2003em space\ttab\nnewline",
			expected: "normal space nbsp em space\ttab\nnewline",
		},
		{
			name:     "multiple unusual spaces in a row",
			input:    "hello\u00A0\u2003\u2009world",
			expected: "hello   world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only regular whitespace",
			input:    " \t\n ",
			expected: " \t\n ",
		},
		{
			name:     "no whitespace",
			input:    "helloworld",
			expected: "helloworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeWhitespace(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeWhitespace(%q) = %q, want %q",
					tt.input, result, tt.expected)
			}
		})
	}
}
