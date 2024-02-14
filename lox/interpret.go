package lox 

type Value interface {}

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

func Interpret(node Stmt) error {
    return nil
}

func Evaluate(node Node) Value {
    switch nt := node.(type) {

    case LiteralExpr[bool]:
        return nt.value
    case LiteralExpr[float64]:
        return nt.value
    case LiteralExpr[string]:
        return nt.value
    case LiteralExpr[*struct{}]:
        return Null(nil)
    case UnaryExpr:
        value := Evaluate(nt)
        if nt.Operation == BANG {
            // Assert that the value is a bool?
            return isTruthy(value)
        }
        if nt.Operation == MINUS {
            v := value.(float64)
            return float64(v)
            
        }
        panic("Bad operand in UnaryExpr")
    case GroupingExpr:
        return Evaluate(nt.Operand)
    case BinaryExpr:
        lhs := Evaluate(nt.Lhs)
        rhs := Evaluate(nt.Rhs)

        // Handle equals first since it may not have numbers
        if nt.Operation == BANG_EQUAL {
            return !isEqual(lhs, rhs)
        } 
        if nt.Operation == EQUAL_EQUAL {
            return isEqual(lhs, rhs)
        } 

        // All the other operations need numbers
        nl, nr := lhs.(float64), rhs.(float64)
        switch nt.Operation {
        case STAR:
            return nl * nr
        case SLASH:
            return nl / nr
        case PLUS:
            return nl + nr
        case MINUS:
            return nl - nr
        case GREATER:
            return nl > nr
        case GREATER_EQUAL:
            return nl >= nr
        case LESS:
            return nl < nr
        case LESS_EQUAL:
            return nl <= nr
        }
        panic("Bad operand in BinaryExpr")
        
    }

    panic("Ahhhh")
}
