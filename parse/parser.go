package parse

import (
	"fmt"
	"slices"

	"github.com/stygian91/xfer/lex"
)

type Parser struct {
	tokens []lex.Token

	idx int
}

func NewParser(tokens []lex.Token) Parser {
	return Parser{
		tokens: tokens,
		idx:    0,
	}
}

func (this Parser) HasMore() bool {
	return this.idx < len(this.tokens)
}

func (this Parser) CurrentToken() (lex.Token, bool) {
	if !this.HasMore() {
		return lex.Token{}, false
	}

	return this.tokens[this.idx], true
}

func (this *Parser) Optional(kind lex.TokenKind) (lex.Token, bool) {
	currToken, exists := this.CurrentToken()
	if !exists {
		return lex.Token{}, false
	}

	if currToken.Kind != kind {
		return lex.Token{}, false
	}

	this.idx += 1
	return currToken, true
}

func (this *Parser) Expect(kind lex.TokenKind) (lex.Token, error) {
	currToken, exists := this.CurrentToken()

	if !exists {
		return lex.Token{}, fmt.Errorf("Unexpected end of token stream, expected %s", lex.KindString(kind))
	}

	if kind != currToken.Kind {
		return lex.Token{}, fmt.Errorf(
			"Unexpected token %s, expected %s at line %d, col %d",
			lex.KindString(currToken.Kind),
			lex.KindString(kind),
			currToken.Line,
			currToken.Col,
		)
	}

	this.idx += 1
	return currToken, nil
}

func (this *Parser) ExpectAny(kinds []lex.TokenKind) (lex.Token, error) {
	currToken, exists := this.CurrentToken()

	if !exists {
		return lex.Token{}, fmt.Errorf("Unexpected end of token stream, expected one of %s", lex.KindsString(kinds))
	}

	if !slices.Contains(kinds, currToken.Kind) {
		return lex.Token{}, fmt.Errorf(
			"Unexpected token %s, expected one of %s at line %d, col %d",
			lex.KindString(currToken.Kind),
			lex.KindsString(kinds),
			currToken.Line,
			currToken.Col,
		)
	}

	this.idx += 1
	return currToken, nil
}

func (this *Parser) ExpectSeq(kinds []lex.TokenKind) ([]lex.Token, error) {
	var err error
	tokens := []lex.Token{}

	startIdx := this.idx
	defer (func() {
		if err != nil {
			this.idx = startIdx
		}
	})()

	for _, kind := range kinds {
		currToken, exists := this.CurrentToken()

		if !exists {
			err = fmt.Errorf("Unexpected end of token stream, expected %s", lex.KindString(kind))
			return tokens, err
		}

		if currToken.Kind != kind {
			err = fmt.Errorf(
				"Unexpected token %s, expected %s at line %d, col %d",
				lex.KindString(currToken.Kind),
				lex.KindString(kind),
				currToken.Line,
				currToken.Col,
			)
			return tokens, err
		}

		tokens = append(tokens, currToken)
		this.idx += 1
	}

	return tokens, err
}
