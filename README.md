# go-bnf

A simple, educational BNF (Backus-Naur Form) parser and matching engine written in Go.

This project implements a parser capable of loading BNF grammar definitions and verifying if an input text matches defined patterns. It supports modern features like Packrat parsing (memoization), left recursion handling, and AST (Abstract Syntax Tree) generation.

## Features

- **BNF Parsing**: generic parser for standard BNF syntax (`<rule> ::= expression`).
- **Pattern Matching**: Supports Sequences, Choices (`|`), Repetitions (`+`, `*`, `?`), grouping `()`, and literals.
- **Packrat Engine**: Efficient matching using memoization to ensure performance even with complex backtracking.
- **Left Recursion**: Direct support for left-recursive rules (e.g., `Expr ::= Expr "+" Term | Term`).
- **AST Generation**: Generate and visualize parse trees for successful matches.
- **Grammar Validation**: "Pre-flight" checks to ensure all referenced rules are defined.
- **CLI Utility**: Robust command-line tool with colored error reporting and visual pointers.

## Installation

```bash
make install
# or
go install
```

## Usage

### Validate a file
Checking if the whole content of `input.txt` matches `grammar.bnf`:
```bash
bnf -g grammar.bnf -i input.txt
```

### Parse and Show AST
Visualize the structure of the match:
```bash
bnf -g examples/numbers.bnf -i input.txt -p
```

### Validate Grammar Integrity
Only check if the grammar file itself is valid (no undefined rules):
```bash
bnf -g grammar.bnf -v
```

### Custom Start Rule
Override the entry point defined in the grammar (useful for testing sub-rules):
```bash
bnf -g examples/numbers.bnf -i input.txt -s "digit"
```

### Standard Input
Useful for piping content:
```bash
echo "123" | bnf -g examples/numbers.bnf
```

### Line-by-Line Mode
Check each line independently (useful for log files or lists):
```bash
bnf -g grammar.bnf -i input.txt -l
```

## Example Grammar

Standard `<rule> ::= ...` syntax. Literals can use `"` or `'`.

```bnf
<number>   ::= <digit>+
<digit>    ::= "0" | "1" | "2" | "3" | "4"
<expr>     ::= <expr> "+" <number> | <number> // Left recursion supported!
```

### Grammar of go-bnf itself

In a simplified form, the syntax supported by this tool is:

```bnf
<grammar>    ::= <rule>+
<rule>       ::= <identifier> "::=" <expression>
<expression> ::= <sequence> ( "|" <sequence> )*
<sequence>   ::= <factor>*
<factor>     ::= <atom> ( "*" | "+" | "?" )?
<atom>       ::= <identifier> | <string> | "(" <expression> ")"
```

## Development

Build and verify the project using the included Makefile:

- `make build`: Compile the binary.
- `make test`: Run unit tests and API tests.
- `make integration-test`: Run CLI tests against examples.
- `make lint`: Run code quality checks.
- `make docs`: Generate API documentation in the `docs/` folder.

## Architecture

The engine uses a modern Packrat parsing approach:

1.  **Grammar Parser**: Reads the `.bnf` file into a raw AST, then builds a linked execution graph.
2.  **Memoization**: Every (Node, Position) match result is cached to prevent exponential time complexity in ambiguous grammars.
3.  **Recursive Growth**: Handles left-recursive rules by iteratively "growing" the match from a seed failure until a fixed point is reached.
4.  **Node Interface**:
    ```go
    type node interface {
        match(ctx *context, pos int) ([]MatchResult, error)
        Expect() []string
    }
    ```

## License

[GNU GPL v3](./LICENSE)
