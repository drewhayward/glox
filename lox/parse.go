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
		s, err := ps.parseDeclaration()
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, s)
	}

	return ProgramNode{Statements: stmts}, nil
}

func (ps *parserState) parseDeclaration() (Stmt, error) {
	if ps.matchToken(VAR) {
		err := ps.consumeToken(IDENTIFIER, "Expected identifier.")
		if err != nil {
			return nil, err
		}
		d := DeclarationStmt{Name: ps.previous().lexeme}

		// Optionally consume the value definition
		if ps.matchToken(EQUAL) {
			expr, err := ps.parseExpr()
			if err != nil {
				return nil, err
			}

			d.Expr = &expr
		}

		err = ps.consumeToken(SEMICOLON, "Expected semicolon.")
		if err != nil {
			return nil, err
		}

		return d, nil
	}

	return ps.parseStmt()
}

func (ps *parserState) parseStmt() (Stmt, error) {
	switch ps.peekToken().type_ {
	case PRINT:
		return ps.parsePrint()
	case WHILE:
		return ps.parseWhile()
	case FOR:
		return ps.parseFor()
	case LEFT_BRACE:
		return ps.parseBlock()
	case IF:
		return ps.parseIf()
	default:
		return ps.parseExprStmt()
	}
}

func (ps *parserState) parseExprStmt() (Stmt, error) {
	expr, err := ps.parseExpr()
	if err != nil {
		return nil, err
	}

	err = ps.consumeToken(SEMICOLON, "Expected semicolon.")
	if err != nil {
		return nil, err
	}

	return ExprStmt{expr}, nil
}

func (ps *parserState) parseIf() (Stmt, error) {
	err := ps.consumeToken(IF, "Expected 'if' keyword to start if statement")
	if err != nil {
		return nil, err
	}

	err = ps.consumeToken(LEFT_PAREN, "Expected '(' after if ")
	if err != nil {
		return nil, err
	}

	condition, err := ps.parseExpr()
	if err != nil {
		return nil, err
	}

	err = ps.consumeToken(RIGHT_PAREN, "Expected ')' after if ")
	if err != nil {
		return nil, err
	}

	thenStmt, err := ps.parseStmt()
	if err != nil {
		return nil, err
	}

	var elseStmt Stmt
	if ps.matchToken(ELSE) {
		elseStmt, err = ps.parseStmt()
		if err != nil {
			return nil, err
		}
	}

	return IfStmt{Condition: condition, ThenBranch: thenStmt, ElseBranch: elseStmt}, nil
}

func (ps *parserState) parseWhile() (Stmt, error) {
	err := ps.consumeToken(WHILE, "Expected 'while' to start")
	if err != nil {
		return nil, err
	}

	err = ps.consumeToken(LEFT_PAREN, "Expected '(' to start")
	if err != nil {
		return nil, err
	}

	expr, err := ps.parseExpr()
	if err != nil {
		return nil, err
	}

	err = ps.consumeToken(RIGHT_PAREN, "Expected ')' to start")
	if err != nil {
		return nil, err
	}

	stmt, err := ps.parseStmt()
	if err != nil {
		return nil, err
	}

	return WhileStmt{
		Condition: expr,
		Body:      stmt,
	}, nil
}

// Parse a for loop as a desugared while because we can
func (ps *parserState) parseFor() (Stmt, error) {
	err := ps.consumeToken(FOR, "Expected 'for' to start loop")
	if err != nil {
		return nil, err
	}

	err = ps.consumeToken(LEFT_PAREN, "Expected '(' to start")
	if err != nil {
		return nil, err
	}

	// Pull out the init
	var init Stmt
	if ps.peekToken().type_ == VAR {
		init, err = ps.parseDeclaration()
		if err != nil {
			return nil, err
		}
	} else if !ps.matchToken(SEMICOLON) {
		init, err = ps.parseExprStmt()
		if err != nil {
			return nil, err
		}

		err = ps.consumeToken(SEMICOLON, "Expected ';' after loop init")
		if err != nil {
			return nil, err
		}
	}

	var cond Expr
	if !ps.matchToken(SEMICOLON) {
		cond, err = ps.parseExpr()
		if err != nil {
			return nil, err
		}

		err = ps.consumeToken(SEMICOLON, "Expected ';' after loop condition")
		if err != nil {
			return nil, err
		}
	}

	var increment Expr
	if !ps.matchToken(RIGHT_PAREN) {
		increment, err = ps.parseExpr()
		if err != nil {
			return nil, err
		}

		err = ps.consumeToken(RIGHT_PAREN, "Expected ')' to start")
		if err != nil {
			return nil, err
		}
	}

	body, err := ps.parseStmt()
	if err != nil {
		return nil, err
	}

	// Add the increment to the end of the while
	if increment != nil {
		body = BlockStmt{
			Statements: []Stmt{body, ExprStmt{Expr: increment}},
		}
	}

	// Create a new while with the condition
	if cond == nil {
		cond = LiteralExpr[bool]{value: true}
	}
	body = WhileStmt{
		Condition: cond,
		Body:      body,
	}

	if init != nil {
		body = BlockStmt{Statements: []Stmt{init, body}}
	}

	return body, nil
}

func (ps *parserState) parsePrint() (Stmt, error) {
	err := ps.consumeToken(PRINT, "Expected 'print'")
	if err != nil {
		return nil, err
	}
	expr, err := ps.parseExpr()
	if err != nil {
		return nil, err
	}

	err = ps.consumeToken(SEMICOLON, "Expected semicolon.")
	if err != nil {
		return nil, err
	}

	return PrintStmt{expr}, nil
}

func (ps *parserState) parseBlock() (Stmt, error) {
	err := ps.consumeToken(LEFT_BRACE, "Expected block to start with '{'")
	if err != nil {
		return nil, err
	}

	stmts := make([]Stmt, 0)
	for !ps.matchToken(RIGHT_BRACE) {
		s, err := ps.parseDeclaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, s)
	}

	return BlockStmt{Statements: stmts}, nil
}

func (ps *parserState) parseExpr() (Expr, error) {
	expr, err := ps.parseAssignment()
	if err != nil {
		return nil, err
	}

	return expr, err
}

func (ps *parserState) parseAssignment() (Expr, error) {
	// This parses the lefthand side of the assignment
	expr, err := ps.parseOr()
	if err != nil {
		return nil, err
	}

	if ps.matchToken(EQUAL) {
		equal := ps.previous()
		// Since this is right-associative, we use recursion to handle
		// further assignment expressions
		value, err := ps.parseAssignment()
		if err != nil {
			return nil, err
		}

		// If the LHS is a variable, we can assign to it
		v, ok := expr.(VarExpr)
		if ok {
			return AssignExpr{
				Name:  v.Name,
				Value: value,
			}, nil

		}
		return nil, ParseError{
			message: "Invalid assignment target",
			token:   equal,
		}
	}

	return expr, nil
}

func (ps *parserState) parseOr() (Expr, error) {
	expr, err := ps.parseAnd()
	if err != nil {
		return nil, err
	}

	for ps.matchToken(OR) {
		fmt.Println("parsing or")
		tok := ps.previous()
		rhs, err := ps.parseAnd()
		if err != nil {
			return nil, err
		}

		expr = LogicalExpr{
			Operation: tok.type_,
			Lhs:       expr,
			Rhs:       rhs,
		}
	}

	return expr, nil
}
func (ps *parserState) parseAnd() (Expr, error) {
	expr, err := ps.parseEquality()
	if err != nil {
		return nil, err
	}

	for ps.matchToken(AND) {
		tok := ps.previous()
		rhs, err := ps.parseEquality()
		if err != nil {
			return nil, err
		}

		expr = LogicalExpr{
			Operation: tok.type_,
			Lhs:       expr,
			Rhs:       rhs,
		}
	}

	return expr, nil
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
		return NewLiteralExpr(false), nil
	}
	if ps.matchToken(TRUE) {
		return NewLiteralExpr(true), nil
	}
	if ps.matchToken(NIL) {
		return NewLiteralExpr[*struct{}](nil), nil
	}

	if ps.matchToken(IDENTIFIER) {
		return VarExpr{Name: ps.previous().lexeme}, nil
	}

	if ps.matchToken(NUMBER) {
		result, err := strconv.ParseFloat(ps.previous().lexeme, 64)
		if err != nil {
			return nil, ParseError{
				message: "Expected float",
			}
		}

		return NewLiteralExpr(result), nil
	}

	if ps.matchToken(STRING) {
		return NewLiteralExpr(ps.previous().lexeme), nil
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
