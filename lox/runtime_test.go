package lox

import (
	"bytes"
	"testing"
    is "github.com/matryer/is"

)

func TestOutut(t * testing.T) {
    cases := []struct{
        name string
        input string
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
        {"if statement", 
        `if (true)
            print 1;
        else
            print 2;
        `,
        "1\n"},
        {"if statement", 
        `if (false)
            print 1;
        else
            print 2;
        `,
        "2\n"},
    }
    is := is.New(t)

    for _, tc := range cases {
        t.Run(tc.name, func (t *testing.T) {
            var buf bytes.Buffer

            s := RuntimeState{
                CurrEnv: NewScopeEnv(nil),
                OutWriter: &buf,
            }            

            s.Run(tc.input)

            is.Equal(string(buf.Bytes()), tc.output)
        })
    }
}
