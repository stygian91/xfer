package lex

import (
	"fmt"
	"strconv"
	"strings"
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

	var err error

	var line, col uint = 1, 0
	var char rune
	var valid bool
	var bytePos int

	var identLine, identCol uint
	var identPos int
	currIdent := ""

	var numberLine, numberCol uint
	var numberPos int
	currNumber := ""
	currNumberIsInt := true

	finishIdent := func() {
		subkind, exists := reservedWords[currIdent]
		if exists {
			tokens = append(tokens, Token{kind: KEYWORD, subkind: subkind, literal: currIdent, byte: uint(identPos), line: identLine, col: identCol})
		} else {
			tokens = append(tokens, Token{kind: IDENT, subkind: NILKIND, literal: currIdent, byte: uint(identPos), line: identLine, col: identCol})
		}
		currIdent = ""
	}

	finishNumber := func() {
		var kind TokenKind
		var value interface{}
		if currNumberIsInt {
			kind = INT
			value, err = strconv.ParseInt(currNumber, 10, 64)
		} else {
			kind = FLOAT
			value, err = strconv.ParseFloat(currNumber, 64)
		}

		if err != nil {
			return
		}

		tokens = append(tokens, Token{kind: kind, subkind: NILKIND, literal: currNumber, byte: uint(numberPos), line: numberLine, col: numberCol, value: value})

		currNumber = ""
		currNumberIsInt = true
	}

	newSimple := func(kind TokenKind, literal string) {
		if len(currIdent) > 0 {
			finishIdent()
		}

		if len(currNumber) > 0 {
			finishNumber()
		}

		tokens = append(tokens, Token{kind: kind, subkind: NILKIND, literal: literal, byte: uint(bytePos), line: line, col: col})
	}

	next, peek, stop := iter.Peek2(strIter(this.input))
	defer stop()

	for {
		if err != nil {
			return tokens, err
		}

		bytePos, char, valid = next()
		if !valid {
			break
		}

		col += 1

		if unicode.IsSpace(char) {
			if len(currIdent) > 0 {
				finishIdent()
			}

			if len(currNumber) > 0 {
				finishNumber()
			}
		}

		switch {
		case isValidFirstIdentRune(char) && len(currIdent) == 0:
			identPos = bytePos
			identLine = line
			identCol = col
			currIdent += string(char)
		case isValidIdentRune(char) && len(currIdent) > 0:
			currIdent += string(char)

		case unicode.IsDigit(char):
			if len(currNumber) == 0 {
				numberPos = bytePos
				numberLine = line
				numberCol = col
				currNumberIsInt = true
			}
			currNumber += string(char)

		case char == '.':
			if len(currNumber) > 0 && !strings.Contains(currNumber, ".") {
				currNumberIsInt = false
				currNumber += string(char)
			} else {
				newSimple(DOT, ".")
			}

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
