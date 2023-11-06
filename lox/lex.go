package lox

import (
	"errors"
	"fmt"
	"os"
)


type TokenType string
const (
    // Single char tokens
	LEFT_PAREN = "LEFT_PAREN"
	RIGHT_PAREN = "RIGHT_PAREN"
	LEFT_BRACE = "LEFT_BRACE"
	RIGHT_BRACE = "RIGHT_BRACE"
	COMMA = "COMMA"
	DOT = "DOT"
	MINUS = "MINUS"
	PLUS = "PLUS"
	SEMICOLON = "SEMICOLON"
	SLASH = "SLASH"
	STAR = "STAR"

    // One or two char tokens
	BANG = "BANG"
	BANG_EQUAL = "BANG_EQUAL"
	EQUAL = "EQUAL"
	EQUAL_EQUAL = "EQUAL_EQUAL"
	GREATER = "GREATER"
	GREATER_EQUAL = "GREATER_EQUAL"
	LESS = "LESS"
	LESS_EQUAL = "LESS_EQUAL"

    // Literals
	IDENTIFIER = "IDENTIFIER"
	STRING = "STRING"
	NUMBER = "NUMBER"

    // Keywords
	AND = "AND"
	CLASS = "CLASS"
	ELSE = "ELSE"
	FALSE = "FALSE"
	FUN = "FUN"
	FOR = "FOR"
	IF = "IF"
	NUL = "NUL"
	OR = "OR"
	PRINT = "PRINT"
	RETURN = "RETURN"
	SUPER = "SUPER"
	THIS = "THIS"
	TRUE = "TRUE"
	VAR = "VAR"
	WHILE = "WHILE"

	EOF = "EOF"
)

type Token struct {
    type_ TokenType 
    lexeme string
    literal string
    line int
}

func (t Token) String() string {
    return fmt.Sprintf("Token{%v '%s' %s}", t.type_, t.lexeme, t.literal)
}

func report(line int, where string, msg string) {
    fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, msg)
}

func reportError(line int, msg string) {
    report(line, "", msg)
}

func consumeComment(source *string, current *int) {
    for ;; {
    }

}


func ScanTokens(source string) ([]Token, error) {
    tokens := []Token{}

    hasError := false
    start := 0
    current := 0
    line := 1

    
    addToken := func(t TokenType) {
        tokens = append(tokens, Token{
            type_: t,
            lexeme: source[start: current + 1],
            line: line,
        })

    }

    peek := func(c byte) (byte, bool) {
        if current + 1 < len(source) && source[current + 1] == c {
            return source[current + 1], true
        }
        return '~', false
    }

    // Conditionally step forward if the next char matches 
    match := func(c byte) bool {
        if current + 1 < len(source) && source[current + 1] == c {
            current += 1
            return true
        }
        return false
    }

    for ; current < len(source); current++ {
        start = current

        switch c := source[current]; c {
            case '(': addToken(LEFT_PAREN)
            case ')': addToken(RIGHT_PAREN)
            case '{': addToken(LEFT_BRACE)
            case '}': addToken(RIGHT_BRACE)
            case ',': addToken(COMMA)
            case '.': addToken(DOT)
            case '-': addToken(MINUS)
            case '+': addToken(PLUS)
            case ';': addToken(SEMICOLON)
            case '*': addToken(STAR)
            case '!':
                if match('=') {
                    addToken(BANG_EQUAL)
                } else {
                    addToken(BANG)
                }
            case '=':
                if match('=') {
                    addToken(EQUAL_EQUAL)
                } else {
                    addToken(EQUAL)
                }
            case '>':
                if match('=') {
                    addToken(GREATER_EQUAL)
                } else {
                    addToken(GREATER)
                }
            case '<':
                if match('=') {
                    addToken(LESS_EQUAL)
                } else {
                    addToken(LESS)
                }
            case '/':
                if match('/') {
                    // Go until ya can't go no more
                    // TODO
                } else {
                    addToken(SLASH)
                }
            default:
                hasError = true

            reportError(line, fmt.Sprintf("Unexpected character: %c", c))

        }
        
    }
    
    if hasError {
        return tokens, errors.New("Unexpected characters in source")
    }

    return tokens, nil
}
