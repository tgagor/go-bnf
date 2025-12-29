package bnf

// A*  -> 0 or more repeats of A
// A+  -> 1 or more repeats of A
type Repeat struct {
	Node Node
	Min  int // 0=*, 1=+
}

func (r *Repeat) Match(input string, pos int) []int {
	// aktualne pozycje po N powtórzeniach
	current := []int{pos}

	// wszystkie pozycje, które spełniają Min
	var results []int

	for i := 0; ; i++ {
		// jeśli osiągnęliśmy minimalną liczbę powtórzeń
		if i >= r.Min {
			results = append(results, current...)
		}

		var next []int

		for _, p := range current {
			matches := r.Node.Match(input, p)
			for _, m := range matches {
				// 🔑 WARUNEK BEZPIECZEŃSTWA
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

	if len(results) == 0 {
		return nil
	}

	return results
}
