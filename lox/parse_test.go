package lox

// func TestParseTokens(t *testing.T) {
// 	testCases := []struct {
// 		input      []Token
// 		expression Expr
// 	}{
// 		// 1
// 		{[]Token{{NUMBER, "1", "1", 0}}, Literal[string]{"1"}},
// 		// "aoeu"
// 		{[]Token{{STRING, "aoeu", "aoeu", 0}}, Literal[string]{"aoeu"}},
// 		// nil
// 		{[]Token{{NIL, "nil", "nil", 0}}, Literal[*struct{}]{nil}},
// 		// ("aoeu")
// 		{[]Token{
// 			{LEFT_PAREN, "(", "(", 0},
// 			{STRING, "aoeu", "aoeu", 0},
// 			{RIGHT_PAREN, ")", ")", 0},
// 		}, Literal[string]{"aoeu"}},
// 		// (1)
// 		{[]Token{
// 			{LEFT_PAREN, "(", "(", 0},
// 			{NUMBER, "1", "1", 0},
// 			{RIGHT_PAREN, ")", ")", 0},
// 		}, Literal[string]{"1"}},
// 		// -1
// 		{[]Token{
// 			{MINUS, "-", "-", 0},
// 			{NUMBER, "1", "1", 0},
// 		}, UnaryExpr{MINUS, Literal[string]{"1"}}},
// 		// !1
// 		{[]Token{
// 			{BANG, "!", "!", 0},
// 			{NUMBER, "1", "1", 0},
// 		}, UnaryExpr{BANG, Literal[string]{"1"}}},
// 		// 1 * 1
// 		{[]Token{
// 			{NUMBER, "1", "1", 0},
// 			{STAR, "*", "*", 0},
// 			{NUMBER, "1", "1", 0},
// 		}, BinaryExpr{STAR, Literal[string]{"1"}, Literal[string]{"1"}}},
// 		// 1 / 1
// 		{[]Token{
// 			{NUMBER, "1", "1", 0},
// 			{SLASH, "/", "/", 0},
// 			{NUMBER, "1", "1", 0},
// 		}, BinaryExpr{SLASH, Literal[string]{"1"}, Literal[string]{"1"}}},
// 		// 1 + 1
// 		{[]Token{
// 			{NUMBER, "1", "1", 0},
// 			{PLUS, "+", "+", 0},
// 			{NUMBER, "1", "1", 0},
// 		}, BinaryExpr{PLUS, Literal[string]{"1"}, Literal[string]{"1"}}},
// 		// 1 - 1
// 		{[]Token{
// 			{NUMBER, "1", "1", 0},
// 			{MINUS, "-", "-", 0},
// 			{NUMBER, "1", "1", 0},
// 		}, BinaryExpr{MINUS, Literal[string]{"1"}, Literal[string]{"1"}}},
// 		// 1 < 1
// 		{[]Token{
// 			{NUMBER, "1", "1", 0},
// 			{LESS, "*", "*", 0},
// 			{NUMBER, "1", "1", 0},
// 		}, BinaryExpr{LESS, Literal[string]{"1"}, Literal[string]{"1"}}},
// 		// 1 <= 1
// 		{[]Token{
// 			{NUMBER, "1", "1", 0},
// 			{LESS_EQUAL, "<=", "<=", 0},
// 			{NUMBER, "1", "1", 0},
// 		}, BinaryExpr{LESS_EQUAL, Literal[string]{"1"}, Literal[string]{"1"}}},
// 		// 1 > 1
// 		{[]Token{
// 			{NUMBER, "1", "1", 0},
// 			{GREATER, ">", ">", 0},
// 			{NUMBER, "1", "1", 0},
// 		}, BinaryExpr{GREATER, Literal[string]{"1"}, Literal[string]{"1"}}},
// 		// 1 >= 1
// 		{[]Token{
// 			{NUMBER, "1", "1", 0},
// 			{GREATER_EQUAL, ">=", ">=", 0},
// 			{NUMBER, "1", "1", 0},
// 		}, BinaryExpr{GREATER_EQUAL, Literal[string]{"1"}, Literal[string]{"1"}}},
// 		// 1 == 1
// 		{[]Token{
// 			{NUMBER, "1", "1", 0},
// 			{EQUAL_EQUAL, "==", "==", 0},
// 			{NUMBER, "1", "1", 0},
// 		}, BinaryExpr{EQUAL_EQUAL, Literal[string]{"1"}, Literal[string]{"1"}}},
// 		// 1 + 1 + 1
// 		{[]Token{
// 			{NUMBER, "1", "1", 0},
// 			{PLUS, "+", "+", 0},
// 			{NUMBER, "1", "1", 0},
// 			{PLUS, "+", "+", 0},
// 			{NUMBER, "1", "1", 0},
// 		}, BinaryExpr{PLUS, BinaryExpr{PLUS, Literal[string]{"1"}, Literal[string]{"1"}}, Literal[string]{"1"}}},
// 	}
//
// 	for _, tc := range testCases {
// 		t.Run(fmt.Sprintf("ParseTokens(%q)", tc.input), func(t *testing.T) {
// 			result := ParseExpr(tc.input)
// 			if result != tc.expression {
// 				t.Errorf("ParseTokens(%q) returned the wrong ast. Got %q, wanted: %q", tc.input, result, tc.expression)
// 			}
// 		})
// 	}
// }
