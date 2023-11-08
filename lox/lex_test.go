package lox

import (
	"fmt"
	"testing"
)

func TestScanTokens(t *testing.T) {
	testCases := []struct {
		source   string
		expected []TokenType
	}{
		{"()", []TokenType{LEFT_PAREN, RIGHT_PAREN}},
		{"{}", []TokenType{LEFT_BRACE, RIGHT_BRACE}},
		{"/", []TokenType{SLASH}},
		{".", []TokenType{DOT}},
		{",", []TokenType{COMMA}},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("ScanTokens(%q)", tc.source), func(t *testing.T) {
			tokens, _ := ScanTokens(tc.source)
			if len(tokens) != len(tc.expected) {
				t.Errorf("ScanTokens(%q) returned %d tokens, expected %d tokens", tc.source, len(tokens), len(tc.expected))
				return
			}

			for i, token := range tokens {
				if token.type_ != tc.expected[i] {
					t.Errorf("ScanTokens(%q) is wrong: got %s, expected %s", tc.source, token, tc.expected[i])
				}
			}
		})
	}
}
