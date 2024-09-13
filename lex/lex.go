package lex

import (
	"fmt"
	stditer "iter"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/stygian91/iter-go"
)

type Lexer struct {
	next func() (int, rune, bool)
	peek func() (int, rune, bool)
	stop func()

	tokens []Token
	err    error

	char    rune
	bytePos int
	valid   bool

	line uint
	col  uint

	currIdent string
	identLine uint
	identCol  uint
	identPos  int

	currNumber      string
	numberLine      uint
	numberCol       uint
	numberPos       int
	currNumberIsInt bool

	currStr    string
	strPos     int
	strLine    uint
	strCol     uint
	strEscNext bool
}

func NewLexer(seq stditer.Seq2[int, rune]) Lexer {
	next, peek, stop := iter.Peek2(seq)
	tokens := []Token{}

	return Lexer{
		next:            next,
		peek:            peek,
		stop:            stop,
		tokens:          tokens,
		line:            1,
		currNumberIsInt: true,
	}
}

func (this *Lexer) Process() ([]Token, error) {
	defer this.stop()

	for {
		if this.err != nil {
			return this.tokens, this.err
		}

		this.bytePos, this.char, this.valid = this.next()
		if !this.valid {
			endErr := this.ending()
			if endErr != nil {
				return this.tokens, endErr
			}

			break
		}

		this.col += 1

		if unicode.IsSpace(this.char) {
			if len(this.currIdent) > 0 {
				this.finishIdent()
			}

			if len(this.currNumber) > 0 {
				this.finishNumber()
			}
		}

		switch {
		case len(this.currStr) == 0 && this.char == '"':
			this.currStr = "\""
			this.strLine = this.line
			this.strCol = this.col
			this.strPos = this.bytePos
		case len(this.currStr) > 0:
			this.handleString()

		case isValidFirstIdentRune(this.char) && len(this.currIdent) == 0:
			this.identPos = this.bytePos
			this.identLine = this.line
			this.identCol = this.col
			this.currIdent += string(this.char)
		case isValidIdentRune(this.char) && len(this.currIdent) > 0:
			this.currIdent += string(this.char)

		case unicode.IsDigit(this.char):
			if len(this.currNumber) == 0 {
				this.numberPos = this.bytePos
				this.numberLine = this.line
				this.numberCol = this.col
				this.currNumberIsInt = true
			}
			this.currNumber += string(this.char)

		case this.char == '.':
			if len(this.currNumber) > 0 && !strings.Contains(this.currNumber, ".") {
				this.currNumberIsInt = false
				this.currNumber += string(this.char)
			} else {
				this.addSimple(DOT, ".")
			}

		case this.char == ',':
			this.addSimple(COMMA, ",")
		case this.char == ';':
			this.addSimple(SEMICOLON, ";")
		case this.char == '=':
			this.addSimple(EQUAL, "=")
		case this.char == '<':
			this.addSimple(LT, "<")
		case this.char == '>':
			this.addSimple(GT, ">")
		case this.char == '!':
			this.addSimple(BANG, "!")
		case this.char == '(':
			this.addSimple(LPAREN, "(")
		case this.char == ')':
			this.addSimple(RPAREN, ")")
		case this.char == '[':
			this.addSimple(LSQUARE, "[")
		case this.char == ']':
			this.addSimple(RSQUARE, "]")
		case this.char == '{':
			this.addSimple(LCURLY, "{")
		case this.char == '}':
			this.addSimple(RCURLY, "}")
		case this.char == '+':
			this.addSimple(PLUS, "+")
		case this.char == '-':
			this.addSimple(MINUS, "-")
		case this.char == '*':
			this.addSimple(ASTERISK, "*")
		case this.char == '/':
			this.addSimple(SLASH, "/")
		case this.char == '\r':
			_, peekChar, peekValid := this.peek()
			if !peekValid {
				goto LOOPEND
			}

			if peekChar == '\n' {
				this.next()
			}

			this.line += 1
			this.col = 0
		case this.char == '\n':
			this.line += 1
			this.col = 0

		case unicode.IsSpace(this.char):

		default:
			return this.tokens, fmt.Errorf("Unexpected character: '%c'", this.char)
		}
	}
LOOPEND:

	return this.tokens, nil
}

func (this *Lexer) addSimple(kind TokenKind, literal string) {
	if len(this.currIdent) > 0 {
		this.finishIdent()
	}

	if len(this.currNumber) > 0 {
		this.finishNumber()
	}

	this.append(Token{
		Kind:    kind,
		Literal: literal,
		Byte:    uint(this.bytePos),
		Line:    this.line,
		Col:     this.col,
	})
}

func (this *Lexer) append(token Token) {
	this.tokens = append(this.tokens, token)
}

func (this *Lexer) finishIdent() {
	if len(this.currIdent) == 0 {
		return
	}

	var kind TokenKind
	reserved, exists := reservedWords[this.currIdent]

	if exists {
		kind = reserved
	} else {
		kind = IDENT
	}

	this.append(Token{
		Kind:    kind,
		Literal: this.currIdent,
		Byte:    uint(this.identPos),
		Line:    this.identLine,
		Col:     this.identCol,
	})

	this.currIdent = ""
}

func (this *Lexer) finishNumber() {
	if len(this.currNumber) == 0 {
		return
	}

	var kind TokenKind
	var value interface{}
	if this.currNumberIsInt {
		kind = INT
		value, this.err = strconv.ParseInt(this.currNumber, 10, 64)
	} else {
		kind = FLOAT
		value, this.err = strconv.ParseFloat(this.currNumber, 64)
	}

	if this.err != nil {
		return
	}

	this.append(Token{
		Kind:    kind,
		Literal: this.currNumber,
		Byte:    uint(this.numberPos),
		Line:    this.numberLine,
		Col:     this.numberCol,
		Value:   value,
	})

	this.currNumber = ""
	this.currNumberIsInt = true
}

func (this *Lexer) handleString() {
	switch {
	case this.strEscNext:
		this.currStr += string(this.char)
		this.strEscNext = false
	case this.char == '\\':
		_, peekV, peekValid := this.peek()
		if !peekValid {
			this.err = fmt.Errorf("Error while parsing string, unexpected EOF")
			return
		}

		if !slices.Contains(strEscapable, peekV) {
			this.err = fmt.Errorf("Unknown escape sequence")
			return
		}

		this.strEscNext = true
		this.currStr += string(this.char)
	case this.char == '"':
		value := unescape(this.currStr[1:])
		this.currStr += string(this.char)
		this.append(Token{
			Byte:    uint(this.strPos),
			Line:    this.strLine,
			Col:     this.strCol,
			Kind:    STRING,
			Literal: this.currStr,
			Value:   value,
		})
		this.currStr = ""

	case this.char == '\r' || this.char == '\n':
		this.err = fmt.Errorf("\\r and \\n are not allowed in strings, consider escaping them")
		return

	default:
		this.currStr += string(this.char)
	}
}

func (this *Lexer) ending() error {
	this.finishIdent()
	this.finishNumber()
	if len(this.currStr) > 0 {
		return fmt.Errorf("Unexpected EOF")
	}

	return nil
}
