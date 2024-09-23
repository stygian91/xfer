package parse

import (
	"fmt"
	"iter"

	"github.com/stygian91/xfer/lex"
)

type ParseFunc func(*Parser) (Node, error)

type ParseRes struct {
	Value Node
	Err   error
}

type ParseIter func(*Parser) iter.Seq[ParseRes]

type ParseAlternative struct {
	Token lex.TokenKind
	Parse ParseFunc
}

func ParseAlternativeIters(alternatives []ParseAlternative) ParseIter {
	getExpectedKinds := func() []lex.TokenKind {
		kinds := []lex.TokenKind{}
		for _, alt := range alternatives {
			kinds = append(kinds, alt.Token)
		}
		return kinds
	}

	return func(p *Parser) iter.Seq[ParseRes] {
		return func(yield func(ParseRes) bool) {
			for _, alt := range alternatives {
				if p.CurrentTokenIs(alt.Token) {
					node, err := alt.Parse(p)
					yield(ParseRes{Value: node, Err: err})
					return
				}
			}

			currToken, _ := p.CurrentToken()

			yield(ParseRes{
				Value: Node{},
				Err:   fmt.Errorf("Parse alternatives failed: expected one of [%s], got %s", lex.KindsString(getExpectedKinds()), lex.KindString(currToken.Kind)),
			})
		}
	}
}

func ParseFuncToIter(fn ParseFunc) ParseIter {
	return func(p *Parser) iter.Seq[ParseRes] {
		return func(yield func(ParseRes) bool) {
			node, err := fn(p)
			yield(ParseRes{Value: node, Err: err})
		}
	}
}

func TryParseIter(fn ParseFunc, start lex.TokenKind) ParseIter {
	return func(p *Parser) iter.Seq[ParseRes] {
		return func(yield func(ParseRes) bool) {
			if !p.CurrentTokenIs(start) {
				return
			}

			node, err := fn(p)
			yield(ParseRes{Value: node, Err: err})
		}
	}
}

func ExpectButSkipIter(kind lex.TokenKind) ParseIter {
	return func(p *Parser) iter.Seq[ParseRes] {
		return func(yield func(ParseRes) bool) {
			if _, err := p.Expect(kind); err != nil {
				yield(ParseRes{Value: Node{}, Err: err})
			}
		}
	}
}

func NewParseListIter(
	parseEl ParseFunc,
	start, end, delimiter lex.TokenKind,
) ParseIter {
	return func(p *Parser) iter.Seq[ParseRes] {
		return func(yield func(ParseRes) bool) {
			_, err := p.Expect(start)
			if err != nil {
				yield(ParseRes{Err: err})
				return
			}

			if _, exists := p.Optional(end); exists {
				return
			}

			firstEl, err := parseEl(p)
			if err != nil {
				yield(ParseRes{Err: err})
				return
			}

			if !yield(ParseRes{Value: firstEl}) {
				return
			}

			for {
				_, delimiterExists := p.Optional(delimiter)
				if delimiterExists {
					if _, endExists := p.Optional(end); endExists {
						break
					}
				} else {
					if _, err := p.Expect(end); err != nil {
						yield(ParseRes{Err: err})
						return
					}
					break
				}

				el, err := parseEl(p)
				if err != nil {
					yield(ParseRes{Err: err})
					return
				}

				if !yield(ParseRes{Value: el}) {
					return
				}
			}
		}
	}
}
