package bnf

type Sequence struct {
	Elements []Node
}

// func (s *Sequence) Match(input string, pos int) []int {
// 	// aktual possible positions
// 	positions := []int{pos}

// 	for _, elem := range s.Elements {
// 		var next []int

// 		for _, p := range positions {
// 			matches := elem.Match(input, p)
// 			if len(matches) > 0 {
// 				next = append(next, matches...)
// 			}
// 		}

// 		// if no variant match, it's not this sequence
// 		if len(next) == 0 {
// 			return nil
// 		}

// 		positions = next
// 	}

// 	return positions
// }

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
