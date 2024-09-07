package parse

import (
	"fmt"

	"github.com/stygian91/xfer/lex"
)

type Parser struct {
	Tokens []lex.Token

	Idx int
}

func NewParser(tokens []lex.Token) Parser {
	return Parser{
		Tokens: tokens,
		Idx:    0,
	}
}

func (this *Parser) Optional(kind lex.TokenKind) (lex.Token, bool) {
	currToken := this.Tokens[this.Idx]

	if currToken.Kind != kind {
		return lex.Token{}, false
	}

	this.Idx += 1
	return currToken, true
}

func (this *Parser) Expect(kind lex.TokenKind) (lex.Token, error) {
	currToken := this.Tokens[this.Idx]

	if kind != currToken.Kind {
		return lex.Token{}, fmt.Errorf(
			"Unexpected token %s, expected %s at line %d, col %d",
			lex.KindString(currToken.Kind),
			lex.KindString(kind),
			currToken.Line,
			currToken.Col,
		)
	}

	this.Idx += 1
	return currToken, nil
}
