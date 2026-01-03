package bnf

type sequence struct {
	Elements []node
}

func (s *sequence) match(ctx *context, pos int) ([]MatchResult, error) {
	// Start with one "empty" result at current position
	currentResults := []MatchResult{{End: pos, Nodes: nil}}

	for _, elem := range s.Elements {
		var nextResults []MatchResult
		for _, res := range currentResults {
			// Try to match element at current result's end position
			matches, err := ctx.Match(elem, res.End)
			if err != nil {
				return nil, err
			}
			for _, m := range matches {
				// Combine existing nodes with new nodes
				newNodes := make([]*ASTNode, len(res.Nodes)+len(m.Nodes))
				copy(newNodes, res.Nodes)
				copy(newNodes[len(res.Nodes):], m.Nodes)

				nextResults = append(nextResults, MatchResult{
					End:   m.End,
					Nodes: newNodes,
				})
			}
		}
		if len(nextResults) == 0 {
			return nil, nil // No path continued
		}
		currentResults = nextResults
	}
	// TODO: Prune duplicate results? Ambiguity might produce same end with different trees.
	// For now kept all.
	return currentResults, nil
}

func (s *sequence) Expect() []string {
	if len(s.Elements) == 0 {
		return nil
	}
	return s.Elements[0].Expect()
}
