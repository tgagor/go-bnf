### 1. Package Structure: Keep It Simple

**Recommendation:** For now, **stick with the single `bnf` package.** Do not split it into `lexer`, `parser`, and `ast`.

**Justification:**

*   **High Cohesion:** Your core types (`Node`, `Context`, `Grammar`, `Rule`, and all the expression nodes like `Sequence`, `Choice`, etc.) are highly interdependent. They all work together to form a single, coherent parsing engine. Splitting them would create a web of cross-package dependencies.
*   **No "Import Pollution":** In Go, a package is a unit of compilation. As long as your types are logically related, they belong together. You aren't polluting anything; you're creating a cohesive API.
*   **The Go Idiom:** It is very common in Go for libraries to be a single package, even a large one, if it exposes a single, primary domain. Think of libraries like `encoding/json` or `net/http`. The Go community generally prefers fewer, more cohesive packages over many small, chatty ones.
*   **The `ast` Exception:** If you were to build a rich, reusable Abstract Syntax Tree that could be used independently of the parser itself, a separate `ast` package might make sense. However, your current "AST" is the internal representation of the grammar rules, which is tightly coupled to the parser's logic.

Splitting the package now would lead to the exact problems you were worried about and make the code harder to navigate, not easier.

### 2. Public API Design: Hide the Internals

A great library exposes only what the user needs and hides the complex machinery.

**Recommendation:** Make most of your current types **internal (unexported)**.

Here is a proposed public vs. internal API:

**Public (Exported):**

*   `type Grammar`: The main entry point for a user.
*   `func LoadGrammar(r io.Reader) (*Grammar, error)`: A flexible way to load a grammar from any source (file, string, network). Your current `LoadGrammarFile` is too specific.
*   `func (g *Grammar) Parse(input string) (*ASTNode, error)`: The primary method. It should return a concrete result (an AST of the *parsed input*) and a detailed error. Right now, your `Match` method returns a `bool`, which is more of a "validate" function.
*   `type ParseError`: Your custom error type is great. Keep it and make it richer (see below).
*   `type ASTNode` (or similar): A new public type representing a node in the parse tree of the *input text*, not the grammar. This is what users typically want from a parser.

**Internal (Unexported):**

*   `type node`: The core interface `Node` should be unexported. Users shouldn't implement their own node types.
*   `type context`: The memoization context is an implementation detail and must be unexported.
*   `terminal`, `sequence`, `choice`, `repeat`, `nonTerminal`: All the concrete grammar node types should be unexported.
*   The `match` method on each node should be unexported.

**Example of an Improved API:**

```go
// in bnf.go

// Grammar is a compiled and ready-to-use BNF grammar.
type Grammar struct {
    // ... internal fields
}

// Parse parses the input string according to the grammar.
// It returns the root of the resulting parse tree or an error.
func (g *Grammar) Parse(input string) (*ASTNode, error) {
    // ... new logic using the internal context ...
}

// Validate checks if the input string fully matches the grammar.
func (g *Grammar) Validate(input string) (bool, error) {
    // This would be similar to your current Match() method
}

// LoadGrammar parses a grammar definition from a reader.
func LoadGrammar(r io.Reader) (*Grammar, error) {
    // ...
}

// ParseError provides detailed information about a parsing failure.
type ParseError struct {
    Line, Column int
    Message      string
    // ... more context
}

func (e *ParseError) Error() string { /* ... */ }

// ASTNode represents a node in the parsed output tree.
type ASTNode struct {
    Type     string
    Value    string
    Children []*ASTNode
}
```

This API is cleaner, hides the complexity, and provides more value to the user.

### 3. Solving Left Recursion: Grow the Seed

You are correct that your current memoization check (`if entry.inProgress { return nil }`) only prevents infinite loops; it doesn't correctly parse left-recursive rules. You need to implement the standard algorithm for handling this in packrat parsers.

**Recommendation:** Modify your `Context.Match` function to detect and handle left recursion by "growing the seed result."

Here is the pseudo-code for the required logic, which you can implement in your `Context.Match` function:

```
func (ctx *Context) Match(node Node, pos int) []int {
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

    // 1. Set up for potential left recursion.
    //    Assume failure (the "seed" result) and mark it as left-recursive.
    ctx.memo[key] = &memoEntry{isLeftRecursive: true, results: nil}

    // 2. Try to parse with the seed.
    var lastResults []int
    currentResults := node.match(ctx, pos)

    // 3. "Grow" the result. Keep parsing as long as the match gets longer.
    for !equal(currentResults, lastResults) {
        lastResults = currentResults
        // Important: Update the cache with the better result BEFORE re-evaluating.
        ctx.memo[key] = &memoEntry{isLeftRecursive: true, results: lastResults}
        currentResults = node.match(ctx, pos)
    }

    // 4. Finalize the result. Mark that we are done with this rule.
    ctx.memo[key] = &memoEntry{isLeftRecursive: false, results: currentResults}

    return currentResults
}
```
*Note: This is a simplified sketch. You'll need to adapt `memoEntry` and the logic carefully. The key is the loop that re-evaluates the rule, allowing the match to grow longer with each iteration until it stabilizes.*

### 4. Testing Strategy: Black-Box First

Your refactoring of the unit tests to use a `Context` was correct. However, you can improve the structure further.

**Recommendation:**

1.  **Create a `bnf_test` package:** For a library, your primary tests should be "black-box" tests. Create files like `grammar_test.go` inside your `bnf` directory, but with the package declaration `package bnf_test`.
2.  **Test the Public API:** In this `_test` package, you can only import `go-bnf/bnf` and use its *public* API, just like a real user would. This is the most important type of test. Your `TestRecursiveGrammar` is a perfect candidate for this package, as it calls `g.Match(...)`.
3.  **Keep `_test.go` files for white-box testing:** For testing tricky internal logic that the public API can't easily reach, it's fine to have your existing test files (e.g., `choice_test.go`) with `package bnf`. This is for "white-box" testing.

This hybrid approach gives you the best of both worlds: you ensure your public API works as advertised, and you can still unit-test complex internal machinery.

### 5. Documentation and Usability

**Recommendation:** Add GoDoc comments and usage examples.

*   **GoDoc:** Add comments to all your exported types and functions. Start each comment with the name of the thing it's describing.
    ```go
    // Grammar represents a compiled BNF grammar.
    type Grammar struct { /* ... */ }

    // LoadGrammar parses a grammar from an io.Reader.
    func LoadGrammar(r io.Reader) (*Grammar, error) { /* ... */ }
    ```
    Run `go doc ./bnf` or use a local `godoc` server to see how it will look.
*   **Examples:** Create an `examples_test.go` file (in `package bnf_test`). Use the special `Example` function format. These are compiled and run as tests, appear in the documentation, and are the best way to show how to use your library.

    ```go
    func ExampleGrammar_Validate() {
        grammar, err := bnf.LoadGrammar(strings.NewReader(`<number> ::= "0" | "1"`))
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(grammar.Validate("1"))
        fmt.Println(grammar.Validate("3"))
        // Output:
        // true
        // false
    }
    ```

By following these steps, you will significantly improve the structure, robustness, and user-friendliness of your library, turning your powerful parser engine into a proper, reusable Go module.
