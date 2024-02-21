package lox

import (
	"fmt"
	"io"
	"os"
)

type Value interface{}

type Null *struct{}

// Determines whether a value is truthy.
// Lox implements Ruby's truthiness rules
// lox-nil and false are falsey
func isTruthy(v any) bool {
	_, ok := v.(Null)
	if ok || v == false {
		return false
	}

	return true
}

func isEqual(lhs any, rhs any) bool {
	_, lnull := lhs.(Null)
	_, rnull := rhs.(Null)
	if lnull && rnull { // nil equals nil
		return true
	}
	if lnull { // nil doesn't equal other values
		return false
	}

	return lhs == rhs
}

type RuntimeError struct {
	message string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError: %s", e.message)
}

type RuntimeState struct {
	// Points to the currently active scope for execution
	CurrEnv   *ScopeEnv
	OutWriter io.Writer
}

func NewRuntimeState() RuntimeState {
	return RuntimeState{
		CurrEnv:   NewScopeEnv(nil),
		OutWriter: os.Stdout,
	}
}

func (rs *RuntimeState) Run(source string) {
	// Tokenize the source string
	tokens, err := ScanTokens(source)
	if err != nil {
		fmt.Fprintf(rs.OutWriter, "Lexing Error %v\n", tokens)
		return
	}

	root, err := Parse(tokens)
	if err != nil {
		fmt.Fprintln(rs.OutWriter, err)
		return
	}

	pStmts := root.(ProgramNode)
	for _, stmt := range pStmts.Statements {
		err := rs.Interpret(stmt)
		if err != nil {
			fmt.Fprintln(rs.OutWriter, err.Error())
		}
	}
}

func (rs *RuntimeState) Interpret(stmt Stmt) error {
	switch stype := stmt.(type) {
	case PrintStmt:
		value, err := rs.Evaluate(stype.Expr)
		if err != nil {
			return err
		}

		fmt.Fprintln(rs.OutWriter, value)
	case ExprStmt:
		// We don't actually do anything with an ExprStmt value
		_, err := rs.Evaluate(stype.Expr)
		if err != nil {
			return err
		}

	case DeclarationStmt:
		var init any
		if stype.Expr != nil {
			v, err := rs.Evaluate(*stype.Expr)
			if err != nil {
				return err
			}
			init = v
		}

		rs.CurrEnv.Declare(stype.Name, init)
	case BlockStmt:
		// Create a new variable scope
		rs.CurrEnv = NewScopeEnv(rs.CurrEnv)
		for _, stmt := range stype.Statements {
			rs.Interpret(stmt)
		}
		rs.CurrEnv = rs.CurrEnv.parent
	case IfStmt:
		cond, err := rs.Evaluate(stype.Condition)
		if err != nil {
			return err
		}

		if isTruthy(cond) {
			err = rs.Interpret(stype.ThenBranch)
		} else {
			err = rs.Interpret(stype.ElseBranch)
		}

		if err != nil {
			return err
		}
	case WhileStmt:
		for {
			cond, err := rs.Evaluate(stype.Condition)
			if err != nil {
				return err
			}

			if !isTruthy(cond) {
				break
			}

			err = rs.Interpret(stype.Body)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (rs *RuntimeState) Evaluate(node Expr) (Value, error) {
	switch nt := node.(type) {
	case LiteralExpr[bool]:
		return nt.value, nil
	case LiteralExpr[float64]:
		return nt.value, nil
	case LiteralExpr[string]:
		return nt.value, nil
	case LiteralExpr[*struct{}]:
		return Null(nil), nil
	case VarExpr:
		return rs.CurrEnv.Lookup(nt.Name)
	case AssignExpr:
		v, err := rs.Evaluate(nt.Value)
		if err != nil {
			return nil, err
		}

		return rs.CurrEnv.Assign(nt.Name, v)
	case UnaryExpr:
		value, err := rs.Evaluate(nt)
		if err != nil {
			return nil, err
		}

		if nt.Operation == BANG {
			// Assert that the value is a bool?
			return isTruthy(value), nil
		}
		if nt.Operation == MINUS {
			v := value.(float64)
			return float64(v), nil

		}

		return nil, RuntimeError{
			message: fmt.Sprintf("Bad operand '%s' in unary expression", nt.Operation),
		}
	case GroupingExpr:
		return rs.Evaluate(nt.Operand)
	case BinaryExpr:
		lhs, err := rs.Evaluate(nt.Lhs)
		if err != nil {
			return nil, err
		}
		rhs, err := rs.Evaluate(nt.Rhs)
		if err != nil {
			return nil, err
		}

		// Handle equals first since it may not have numbers
		if nt.Operation == BANG_EQUAL {
			return !isEqual(lhs, rhs), nil
		}
		if nt.Operation == EQUAL_EQUAL {
			return isEqual(lhs, rhs), nil
		}

		// All the other operations need numbers
		nl, nr := lhs.(float64), rhs.(float64)
		switch nt.Operation {
		case STAR:
			return nl * nr, nil
		case SLASH:
			if nr == 0 {
				return nil, RuntimeError{message: "Division by zero"}
			}
			return nl / nr, nil
		case PLUS:
			return nl + nr, nil
		case MINUS:
			return nl - nr, nil
		case GREATER:
			return nl > nr, nil
		case GREATER_EQUAL:
			return nl >= nr, nil
		case LESS:
			return nl < nr, nil
		case LESS_EQUAL:
			return nl <= nr, nil
		}

		return nil, RuntimeError{
			message: fmt.Sprintf("Bad operand '%s' in unary expression", nt.Operation),
		}
	case LogicalExpr:
		left, err := rs.Evaluate(nt.Lhs)
		if err != nil {
			return nil, err
		}

		// Short circuit
		if nt.Operation == OR && isTruthy(left) {
			return isTruthy(left), nil
		} else if nt.Operation == AND && !isTruthy(left) {
			return isTruthy(left), nil
		}

		right, err := rs.Evaluate(nt.Rhs)
		if err != nil {
			return nil, err
		}

		if nt.Operation == OR {
			return isTruthy(left) || isTruthy(right), nil
		}
		return isTruthy(left) && isTruthy(right), nil
	}

	return nil, RuntimeError{
		message: fmt.Sprintf("Attempting to evaluate an unsupported type %+v", node),
	}
}
