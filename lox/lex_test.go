package lox

import (
	"fmt"
	"github.com/gkampitakis/go-snaps/snaps"
	"testing"
)

func TestLexSnapshot(t *testing.T) {
	testCases := []struct {
		source string
	}{
		{"()"},
		{"{}"},
		{"/"},
		{"."},
		{","},
		{"-"},
		{"+"},
		{";"},
		{"/"},
		{"*"},
		{"!"},
		{"!="},
		{"="},
		{"=="},
		{">="},
		{">"},
		{"<"},
		{"<="},
		{"\"testing\""},
		{"123"},
		{"123."},
		{"123.09"},
		{"testing"},
		{"for"},
		{"and"},
		{"class"},
		{"else"},
		{"false"},
		{"fun"},
		{"for"},
		{"if"},
		{"nil"},
		{"or"},
		{"print"},
		{"return"},
		{"super"},
		{"this"},
		{"true"},
		{"var"},
		{"while "},
		{"var test = \"foobar\";"},
		{"var\nvar"},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("ScanTokens(%q)", tc.source), func(t *testing.T) {
			tokens, err := ScanTokens(tc.source)
			if err != nil {
				t.Fatal(err.Error())
			}
			snaps.MatchSnapshot(t, tokens)
		})
	}
}
