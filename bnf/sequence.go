package bnf

type Sequence struct {
	Elements []Node
}

func (s *Sequence) match(ctx *Context, pos int) []int {
	positions := []int{pos}

	for _, elem := range s.Elements {
		var next []int
		for _, p := range positions {
			matches := ctx.Match(elem, p)
			next = append(next, matches...)
		}
		if len(next) == 0 {
			return nil
		}
		positions = next
	}
	return positions
}

func (s *Sequence) Expect() []string {
    if len(s.Elements) == 0 {
        return nil
    }
    return s.Elements[0].Expect()
}
