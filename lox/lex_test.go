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
		{"-", []TokenType{MINUS}},
		{"+", []TokenType{PLUS}},
		{";", []TokenType{SEMICOLON}},
		{"/", []TokenType{SLASH}},
		{"*", []TokenType{STAR}},
		{"!", []TokenType{BANG}},
		{"!=", []TokenType{BANG_EQUAL}},
		{"=", []TokenType{EQUAL}},
		{"==", []TokenType{EQUAL_EQUAL}},
		{">=", []TokenType{GREATER_EQUAL}},
		{">", []TokenType{GREATER}},
		{"<", []TokenType{LESS}},
		{"<=", []TokenType{LESS_EQUAL}},
		{"\"testing\"", []TokenType{STRING}},
		{"123", []TokenType{NUMBER}},
		{"123.", []TokenType{NUMBER, DOT}},
		{"123.09", []TokenType{NUMBER}},
		{"testing", []TokenType{IDENTIFIER}},
		{"for", []TokenType{FOR}},
		{"and", []TokenType{AND}},
		{"class", []TokenType{CLASS}},
		{"else", []TokenType{ELSE}},
		{"false", []TokenType{FALSE}},
		{"fun", []TokenType{FUN}},
		{"for", []TokenType{FOR}},
		{"if", []TokenType{IF}},
		{"nul", []TokenType{NUL}},
		{"or", []TokenType{OR}},
		{"print", []TokenType{PRINT}},
		{"return", []TokenType{RETURN}},
		{"super", []TokenType{SUPER}},
		{"this", []TokenType{THIS}},
		{"true", []TokenType{TRUE}},
		{"var", []TokenType{VAR}},
		{"while", []TokenType{WHILE}},
		{"var test = \"foobar\"", []TokenType{VAR, IDENTIFIER, EQUAL, STRING}},
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
