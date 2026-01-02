package bnf

type optional struct {
	Node node
}

func (o *optional) match(ctx *context, pos int) []int {
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

func (o *optional) Expect() []string {
	return o.Node.Expect()
}
