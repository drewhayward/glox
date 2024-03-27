package lox

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestParseSnapshot(t *testing.T) {
	testCases := []struct {
		source string
	}{
		{`
1;
        `},
		{`
1 + 1;
        `},
		{`
var a = 1;
        `},
		{`
fun foo() {}
        `},
		{`
fun foo() {
print "hello";
}
        `},
		{`
var a = 1;
{
var a = 2;
print a;
}
print a;
        `},
		{`
class Foo {}
        `},
		{`
class Foo {
bar() {}
}
        `},
	}

	for _, tc := range testCases {
		s := snaps.WithConfig()

		t.Run(fmt.Sprintf("Parse(%q)", strings.Trim(tc.source, " \n")), func(t *testing.T) {
			tokens, err := ScanTokens(tc.source)
			if err != nil {
				t.Fatalf(err.Error())
				t.FailNow()
			}

			node, err := Parse(tokens)
			if err != nil {
				t.Fatalf(err.Error())
				t.FailNow()
			}

			s.MatchSnapshot(t, node)
		})
	}
}
