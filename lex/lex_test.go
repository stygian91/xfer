package lex_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	l "github.com/stygian91/xfer/lex"
)

func TestSimples(t *testing.T) {
	it := l.StrIter2("()[]{}+-*/.,=<>!")
	lex := l.NewLexer(it)
	tokens, err := lex.Process()

	if err != nil {
		t.Errorf("Error while testing simples: %s", err)
		return
	}

	expected := []l.Token{
		{Kind: l.LPAREN, Literal: "(", Line: 1, Col: 1, Byte: 0},
		{Kind: l.RPAREN, Literal: ")", Line: 1, Col: 2, Byte: 1},
		{Kind: l.LSQUARE, Literal: "[", Line: 1, Col: 3, Byte: 2},
		{Kind: l.RSQUARE, Literal: "]", Line: 1, Col: 4, Byte: 3},
		{Kind: l.LCURLY, Literal: "{", Line: 1, Col: 5, Byte: 4},
		{Kind: l.RCURLY, Literal: "}", Line: 1, Col: 6, Byte: 5},
		{Kind: l.PLUS, Literal: "+", Line: 1, Col: 7, Byte: 6},
		{Kind: l.MINUS, Literal: "-", Line: 1, Col: 8, Byte: 7},
		{Kind: l.ASTERISK, Literal: "*", Line: 1, Col: 9, Byte: 8},
		{Kind: l.SLASH, Literal: "/", Line: 1, Col: 10, Byte: 9},
		{Kind: l.DOT, Literal: ".", Line: 1, Col: 11, Byte: 10},
		{Kind: l.COMMA, Literal: ",", Line: 1, Col: 12, Byte: 11},
		{Kind: l.EQUAL, Literal: "=", Line: 1, Col: 13, Byte: 12},
		{Kind: l.LT, Literal: "<", Line: 1, Col: 14, Byte: 13},
		{Kind: l.GT, Literal: ">", Line: 1, Col: 15, Byte: 14},
		{Kind: l.BANG, Literal: "!", Line: 1, Col: 16, Byte: 15},
	}

	if diff := cmp.Diff(expected, tokens); diff != "" {
		t.Errorf("TestSimples() mismatch (-want +got):\n%s", diff)
	}
}
