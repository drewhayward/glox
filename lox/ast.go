package lox

// Got the trick to use hidden functions to limit
// the Node/Expr/Stmt interfaces from Eli Bendersky
// https://eli.thegreenplace.net/2018/go-and-algebraic-data-types/
// We'll give it a shot...

type Node interface {
	isNode()
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type ProgramNode struct{
    Node
    Statements []Stmt
}
func (_ ProgramNode) isNode() {}

type ExprStmt struct {
	Expr Node
}

func (_ ExprStmt) isNode()   {}
func (_ ExprStmt) stmtNode() {}

type PrintStmt struct {
	Expr Node
}

func (_ PrintStmt) isNode()   {}
func (_ PrintStmt) stmtNode() {}

type UnaryExpr struct {
	Operation TokenType
	Operand   Node
}

func (_ UnaryExpr) isNode()   {}
func (_ UnaryExpr) exprNode() {}

type GroupingExpr struct {
	Operand Node
}

func (_ GroupingExpr) isNode()   {}
func (_ GroupingExpr) exprNode() {}

type BinaryExpr struct {
	Operation TokenType
	Lhs       Node
	Rhs       Node
}

func (_ BinaryExpr) isNode()   {}
func (_ BinaryExpr) exprNode() {}

type LiteralExpr[T any] struct {
	value T
}

func (_ LiteralExpr[T]) isNode()   {}
func (_ LiteralExpr[T]) exprNode() {}

func NewLiteral[T any](val T) LiteralExpr[T] {
	return LiteralExpr[T]{val}
}
