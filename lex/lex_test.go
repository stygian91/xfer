package lex_test

import (
	"strings"
	"testing"

	"github.com/stygian91/iter-go"
	l "github.com/stygian91/xfer/lex"
	"github.com/stygian91/xfer/test"
)

func checkLexResults(t *testing.T, input string, expected []l.Token) {
	it := iter.StrRuneIter2(input)
	lex := l.NewLexer(it)
	tokens, err := lex.Process()

	if err != nil {
		t.Errorf("Error while running lex.Process(): %s", err)
		return
	}

	test.CheckDiff(t, expected, tokens)
}

func TestSimples(t *testing.T) {
	checkLexResults(t, "()[]{}+-*/.,=<>!;", []l.Token{
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
		{Kind: l.SEMICOLON, Literal: ";", Line: 1, Col: 17, Byte: 16},
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
		{Kind: l.IF, Literal: "if", Line: 1, Col: 5, Byte: 4},
		{Kind: l.ELSE, Literal: "else", Line: 1, Col: 8, Byte: 7},
		{Kind: l.EXPORT, Literal: "export", Line: 1, Col: 13, Byte: 12},
		{Kind: l.TRUE, Literal: "true", Line: 1, Col: 20, Byte: 19},
		{Kind: l.FALSE, Literal: "false", Line: 1, Col: 25, Byte: 24},
		{Kind: l.STRUCT, Literal: "struct", Line: 1, Col: 31, Byte: 30},
		{Kind: l.ENUM, Literal: "enum", Line: 1, Col: 38, Byte: 37},
	})
}

func TestString(t *testing.T) {
	checkLexResults(t, "\"asd\"", []l.Token{{
		Kind: l.STRING, Literal: "\"asd\"", Value: "asd", Line: 1, Col: 1, Byte: 0,
	}})

	checkLexResults(t, "\"asd\\r\"", []l.Token{{
		Kind: l.STRING, Literal: "\"asd\\r\"", Value: "asd\r", Line: 1, Col: 1, Byte: 0,
	}})

	checkLexResults(t, "\"asd\\\\r\"", []l.Token{{
		Kind: l.STRING, Literal: "\"asd\\\\r\"", Value: "asd\\r", Line: 1, Col: 1, Byte: 0,
	}})
}

func BenchmarkStrIter(b *testing.B) {
	inputs := []string{
		"()[]{}+-*/.,=<>!;",
		"123 42.69",
		"asd if else export true false struct enum",
	}

	for i := 0; i < b.N; i++ {
		for _, s := range inputs {
			it := iter.StrRuneIter2(s)
			lex := l.NewLexer(it)
			_, err := lex.Process()

			if err != nil {
				b.Error(err)
				return
			}
		}
	}
}

func BenchmarkReaderStrIter(b *testing.B) {
	inputs := []string{
		"()[]{}+-*/.,=<>!;",
		"123 42.69",
		"asd if else export true false struct enum",
	}

	for i := 0; i < b.N; i++ {
		for _, s := range inputs {
			r := strings.NewReader(s)
			it := iter.Utf8ReaderToRuneIter2(r, 4096)
			lex := l.NewLexer(it)
			_, err := lex.Process()

			if err != nil {
				b.Error(err)
				return
			}
		}
	}
}
