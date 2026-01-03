package bnf

type optional struct {
	Node node
}

func (o *optional) match(ctx *context, pos int) ([]int, error) {
	// we can always match nothing
	results := []int{pos}

	// try to match Node
	matches, err := ctx.Match(o.Node, pos)
	if err != nil {
		return nil, err
	}
	for _, m := range matches {
		if m != pos { // protection against empty progress
			results = append(results, m)
		}
	}

	return results, nil
}

func (o *optional) Expect() []string {
	return o.Node.Expect()
}
