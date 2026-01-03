package bnf

func testMatch(node node, input string, pos int) []int {
	ctx := NewContext(input)
	res, _ := ctx.Match(node, pos)
	var out []int
	for _, r := range res {
		out = append(out, r.End)
	}
	return out
}
