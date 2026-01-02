package bnf

import (
	"fmt"
	"slices"
)

type memoKey struct {
	node node
	pos  int
}

type memoEntry struct {
	results         []int // remember results`
	isLeftRecursive bool  // detect left sided recurency
}

type context struct {
	input string
	memo  map[memoKey]*memoEntry // cache (node, pos)

	// debug
	FarthestPos int // farthest position reached during parsing
	CurrentPos  int // current position during parsing (deepest, even for Sequence, Choice, Repeat)
	Expected    []string
	error       *ParseError

	// call stack
	stack []string
}

func NewContext(input string) *context {
	return &context{
		input: input,
		memo:  make(map[memoKey]*memoEntry),
	}
}

func (ctx *context) Match(node node, pos int) []int {
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
		if entry.isLeftRecursive {
			// This is the crucial part for handling left recursion.
			// We've been here before in a recursive call.
			// We return the "seed" result but mark that the rule needs re-evaluation.
			return entry.results
		}
		return entry.results
	}


    // 2. Set up for potential left recursion.
    //    Assume failure (the "seed" result) and mark it as left-recursive.
	entry := &memoEntry{isLeftRecursive: true, results: nil}
	ctx.memo[key] = entry

    // 2. Try to parse with the seed.
	//	  This may recursively call back into this same function.
    // 2. Try to parse with the seed.
    var lastResults []int
    currentResults := node.match(ctx, pos)

    // 3. "Grow" the result. Keep parsing as long as the match gets longer.
    for !slices.Equal(currentResults, lastResults) {
        lastResults = currentResults
        // Important: Update the cache with the better result BEFORE re-evaluating.
        ctx.memo[key] = &memoEntry{isLeftRecursive: true, results: lastResults}
        currentResults = node.match(ctx, pos)
    }

    // 4. Finalize the result. Mark that we are done with this rule.
    ctx.memo[key] = &memoEntry{isLeftRecursive: false, results: currentResults}

	// 5. error tracking
	// we're looking for the farthest position reached
	// if pos < farthest, ignore
	// if pos > farthest, update farthest
	// if pos == farthest, merge expected tokens, to list them all
	if len(currentResults) == 0 {
		if ctx.error == nil || ctx.CurrentPos > ctx.FarthestPos {
			ctx.FarthestPos = ctx.CurrentPos
			ctx.error = ctx.makeError(node)
		} else if ctx.CurrentPos == ctx.FarthestPos {
			// merge expected tokens
			ctx.error.Expected = mergeExpected(ctx.error.Expected, node.Expect())
		}
	}

	return currentResults
}

func (ctx *context) makeError(n node) *ParseError {
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

func (ctx *context) line(pos int) int {
	line, _ := lineCol(ctx.input, pos)
	return line
}

func (ctx *context) col(pos int) int {
	_, col := lineCol(ctx.input, pos)
	return col
}

func (ctx *context) foundAt(pos int) string {
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

func (ctx *context) push(name string) {
	ctx.stack = append(ctx.stack, name)
}

func (ctx *context) pop() {
	if len(ctx.stack) == 0 {
		panic("pop on empty context stack")
	}
	ctx.stack = ctx.stack[:len(ctx.stack)-1]
}

func (ctx *context) stackSnapshot() []string {
	if len(ctx.stack) == 0 {
		return nil
	}

	snap := make([]string, len(ctx.stack))
	copy(snap, ctx.stack)
	return snap
}
