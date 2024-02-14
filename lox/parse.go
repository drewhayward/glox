package lox

import (
	"fmt"
	"strconv"
)

type ParseError struct {
	message string
	line    int
	token   Token
}

func (p ParseError) Error() string {
	return fmt.Sprintf("Parse Error: %s", p.message)
}

// Parse parses the tokens returned by the lexer into an AST.
func Parse(tokens []Token) (Node, error) {
	state := parserState{tokens, 0}

	expr, err := state.parseProgram()
	if err != nil {
		return nil, err
	}

	if !state.Done() {
		return nil, ParseError{message: "Leftover tokens after parsing"}
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

func (ps *parserState) consumeToken(ttype TokenType, errorMsg string) error {
	if !ps.checkTokenType(ttype) {
		return ParseError{
			message: errorMsg,
		}
	}

	ps.advanceToken()
	return nil
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

// Production Rule Functions

func (ps *parserState) parseProgram() (Node, error) {
	stmts := make([]Stmt, 0)
	for !ps.Done() {
		s, err := ps.parseStmt()
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, s)
	}

	return ProgramNode{Statements: stmts}, nil
}

func (ps *parserState) parseStmt() (Stmt, error) {
	isPrint := ps.matchToken(PRINT)
	expr, err := ps.parseExpr()
	if err != nil {
		return nil, err
	}

	err = ps.consumeToken(SEMICOLON, "Expected semicolon.")
	if err != nil {
		return nil, err
	}

	if isPrint {
		return PrintStmt{expr}, nil
	}

	return ExprStmt{expr}, nil
}

func (ps *parserState) parseExpr() (Expr, error) {
	expr, err := ps.parseEquality()
	if err != nil {
		return nil, err
	}

	return expr, err
}

func (ps *parserState) parseEquality() (Expr, error) {
	expr, err := ps.parseComparison()
	if err != nil {
		return nil, err
	}

	for ps.matchToken(EQUAL_EQUAL) {
		op := ps.previous().type_
		rhs, err := ps.parseComparison()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{
			op,
			expr,
			rhs,
		}
	}

	return expr, nil
}

func (ps *parserState) parseComparison() (Expr, error) {
	expr, err := ps.parseTerm()
	if err != nil {
		return nil, err
	}

	for ps.matchToken(LESS, LESS_EQUAL, GREATER_EQUAL, GREATER) {
		op := ps.previous().type_
		rhs, err := ps.parseTerm()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			op,
			expr,
			rhs,
		}
	}

	return expr, nil
}

func (ps *parserState) parseTerm() (Expr, error) {
	expr, err := ps.parseFactor()
	if err != nil {
		return nil, err
	}

	for ps.matchToken(PLUS, MINUS) {
		op := ps.previous().type_
		rhs, err := ps.parseFactor()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{
			op,
			expr,
			rhs,
		}
	}

	return expr, nil
}

func (ps *parserState) parseFactor() (Expr, error) {
	expr, err := ps.parseUnary()
	if err != nil {
		return nil, err
	}

	for ps.matchToken(SLASH, STAR) {
		op := ps.previous().type_
		rhs, err := ps.parseUnary()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{
			op,
			expr,
			rhs,
		}
	}

	return expr, nil
}

func (ps *parserState) parseUnary() (Expr, error) {
	if ps.matchToken(MINUS, BANG) {
		expr, err := ps.parsePrimary()
		if err != nil {
			return nil, err
		}

		return UnaryExpr{ps.previous().type_, expr}, nil
	}

	expr, err := ps.parsePrimary()
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (ps *parserState) parsePrimary() (Expr, error) {
	if ps.matchToken(FALSE) {
		return NewLiteral(false), nil
	}
	if ps.matchToken(TRUE) {
		return NewLiteral(true), nil
	}
	if ps.matchToken(NIL) {
		return NewLiteral[*struct{}](nil), nil
	}

	if ps.matchToken(NUMBER, STRING) {
		result, err := strconv.ParseFloat(ps.previous().lexeme, 64)
		if err != nil {
			return nil, ParseError{
				message: "Expected float",
			}
		}

		return NewLiteral(result), nil
	}

	if ps.matchToken(LEFT_PAREN) {
		expr, err := ps.parseExpr()
		if err != nil {
			return nil, err
		}

		err = ps.consumeToken(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return GroupingExpr{expr}, nil
	}

	return nil, ParseError{message: "Couldn't parse expression"}
}
