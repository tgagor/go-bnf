# go-bnf

A simple, educational BNF (Backus-Naur Form) parser and validator written in Go.

This project implements a parser engine capable of loading BNF grammar definitions and verifying if an input text matches the defined patterns. It was created as a learning exercise in compiler design, exploring recursive descent parsing, backtracking, and AST construction.

## Features

- **BNF Parsing**: Generic parser for standard BNF-like syntax.
- **Pattern Matching**: Supports Sequences, Choices (`|`), Repetitions (`+`, `*`), Optionals, and grouping.
- **CLI Utility**: Command-line tool to verify inputs against grammar files.
- **Backtracking support**: Handles ambiguous grammars via backtracking choices.
- **Line-by-line mode**: Validate file line by line or as a whole blob.

## Installation

```bash
make install
# or
go install
```

## Usage

### Validate a file

Checking if whole content of `input.txt` matches `grammar.bnf`:

```bash
bnf -g grammar.bnf -i input.txt
```

### Validate from Standard Input

Useful for piping content:

```bash
echo "123" | bnf -g examples/numbers.bnf
```

### Validate Line-by-Line

Check each line independently (useful for log files or lists):

```bash
bnf -g grammar.bnf -i input.txt -l
```

## Example Grammar

See `examples/` for more. Use standard `<rule> ::= ...` syntax.

```bnf
<number> ::= <digit>+
<digit>  ::= "0" | "1" | "2" | "3" | "4"
```

## Development

Build and verify the project using the included Makefile:

- `make build`: Compile the binary.
- `make test`: Run unit tests.
- `make integration-test`: Run CLI integration tests against examples.
- `make lint`: Run linters.

## Architecture

The core design separates **Grammar Parsing** from **Text Matching**:

1.  **Grammar Parser**: Reads the `.bnf` file and constructs an AST of the grammar rules. Recursive references are resolved in a second pass.
2.  **Matching Engine**: A recursive engine where every grammar component (`Sequence`, `Choice`, `Terminal`) implements a common `Node` interface:
    ```go
    type Node interface {
        Match(input string, pos int) []int
    }
    ```
    Returning a list of possible end positions allows for backtracking and handling ambiguity.
