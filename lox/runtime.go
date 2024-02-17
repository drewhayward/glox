package lox

import "fmt"

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
	Env map[string]any
}

func NewRuntimeState() RuntimeState {
	return RuntimeState{
		Env: make(map[string]any),
	}
}

func (rs *RuntimeState) Run(source string) {
	// Tokenize the source string
	tokens, err := ScanTokens(source)
	if err != nil {
		fmt.Printf("Lexing Error %v\n", tokens)
		return
	}

	root, err := Parse(tokens)
	if err != nil {
		fmt.Println(err)
		return
	}

	pStmts := root.(ProgramNode)
	for _, stmt := range pStmts.Statements {
		err := rs.Interpret(stmt)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (state *RuntimeState) Interpret(stmt Stmt) error {
	switch stype := stmt.(type) {
	case PrintStmt:
		value, err := state.Evaluate(stype.Expr)
		if err != nil {
			return err
		}

		fmt.Println(value)
	case ExprStmt:
		// We don't actually do anything with an ExprStmt value
		_, err := state.Evaluate(stype.Expr)
		if err != nil {
			return err
		}

	case DeclarationStmt:
		_, exists := state.Env[stype.Name]
		if exists {
			return RuntimeError{message: fmt.Sprintf("Var %s has already been declared", stype.Name)}
		}

		var init any
		if stype.Expr != nil {
			v, err := state.Evaluate(*stype.Expr)
			if err != nil {
				return err
			}
			init = v
		}

		state.Env[stype.Name] = init
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
		value, ok := rs.Env[nt.Name]
		if !ok {
			return nil, RuntimeError{message: fmt.Sprintf("Var %s has never been declared", nt.Name)}
		}

		return value, nil
	case AssignExpr:
		_, ok := rs.Env[nt.Name]
		if !ok {
			return nil, RuntimeError{message: fmt.Sprintf("Var %s has never been declared", nt.Name)}
		}
		value, err := rs.Evaluate(nt.Value)
		if err != nil {
			return nil, err
		}

		rs.Env[nt.Name] = value
		return value, nil
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

	}

	return nil, RuntimeError{
		message: fmt.Sprintf("Attempting to evaluate an unsupported type %+v", node),
	}
}
