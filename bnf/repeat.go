package bnf

// A*  -> 0 or more repeats of A
// A+  -> 1 or more repeats of A
type repeat struct {
	Node node
	Min  int // 0=*, 1=+
}

func (r *repeat) match(ctx *context, pos int) []int {
	current := []int{pos}
	var results []int

	for i := 0; ; i++ {
		if i >= r.Min {
			results = append(results, current...)
		}

		var next []int
		for _, p := range current {
			for _, m := range ctx.Match(r.Node, p) {
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
	return results
}

func (r *repeat) Expect() []string {
	return r.Node.Expect()
}
