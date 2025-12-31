package bnf

import "fmt"

type memoKey struct {
	node Node
	pos  int
}

type memoEntry struct {
	results    []int // remember results`
	inProgress bool  // detect left sided recurency
}

type Context struct {
	input string
	memo  map[memoKey]*memoEntry // cache (node, pos)

	// debug
	FarthestPos int // farthest position reached during parsing
	CurrentPos int // current position during parsing (deepest, even for Sequence, Choice, Repeat)
	Expected    []string
	error       *ParseError

	// call stack
	stack []string
}


func NewContext(input string) *Context {
	return &Context{
		input: input,
		memo:  make(map[memoKey]*memoEntry),
	}
}

func (ctx *Context) Match(node Node, pos int) []int {
	// fmt.Printf("MATCH %T %p @ %d\n", node, node, pos)

	// just in case
	if node == nil {
		panic("nil node in Context!")
	}

	if pos > ctx.CurrentPos {
		ctx.CurrentPos = pos
	}

	// 1. If already calculared - return cache
	key := memoKey{node: node, pos: pos}
	if entry, ok := ctx.memo[key]; ok {
		// if still counting -> left recurency
		// error tracker treats as a failure
		if entry.inProgress {
			fmt.Printf("LEFT RECURSION DETECTED: %T %p @ %d\n", node, node, pos)
			return nil
		}
		return entry.results
	}

	// 2. save calculations
	entry := &memoEntry{inProgress: true}
	ctx.memo[key] = entry

	// 3. calculate result (delegate to node)
	results := node.match(ctx, pos)

	// 4. save result
	entry.results = results
	entry.inProgress = false

	// 5. error tracking
	// we're looking for the farthest position reached
	// if pos < farthest, ignore
	// if pos > farthest, update farthest
	// if pos == farthest, merge expected tokens, to list them all
	if len(results) == 0 {
		if ctx.error == nil || ctx.CurrentPos > ctx.FarthestPos {
			ctx.FarthestPos = ctx.CurrentPos
            ctx.error = ctx.makeError(node)
		} else if ctx.CurrentPos == ctx.FarthestPos {
			// merge expected tokens
			ctx.error.Expected = mergeExpected(ctx.error.Expected, node.Expect())
		}
	}

	return results
}

func (ctx *Context) makeError(n Node) *ParseError {
	line, col := lineCol(ctx.input, ctx.FarthestPos)

	return &ParseError{
		Pos:       ctx.FarthestPos,
		Line:      line,
		Column:    col,
		RuleStack: ctx.stackSnapshot(),
		Expected:  n.Expect(),
		Found:     ctx.foundAt(ctx.FarthestPos),
		Width:     expectedWidth(n.Expect()),
	}
}

func mergeExpected(a, b []string) []string {
	seen := make(map[string]bool)
	var out []string

	for _, x := range a {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	for _, x := range b {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}


func lineCol(input string, pos int) (line, col int) {
	line = 1
	col = 1
	for i, r := range input {
		if i >= pos {
			break
		}
		if r == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return
}

func (ctx *Context) line(pos int) int {
	line, _ := lineCol(ctx.input, pos)
	return line
}

func (ctx *Context) col(pos int) int {
	_, col := lineCol(ctx.input, pos)
	return col
}

func (ctx *Context) foundAt(pos int) string {
	if pos >= len(ctx.input) {
		return "EOF"
	}
	r, _ := runeAt(ctx.input, pos)
	return fmt.Sprintf("%q", r)
}

func runeAt(s string, pos int) (rune, int) {
	for i, r := range s {
		if i == pos {
			return r, i
		}
	}
	return 0, 0
}

func expectedWidth(expected []string) int {
	max := 0
	for _, e := range expected {
		// interesują nas tylko string literals
		if len(e) >= 2 && e[0] == '"' && e[len(e)-1] == '"' {
			w := len([]rune(e[1 : len(e)-1]))
			if w > max {
				max = w
			}
		}
	}
	if max == 0 {
		return 1
	}
	return max
}



func (ctx *Context) push(name string) {
	ctx.stack = append(ctx.stack, name)
}

func (ctx *Context) pop() {
	if len(ctx.stack) == 0 {
		panic("pop on empty context stack")
	}
	ctx.stack = ctx.stack[:len(ctx.stack)-1]
}

func (ctx *Context) stackSnapshot() []string {
	if len(ctx.stack) == 0 {
		return nil
	}

	snap := make([]string, len(ctx.stack))
	copy(snap, ctx.stack)
	return snap
}
