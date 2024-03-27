package lox

import (
	"fmt"
	"github.com/gkampitakis/go-snaps/snaps"
	"testing"
)

func TestParseSnapshot(t *testing.T) {
	testCases := []struct {
		source string
	}{
		{`
            var a = 1;
        `},
	}

	for _, tc := range testCases {
		s := snaps.WithConfig()
		t.Run(fmt.Sprintf("Parse(%q)", tc.source), func(t *testing.T) {
			tokens, _ := ScanTokens(tc.source)
			node, err := Parse(tokens)
			if err != nil {
				t.Fatalf(err.Error())
			}

			s.MatchSnapshot(t, node)
		})
	}
}
