package bnf

type optional struct {
	Node node
}

func (o *optional) match(ctx *context, pos int) ([]MatchResult, error) {
	// we can always match nothing
	results := []MatchResult{{End: pos, Nodes: nil}}

	// try to match Node
	matches, err := ctx.Match(o.Node, pos)
	if err != nil {
		return nil, err
	}
	for _, m := range matches {
		if m.End != pos { // protection against empty progress, although optional usually allows it?
			// Actually optional *is* the empty match. If node matches something, we return it.
			results = append(results, m)
		}
	}

	return results, nil
}

func (o *optional) Expect() []string {
	return o.Node.Expect()
}
