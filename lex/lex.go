package lex

import (
	"fmt"
	stditer "iter"

	"github.com/stygian91/iter-go"
)

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

func strIter(input string) stditer.Seq2[int, rune] {
	return func(yield func(int, rune) bool) {
		for i, char := range input {
			if !yield(i, char) {
				return
			}
		}
	}
}

func (this *Lexer) Process() ([]Token, error) {
	tokens := []Token{}
	var line, col uint = 1, 0

	next, peek, stop := iter.Peek2(strIter(this.input))
	defer stop()

	for {
		bytePos, char, valid := next()
		if !valid {
			break
		}

		newSimple := func(kind TokenKind, literal string) {
			tokens = append(tokens, SimpleToken{kind: kind, literal: literal, byte: uint(bytePos), line: line, col: col})
		}

		col += 1

		switch {
		case char == '(':
			newSimple(LPAREN, "(")
		case char == ')':
			newSimple(RPAREN, ")")
		case char == '[':
			newSimple(LSQUARE, "[")
		case char == ']':
			newSimple(RSQUARE, "]")
		case char == '\r':
			_, peekChar, peekValid := peek()
			if !peekValid {
				goto LOOPEND
			}

			if peekChar == '\n' {
				next()
			}

			line += 1
			col = 0
		case char == '\n':
			line += 1
			col = 0

		default:
			return tokens, fmt.Errorf("Unexpected character: '%c'", char)
		}
	}
LOOPEND:

	return tokens, nil
}
