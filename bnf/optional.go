package bnf

type Optional struct {
	Node Node
}

func (o *Optional) match(ctx *Context, pos int) []int {
	// we can always match nothing
	results := []int{pos}

	// try to match Node
	matches := ctx.Match(o.Node, pos)
	for _, m := range matches {
		if m != pos { // protection against empty progress
			results = append(results, m)
		}
	}

	return results
}
