package lex

import "fmt"

type TokenKind int
type TokenSubkind int

type Lexer struct {
	input  string
	cursor uint
}

type Token interface {
	Byte() uint
	Line() uint
	Col() uint
	Kind() TokenKind
	Subkind() TokenSubkind
	Literal() string
	Value() interface{}
}

type SimpleToken struct {
	byte    uint
	line    uint
	col     uint
	kind    TokenKind
	literal string
}

func (this SimpleToken) Byte() uint            { return this.byte }
func (this SimpleToken) Line() uint            { return this.line }
func (this SimpleToken) Col() uint             { return this.col }
func (this SimpleToken) Kind() TokenKind       { return this.kind }
func (this SimpleToken) Subkind() TokenSubkind { return NILKIND }
func (this SimpleToken) Literal() string       { return this.literal }
func (this SimpleToken) Value() interface{}    { return nil }

const NILKIND = 0

// tokens
const (
	LPAREN = iota + 1
	RPAREN
	LSQUARE
	RSQUARE
	LCURLY
	RCURLY
	IDENT
	KEYWORD
	INT
	FLOAT
	STRING
	BOOL
	PLUS
	MINUS
	ASTERISK
	SLASH
)

// keywords
const (
	STRUCT = iota + 1
	IF
	ELSE
	EXPORT
)

func NewLexer(input string) Lexer {
	return Lexer{
		input:  input,
		cursor: 0,
	}
}

func (this *Lexer) Process() ([]Token, error) {
	tokens := []Token{}

	// TODO: track current line & col
	for bytePos, char := range this.input {
		switch {
		case char == '(':
			tokens = append(tokens, SimpleToken{
				kind:    LPAREN,
				literal: "(",
				byte:    uint(bytePos),
			})

		case char == ')':
			tokens = append(tokens, SimpleToken{
				kind:    LPAREN,
				literal: ")",
				byte:    uint(bytePos),
			})

		case char == '\n':
			// TODO:

		default:
			return tokens, fmt.Errorf("Unexpected character: '%c'", char)
		}
	}

	return tokens, nil
}
