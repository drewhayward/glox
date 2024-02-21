package lox

import (
	"bytes"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestOutut(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		output string
	}{
		{"empty", "", ""},
		{"exp statment", "1;", ""},
		{"print", "print 1;", "1\n"},
		{"block statement",
			`var a = 1;
        {
            var a = 2;
            print a;
        }
        print a;
        `,
			"2\n1\n"},
		{"if statement: true",
			`if (true)
            print 1;
        else
            print 2;
        `,
			"1\n"},
		{"if statement: false",
			`if (false)
            print 1;
        else
            print 2;
        `,
			"2\n"},
		{"while loop",
			`
            var a = 1;
            while (a <= 5) {
                print a;
                a = a + 1;
            }
        `,
			"1\n2\n3\n4\n5\n"},
		{"for loop",
			`
            for (var a = 1; a <= 5; a = a + 1) {
                print a;
            }
        `,
			"1\n2\n3\n4\n5\n"},
		{"logical or shortcircuit",
			"print 1 or (1 / 0);",
			"true\n"},
		{"logical and shortcircuit",
			"print false and (1 / 0);",
			"false\n"},
		{"assignment as expression",
			"var a; print a = 1;",
			"1\n"},
	}
	is := is.New(t)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			s := RuntimeState{
				CurrEnv:   NewScopeEnv(nil),
				OutWriter: &buf,
			}

			s.Run(tc.input)

			is.Equal(
				strings.ReplaceAll(string(buf.Bytes()), "\n", "\\n"),
				strings.ReplaceAll(tc.output, "\n", "\\n"),
			)
		})
	}
}
