package bnf

func testMatch(node node, input string, pos int) []int {
	ctx := NewContext(input)
	res, _ := ctx.Match(node, pos)
	return res
}
