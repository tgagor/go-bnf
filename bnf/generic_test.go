package bnf

func testMatch(node node, input string, pos int) []int {
	ctx := NewContext(input)
	return ctx.Match(node, pos)
}
