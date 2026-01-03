package bnf

import (
	"fmt"
)

type memoKey struct {
	node node
	pos  int
}

type memoEntry struct {
	results         []MatchResult // remember results`
	isLeftRecursive bool          // detect left sided recurency
}

type context struct {
	input        string
	memo         map[memoKey]*memoEntry // cache (node, pos)
	activeCounts map[int]int            // recursive rules count per position

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
		input:        input,
		memo:         make(map[memoKey]*memoEntry),
		activeCounts: make(map[int]int),
	}
}

func (ctx *context) Match(node node, pos int) ([]MatchResult, error) {
	// Simple safety check.
	if node == nil {
		return nil, fmt.Errorf("nil node in Context")
	}

	// Track the deepest position reached for error reporting.
	if pos > ctx.CurrentPos {
		ctx.CurrentPos = pos
	}

	// 1. Check Memoization Cache
	// If we have a result for this (node, pos) pair, return it immediately.
	key := memoKey{node: node, pos: pos}
	if entry, ok := ctx.memo[key]; ok {
		if entry.isLeftRecursive {
			// CRITICAL: Left Recursion Handling.
			// If we hit a cache entry marked 'isLeftRecursive', it means we are inside
			// a recursive call for this same rule at this same position.
			// We return the current "seed" result (which starts as failure/nil)
			// to allow the parser to make progress and potentially find a base case.
			return entry.results, nil
		}
		return entry.results, nil
	}

	// 2. Initialize for Left Recursion
	// We haven't visited this node at this pos yet. To handle potential direct or indirect
	// left recursion, we insert a "seed" entry into the memo table. This seed assumes
	// failure (nil results) initially. If we recurse back here, step 1 will return this nil.
	entry := &memoEntry{isLeftRecursive: true, results: nil}
	ctx.memo[key] = entry
	ctx.activeCounts[pos]++

	// 3. First Parse Attempt
	// Compute the result using the current seed.
	var lastResults []MatchResult
	currentResults, err := node.match(ctx, pos)
	if err != nil {
		ctx.activeCounts[pos]--
		delete(ctx.memo, key)
		return nil, err
	}

	// 4. Grow the Seed (Iterative Step)
	// If this rule is left-recursive, we might have computed a result based on the nil seed.
	// Now we update the cache with this new result and try to parse AGAIN.
	// This allows the recursion to "grow" (unroll) one more level.
	// We repeat this until the results stop changing (stabilize).
	// We compare 'End' positions to detect stability.
	for !resultsEndEqual(currentResults, lastResults) {
		lastResults = currentResults

		// Update cache with the latest better result before re-evaluating.
		ctx.memo[key] = &memoEntry{isLeftRecursive: true, results: lastResults}

		currentResults, err = node.match(ctx, pos)
		if err != nil {
			ctx.activeCounts[pos]--
			delete(ctx.memo, key)
			return nil, err
		}
	}

	// 5. Cleanup and Finalize
	// We are done with this node at this position.
	ctx.activeCounts[pos]--
	if ctx.activeCounts[pos] > 0 {
		// If other rules are still active at this position, our result might depend on
		// their temporary seeds. We cannot permanently memoize this result yet because
		// it might change when those parent rules grow.
		delete(ctx.memo, key)
	} else {
		// No active recursion stack at this position, so this result is final and safe to cache.
		ctx.memo[key] = &memoEntry{isLeftRecursive: false, results: currentResults}
	}

	// 6. Error Tracking
	// If we failed to match anything, we record error information.
	// We track the "farthest failure" (deepest position) to provide helpful error messages.
	// If we reached the same farthest position again, we merge the expected tokens.
	if len(currentResults) == 0 {
		if ctx.error == nil || ctx.CurrentPos > ctx.FarthestPos {
			ctx.FarthestPos = ctx.CurrentPos
			ctx.error = ctx.makeError(node)
		} else if ctx.CurrentPos == ctx.FarthestPos {
			ctx.error.Expected = mergeExpected(ctx.error.Expected, node.Expect())
		}
	}

	return currentResults, nil
}

func resultsEndEqual(a, b []MatchResult) bool {
	if len(a) != len(b) {
		return false
	}
	// Sort? Order matters? Standard logic usually assumes ordered or stable order.
	// Our match implementations (sequence, choice) are deterministic in order.
	for i := range a {
		if a[i].End != b[i].End {
			return false
		}
	}
	return true
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

func (ctx *context) pop() error {
	if len(ctx.stack) == 0 {
		return fmt.Errorf("pop on empty context stack")
	}
	ctx.stack = ctx.stack[:len(ctx.stack)-1]
	return nil
}

func (ctx *context) stackSnapshot() []string {
	if len(ctx.stack) == 0 {
		return nil
	}

	snap := make([]string, len(ctx.stack))
	copy(snap, ctx.stack)
	return snap
}
