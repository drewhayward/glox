package lox

// Got the trick to use hidden functions to limit
// the Node/Expr/Stmt interfaces from Eli Bendersky
// https://eli.thegreenplace.net/2018/go-and-algebraic-data-types/ 

type Node interface{
    isNode()
}

type Expr interface{
    Node
    exprNode()
}

type UnaryExpr struct {
	Operation TokenType
	Operand      Node
}
func (_ UnaryExpr) isNode() {}
func (_ UnaryExpr) exprNode() {}

type GroupingExpr struct {
    Operand Node
}
func (_ GroupingExpr) isNode() {}
func (_ GroupingExpr) exprNode() {}

type BinaryExpr struct {
	Operation TokenType
	Lhs       Node
	Rhs       Node
}
func (_ BinaryExpr) isNode() {}
func (_ BinaryExpr) exprNode() {}


type LiteralExpr[T any] struct {
	value T
}
func (_ LiteralExpr[T]) isNode() {}
func (_ LiteralExpr[T]) exprNode() {}

func NewLiteral[T any](val T) LiteralExpr[T] {
	return LiteralExpr[T]{val}
}
