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
		kind:    kind,
		subkind: NILKIND,
		literal: literal,
		byte:    uint(this.bytePos),
		line:    this.line,
		col:     this.col,
	})
}

func (this *Lexer) append(token Token) {
	this.tokens = append(this.tokens, token)
}

func (this *Lexer) finishIdent() {
	var kind TokenKind
	subkind, exists := reservedWords[this.currIdent]

	if exists {
		kind = KEYWORD
	} else {
		kind = IDENT
	}

	this.append(Token{
		kind:    kind,
		subkind: subkind,
		literal: this.currIdent,
		byte:    uint(this.identPos),
		line:    this.identLine,
		col:     this.identCol,
	})

	this.currIdent = ""
}

func (this *Lexer) finishNumber() {
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
		kind:    kind,
		subkind: NILKIND,
		literal: this.currNumber,
		byte:    uint(this.numberPos),
		line:    this.numberLine,
		col:     this.numberCol,
		value:   value,
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

		// TODO:
		slices.Contains(strEscapable, peekV)
	}
}
