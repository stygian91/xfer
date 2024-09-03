package lex

import (
	"fmt"
	"unicode"

	"github.com/stygian91/iter-go"
)

type Lexer struct {
	input string
}

func NewLexer(input string) Lexer {
	return Lexer{
		input: input,
	}
}

func (this *Lexer) Process() ([]Token, error) {
	tokens := []Token{}
	var line, col, identLine, identCol uint = 1, 0, 0, 0
	var bytePos, identPos int
	var char rune
	var valid bool
	currIdent := ""

	finishIdent := func() {
		subkind, exists := allKeywords[currIdent]
		if exists {
			tokens = append(tokens, Token{kind: KEYWORD, subkind: subkind, literal: currIdent, byte: uint(identPos), line: identLine, col: identCol})
		} else {
			tokens = append(tokens, Token{kind: IDENT, subkind: NILKIND, literal: currIdent, byte: uint(identPos), line: identLine, col: identCol})
		}
		currIdent = ""
	}

	newSimple := func(kind TokenKind, literal string) {
		if len(currIdent) > 0 {
			finishIdent()
		}

		tokens = append(tokens, Token{kind: kind, subkind: NILKIND, literal: literal, byte: uint(bytePos), line: line, col: col})
	}

	next, peek, stop := iter.Peek2(strIter(this.input))
	defer stop()

	for {
		bytePos, char, valid = next()
		if !valid {
			break
		}

		col += 1

		if unicode.IsSpace(char) && len(currIdent) > 0 {
			finishIdent()
		}

		switch {
		case isValidIdentRune(char):
			if len(currIdent) == 0 {
				identPos = bytePos
				identLine = line
				identCol = col
			}
			currIdent += string(char)

		case char == '(':
			newSimple(LPAREN, "(")
		case char == ')':
			newSimple(RPAREN, ")")
		case char == '[':
			newSimple(LSQUARE, "[")
		case char == ']':
			newSimple(RSQUARE, "]")
		case char == '{':
			newSimple(LCURLY, "{")
		case char == '}':
			newSimple(RCURLY, "}")
		case char == '+':
			newSimple(PLUS, "+")
		case char == '-':
			newSimple(MINUS, "-")
		case char == '*':
			newSimple(ASTERISK, "*")
		case char == '/':
			newSimple(SLASH, "/")
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

		case unicode.IsSpace(char):

		default:
			return tokens, fmt.Errorf("Unexpected character: '%c'", char)
		}
	}
LOOPEND:

	return tokens, nil
}
