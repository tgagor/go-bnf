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
}

func NewContext(input string) *Context {
	return &Context{
		input: input,
		memo:  make(map[memoKey]*memoEntry),
	}
}

func (ctx *Context) Match(node Node, pos int) []int {
	// fmt.Printf("MATCH %T %p @ %d\n", node, node, pos)
	key := memoKey{node: node, pos: pos}

	// just in case
	if node == nil {
		panic("nil node in Context!")
	}

	// 1. If already calculared - return cache
	if entry, ok := ctx.memo[key]; ok {
		// if still counting -> left recurency
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

	return results
}
