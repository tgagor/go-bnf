package bnf

type sequence struct {
	Elements []node
}

func (s *sequence) match(ctx *context, pos int) ([]int, error) {
	positions := []int{pos}

	for _, elem := range s.Elements {
		var next []int
		for _, p := range positions {
			matches, err := ctx.Match(elem, p)
			if err != nil {
				return nil, err
			}
			next = append(next, matches...)
		}
		if len(next) == 0 {
			return nil, nil
		}
		positions = next
	}
	return positions, nil
}

func (s *sequence) Expect() []string {
	if len(s.Elements) == 0 {
		return nil
	}
	return s.Elements[0].Expect()
}
