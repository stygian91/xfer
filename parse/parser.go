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

type ExpectKind int

const (
	EXPECSEQNILKIND = iota
	OPTIONAL
	EXPECT
	EXPECTANY
)

type ExpectSeqHandler interface {
	ExpectKind() ExpectKind
	TokenKinds() []lex.TokenKind
	Handle(*Parser, lex.Token) (Node, error)
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

func (this *Parser) ParseSeq(handlers []ExpectSeqHandler) ([]Node, error) {
	nodes := []Node{}

	for _, handler := range handlers {
		tokenKinds := handler.TokenKinds()
		if len(tokenKinds) == 0 {
			return []Node{}, fmt.Errorf("Parser.ExpectSeq(): no token kinds provided")
		}

		switch handler.ExpectKind() {
		case OPTIONAL:
			token, exists := this.Optional(tokenKinds[0])
			if !exists {
				break
			}
			node, err := handler.Handle(this, token)
			if err != nil {
				return []Node{}, fmt.Errorf("Parser.ExpectSeq(): error in Optional handler: %w", err)
			}

			nodes = append(nodes, node)
		case EXPECT:
			token, err := this.Expect(tokenKinds[0])
			if err != nil {
				return []Node{}, fmt.Errorf("Parser.ExpectSeq(): error in Expect: %w", err)
			}

			node, err := handler.Handle(this, token)
			if err != nil {
				return []Node{}, fmt.Errorf("Parser.ExpectSeq(): error in Expect handler: %w", err)
			}

			nodes = append(nodes, node)
		case EXPECTANY:
			token, err := this.ExpectAny(tokenKinds)
			if err != nil {
				return []Node{}, fmt.Errorf("Parser.ExpectSeq(): error in ExpectAny: %w", err)
			}

			node, err := handler.Handle(this, token)
			if err != nil {
				return []Node{}, fmt.Errorf("Parser.ExpectSeq(): error in ExpectAny handler: %w", err)
			}

			nodes = append(nodes, node)

		default:
			return []Node{}, fmt.Errorf("Invalid expect kind %+v", handler.ExpectKind())
		}
	}

	return nodes, nil
}
