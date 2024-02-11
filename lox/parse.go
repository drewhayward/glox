package lox

import (
	"errors"
	"fmt"
	"strconv"
)

func ParseExpr(tokens []Token) (Node, error) {
	state := parserState{tokens, 0}

	expr := state.parseExpr()

	if !state.Done() {
		return nil, errors.New("Leftover tokens after parsing")
	}

	return expr, nil
}

type parserState struct {
	tokens  []Token
	current int
}

func (ps *parserState) Done() bool {
	if ps.current < len(ps.tokens) {
		return false
	}
	return true
}

// Returns the current token
func (ps *parserState) peekToken() Token {
	return ps.tokens[ps.current]
}

// Checks the current token for a specific token type
func (ps *parserState) checkTokenType(ttype TokenType) bool {
	if ps.current >= len(ps.tokens) {
		return false
	}

	return ps.peekToken().type_ == ttype
}

// Moves the state forward
func (ps *parserState) advanceToken() {
	if ps.current < len(ps.tokens) {
		ps.current++
	}
}

func (ps *parserState) consumeToken(ttype TokenType, errorMsg string) {
	if !ps.checkTokenType(ttype) {
		panic(errorMsg)
	}
}

// Matches any of the token types and consumes it
func (ps *parserState) matchToken(types ...TokenType) bool {
	for _, ttype := range types {
		if ps.checkTokenType(ttype) {
			ps.advanceToken()
			return true
		}
	}
	return false
}

func (ps *parserState) previous() Token {
	return ps.tokens[ps.current-1]
}

func (ps *parserState) parseExpr() Expr {
	return ps.parseEquality()
}

func (ps *parserState) parseEquality() Expr {
	expr := ps.parseComparison()

	for ps.matchToken(EQUAL_EQUAL) {
		op := ps.previous().type_
		rhs := ps.parseComparison()
		expr = BinaryExpr{
			op,
			expr,
			rhs,
		}
	}

	return expr
}

func (ps *parserState) parseComparison() Expr {
	expr := ps.parseTerm()

	for ps.matchToken(LESS, LESS_EQUAL, GREATER_EQUAL, GREATER) {
		op := ps.previous().type_
		rhs := ps.parseTerm()
		expr = BinaryExpr{
			op,
			expr,
			rhs,
		}
	}

	return expr
}

func (ps *parserState) parseTerm() Expr {
	expr := ps.parseFactor()

	for ps.matchToken(PLUS, MINUS) {
		op := ps.previous().type_
		rhs := ps.parseFactor()
		expr = BinaryExpr{
			op,
			expr,
			rhs,
		}
	}

	return expr
}

func (ps *parserState) parseFactor() Expr {
	expr := ps.parseUnary()

	for ps.matchToken(SLASH, STAR) {
		op := ps.previous().type_
		rhs := ps.parseUnary()
		expr = BinaryExpr{
			op,
			expr,
			rhs,
		}
	}

	return expr
}

func (ps *parserState) parseUnary() Expr {
	if ps.matchToken(MINUS, BANG) {
		return UnaryExpr{ps.previous().type_, ps.parsePrimary()}
	}
	return ps.parsePrimary()
}

func (ps *parserState) parsePrimary() Expr {
	if ps.matchToken(FALSE) {
		return NewLiteral(false)
	}
	if ps.matchToken(TRUE) {
		return NewLiteral(true)
	}
	if ps.matchToken(NIL) {
		return NewLiteral[*struct{}](nil)
	}

	if ps.matchToken(NUMBER, STRING) {
		result, err := strconv.ParseFloat(ps.previous().lexeme, 64)
		if err != nil {
			fmt.Println("Failed to parse float")
		}
		return NewLiteral(result)
	}

	if ps.matchToken(LEFT_PAREN) {
		expr := ps.parseExpr()
		ps.consumeToken(RIGHT_PAREN, "Expect ')' after expression.")
		return GroupingExpr{expr}
	}

	panic("Fell through")
}
