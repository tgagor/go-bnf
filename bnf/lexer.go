package bnf

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// TokenType represents the category of a lexed token.
type TokenType int

const (
	EOF      TokenType = iota
	IDENT              // generic identifier or BNF rule name
	NT_IDENT           // BNF rule name explicitly in angle brackets (<rule>)
	STRING             // quoted string literal
	REGEX              // regex pattern enclosed in /.../
	ASSIGN             // the ::= operator
	PIPE               // the | operator
	STAR               // the * operator
	PLUS               // the + operator
	QMARK              // the ? operator
	LPAREN             // the ( operator
	RPAREN             // the ) operator
)

// Token represents a single atom (lexeme) in the input BNF grammar.
type Token struct {
	Type TokenType
	Text string
}

// Lexer breaks the BNF grammar input into a stream of tokens.
type Lexer struct {
	r *bufio.Reader
}

// NewLexer creates a new Lexer for the given reader.
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(r)}
}

// Next returns the next token from the input stream.
func (l *Lexer) Next() (Token, error) {
	// 1. skip whitespace and comments
	var ch rune
	var err error
	for {
		ch, _, err = l.r.ReadRune()
		if err == io.EOF {
			return Token{Type: EOF}, nil
		}
		if err != nil {
			return Token{}, err
		}

		if isWhitespace(ch) {
			continue
		}

		// line comment
		if ch == ';' || ch == '#' {
			if err := l.skipUntilEOL(); err != nil {
				if err == io.EOF {
					return Token{Type: EOF}, nil
				}
				return Token{}, err
			}
			continue
		}

		if ch == '/' {
			// Peek at the next byte to see if it's another '/'
			peek, err := l.r.Peek(1)
			if err == nil && len(peek) > 0 && peek[0] == '/' {
				l.r.ReadRune() // consume the second '/'
				if err := l.skipUntilEOL(); err != nil {
					if err == io.EOF {
						return Token{Type: EOF}, nil
					}
					return Token{}, err
				}
				continue
			}
			// Not a comment, so the / is the start of a token (likely a regex)
			// We have '/' in 'ch', and we break the loop to process it.
			break
		}

		// ch is meaningful, so we break the loop to process it in the main switch
		break
	}

	// 2. identificator
	if isIdentStart(ch) {
		var sb strings.Builder
		sb.WriteRune(ch)

		for {
			ch, _, err := l.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					break // legitimate end of ident
				}
				return Token{}, err
			}
			if !isIdentPart(ch) {
				l.r.UnreadRune()
				break
			}
			sb.WriteRune(ch)
		}

		return Token{
			Type: IDENT,
			Text: sb.String(),
		}, nil
	}

	// 2.5. <identificator>
	if ch == '<' {
		var sb strings.Builder

		for {
			ch, _, err := l.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return Token{}, fmt.Errorf("unterminated <identifier>")
				}
				return Token{}, err
			}
			if ch == '>' {
				break
			}
			if !isIdentPart(ch) {
				return Token{}, fmt.Errorf("invalid character in <identifier>: %q", ch)
			}
			sb.WriteRune(ch)
		}

		if sb.Len() == 0 {
			return Token{}, fmt.Errorf("empty <identifier>")
		}

		return Token{
			Type: IDENT, // normalize NT_IDENT as <> are just a sugar here
			Text: sb.String(),
		}, nil
	}

	// 3. regex literal /pattern/
	if ch == '/' {
		var sb strings.Builder

		for {
			ch, _, err := l.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return Token{}, fmt.Errorf("unterminated regex literal")
				}
				return Token{}, err
			}

			if ch == '/' {
				break
			}

			// Handle escape sequences in regex
			if ch == '\\' {
				esc, _, err := l.r.ReadRune()
				if err != nil {
					if err == io.EOF {
						return Token{}, fmt.Errorf("unterminated escape sequence in regex")
					}
					return Token{}, err
				}
				// Preserve the backslash for regex engine
				sb.WriteRune('\\')
				sb.WriteRune(esc)
				continue
			}

			sb.WriteRune(ch)
		}

		return Token{
			Type: REGEX,
			Text: sb.String(),
		}, nil
	}

	// 4. string literal
	if ch == '"' || ch == '\'' {
		quote := ch
		var sb strings.Builder

		for {
			ch, _, err := l.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return Token{}, fmt.Errorf("unterminated string literal")
				}
				return Token{}, err
			}

			if ch == quote {
				break
			}

			if ch == '\\' {
				esc, _, err := l.r.ReadRune()
				if err != nil {
					if err == io.EOF {
						return Token{}, fmt.Errorf("unterminated escape sequence")
					}
					return Token{}, err
				}
				switch esc {
				case '"':
					sb.WriteRune('"')
				case '\'':
					sb.WriteRune('\'')
				case '\\':
					sb.WriteRune('\\')
				case 'n':
					sb.WriteRune('\n')
				case 't':
					sb.WriteRune('\t')
				default:
					return Token{}, fmt.Errorf("unknown escape sequence: \\%c", esc)
				}
				continue
			}

			sb.WriteRune(ch)
		}

		return Token{
			Type: STRING,
			Text: sb.String(),
		}, nil
	}

	// 5. ASSIGN ::=
	if ch == ':' {
		ch2, _, err := l.r.ReadRune()
		if err == nil && ch2 == ':' {
			ch3, _, err := l.r.ReadRune()
			if err == nil && ch3 == '=' {
				return Token{Type: ASSIGN, Text: "::="}, nil
			}
		}
		return Token{}, fmt.Errorf("expected ::=")
	}

	// 6. Single symbols
	switch ch {
	case '|':
		return Token{Type: PIPE, Text: "|"}, nil
	case '*':
		return Token{Type: STAR, Text: "*"}, nil
	case '+':
		return Token{Type: PLUS, Text: "+"}, nil
	case '?':
		return Token{Type: QMARK, Text: "?"}, nil
	case '(':
		return Token{Type: LPAREN, Text: "("}, nil
	case ')':
		return Token{Type: RPAREN, Text: ")"}, nil
	}

	return Token{}, fmt.Errorf("unexpected character: %q", ch)
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}

func isIdentStart(ch rune) bool {
	return ch == '_' ||
		(ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z')
}

func isIdentPart(ch rune) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9') || ch == '-'
}

func (l *Lexer) skipUntilEOL() error {
	for {
		ch, _, err := l.r.ReadRune()
		if err != nil {
			return err
		}

		// Linux new line
		if ch == '\n' {
			return nil
		}
		// Windows, or old Mac style new line
		if ch == '\r' {
			// check if it's not \r\n
			next, _, err := l.r.ReadRune()
			if err == nil && next != '\n' {
				l.r.UnreadRune()
			}
			return nil
		}
	}
}
