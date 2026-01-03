package bnf

// A*  -> 0 or more repeats of A
// A+  -> 1 or more repeats of A
type repeat struct {
	Node node
	Min  int // 0=*, 1=+
}

func (r *repeat) match(ctx *context, pos int) ([]MatchResult, error) {
	// Start with empty match
	currentResults := []MatchResult{{End: pos, Nodes: nil}}
	var finalResults []MatchResult

	for i := 0; ; i++ {
		if i >= r.Min {
			// Current state is valid, add to final results
			finalResults = append(finalResults, currentResults...)
		}

		var nextResults []MatchResult
		for _, res := range currentResults {
			// Try to match more
			matches, err := ctx.Match(r.Node, res.End)
			if err != nil {
				return nil, err
			}
			for _, m := range matches {
				// Prevent infinite loops on empty matches
				if m.End > res.End {
					// Combine nodes
					newNodes := make([]*ASTNode, len(res.Nodes)+len(m.Nodes))
					copy(newNodes, res.Nodes)
					copy(newNodes[len(res.Nodes):], m.Nodes)

					nextResults = append(nextResults, MatchResult{
						End:   m.End,
						Nodes: newNodes,
					})
				}
			}
		}

		if len(nextResults) == 0 {
			break
		}
		currentResults = nextResults
	}
	return finalResults, nil
}

func (r *repeat) Expect() []string {
	return r.Node.Expect()
}
