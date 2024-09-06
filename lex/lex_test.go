package lex_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	l "github.com/stygian91/xfer/lex"
)

func checkLexResults(t *testing.T, input string, expected []l.Token) {
	it := l.StrIter2(input)
	lex := l.NewLexer(it)
	tokens, err := lex.Process()

	if err != nil {
		t.Errorf("Error while running lex.Process(): %s", err)
		return
	}

	if diff := cmp.Diff(expected, tokens); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestSimples(t *testing.T) {
	checkLexResults(t, "()[]{}+-*/.,=<>!", []l.Token{
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
	})
}

func TestNumbers(t *testing.T) {
	checkLexResults(t, "123 42.69", []l.Token{
		{Kind: l.INT, Literal: "123", Value: int64(123), Line: 1, Col: 1, Byte: 0},
		{Kind: l.FLOAT, Literal: "42.69", Value: 42.69, Line: 1, Col: 5, Byte: 4},
	})
}

func TestIdents(t *testing.T) {
	checkLexResults(t, "asd if else export true false struct enum", []l.Token{
		{Kind: l.IDENT, Literal: "asd", Line: 1, Col: 1, Byte: 0},
		{Kind: l.KEYWORD, Subkind: l.IF, Literal: "if", Line: 1, Col: 5, Byte: 4},
		{Kind: l.KEYWORD, Subkind: l.ELSE, Literal: "else", Line: 1, Col: 8, Byte: 7},
		{Kind: l.KEYWORD, Subkind: l.EXPORT, Literal: "export", Line: 1, Col: 13, Byte: 12},
		{Kind: l.KEYWORD, Subkind: l.TRUE, Literal: "true", Line: 1, Col: 20, Byte: 19},
		{Kind: l.KEYWORD, Subkind: l.FALSE, Literal: "false", Line: 1, Col: 25, Byte: 24},
		{Kind: l.KEYWORD, Subkind: l.STRUCT, Literal: "struct", Line: 1, Col: 31, Byte: 30},
		{Kind: l.KEYWORD, Subkind: l.ENUM, Literal: "enum", Line: 1, Col: 38, Byte: 37},
	})
}

func TestString(t *testing.T) {
	checkLexResults(t, "\"asd\"", []l.Token{{
		Kind:l.STRING, Literal: "\"asd\"", Value: "asd", Line: 1, Col: 1, Byte: 0,
	}})

	checkLexResults(t, "\"asd\\r\"", []l.Token{{
		Kind:l.STRING, Literal: "\"asd\\r\"", Value: "asd\r", Line: 1, Col: 1, Byte: 0,
	}})
	
	checkLexResults(t, "\"asd\\\\r\"", []l.Token{{
		Kind:l.STRING, Literal: "\"asd\\\\r\"", Value: "asd\\r", Line: 1, Col: 1, Byte: 0,
	}})
}
