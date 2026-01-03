package bnf

// A*  -> 0 or more repeats of A
// A+  -> 1 or more repeats of A
type repeat struct {
	Node node
	Min  int // 0=*, 1=+
}

func (r *repeat) match(ctx *context, pos int) ([]int, error) {
	current := []int{pos}
	var results []int

	for i := 0; ; i++ {
		if i >= r.Min {
			results = append(results, current...)
		}

		var next []int
		for _, p := range current {
			matches, err := ctx.Match(r.Node, p)
			if err != nil {
				return nil, err
			}
			for _, m := range matches {
				// safety check to prevent infinite loops
				if m > p {
					next = append(next, m)
				}
			}
		}

		if len(next) == 0 {
			break
		}
		current = next
	}
	return results, nil
}

func (r *repeat) Expect() []string {
	return r.Node.Expect()
}
