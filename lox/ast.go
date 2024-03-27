package lox

// Got the trick to use hidden functions to limit
// the Node/Expr/Stmt interfaces from Eli Bendersky
// https://eli.thegreenplace.net/2018/go-and-algebraic-data-types/
// We'll give it a shot...

type Node interface {
	isNode()
}

type Expr interface {
	exprNode()
}

type Stmt interface {
	stmtNode()
}

type ProgramNode struct {
	Statements []Stmt
}

func (_ ProgramNode) isNode() {}

type ExprStmt struct {
	Expr Expr
}

func (_ ExprStmt) isNode()   {}
func (_ ExprStmt) stmtNode() {}

type PrintStmt struct {
	Expr Expr
}

func (_ PrintStmt) isNode()   {}
func (_ PrintStmt) stmtNode() {}

type BlockStmt struct {
	Statements []Stmt
}

func (_ BlockStmt) isNode()   {}
func (_ BlockStmt) stmtNode() {}

type DeclarationStmt struct {
	Name string
	Expr *Expr
}

func (_ DeclarationStmt) isNode()   {}
func (_ DeclarationStmt) stmtNode() {}

type FunctionDeclarationStmt struct {
	Name       string
	Parameters []string
	Body       BlockStmt
}

func (_ FunctionDeclarationStmt) isNode()   {}
func (_ FunctionDeclarationStmt) stmtNode() {}

type ClassDeclarationStmt struct {
	Name      string
	Functions []FunctionDeclarationStmt
}

func (_ ClassDeclarationStmt) isNode()   {}
func (_ ClassDeclarationStmt) stmtNode() {}

type ReturnStmt struct {
	Value Expr
}

func (_ ReturnStmt) isNode()   {}
func (_ ReturnStmt) stmtNode() {}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (_ IfStmt) isNode()   {}
func (_ IfStmt) stmtNode() {}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (_ WhileStmt) isNode()   {}
func (_ WhileStmt) stmtNode() {}

type CallExpr struct {
	Callee Expr
	Args   []Expr
}

func (_ CallExpr) isNode()   {}
func (_ CallExpr) exprNode() {}

type UnaryExpr struct {
	Operation TokenType
	Operand   Expr
}

func (_ UnaryExpr) isNode()   {}
func (_ UnaryExpr) exprNode() {}

type GroupingExpr struct {
	Operand Expr
}

func (_ GroupingExpr) isNode()   {}
func (_ GroupingExpr) exprNode() {}

type BinaryExpr struct {
	Operation TokenType
	Lhs       Expr
	Rhs       Expr
}

func (_ BinaryExpr) isNode()   {}
func (_ BinaryExpr) exprNode() {}

type VarExpr struct {
	Name string
}

func (_ VarExpr) isNode()   {}
func (_ VarExpr) exprNode() {}

type LogicalExpr struct {
	Operation TokenType
	Lhs       Expr
	Rhs       Expr
}

func (_ LogicalExpr) isNode()   {}
func (_ LogicalExpr) exprNode() {}

type AssignExpr struct {
	Name  string
	Value Expr
}

func (_ AssignExpr) isNode()   {}
func (_ AssignExpr) exprNode() {}

type LiteralExpr[T any] struct {
	value T
}

func (_ LiteralExpr[T]) isNode()   {}
func (_ LiteralExpr[T]) exprNode() {}

func NewLiteralExpr[T any](val T) LiteralExpr[T] {
	return LiteralExpr[T]{val}
}
