package lox

// AST struct types
type (
	Equality   struct{}
	Comparison struct{}
	Term       struct{}
	Factor     struct{}
	Unary      struct{}
	Primary    struct{}
)

type Expr interface{}

type UnaryExpr struct {
	operation TokenType
	expr      Expr
}

type BinaryExpr struct {
	operation TokenType
	lhs       Expr
	rhs       Expr
}

type Literal[T any] struct {
	value T
}

func NewLiteral[T any](val T) Literal[T] {
	return Literal[T]{val}
}
