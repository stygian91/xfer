package parse

import (
	"iter"

	"github.com/stygian91/xfer/lex"
)

func ParseFuncToIter(fn ParseFunc) ParseIter {
	return func(p *Parser) iter.Seq[ParseRes] {
		return func(yield func(ParseRes) bool) {
			node, err := fn(p)
			yield(ParseRes{Value: node, Err: err})
		}
	}
}

func ParseFuncToOptionalIter(fn ParseFunc) ParseIter {
	return func(p *Parser) iter.Seq[ParseRes] {
		return func(yield func(ParseRes) bool) {
			node, err := fn(p)
			if err != nil {
				yield(ParseRes{Value: Node{Kind: NILKIND}, Err: nil})
				return
			}

			yield(ParseRes{Value: node, Err: err})
		}
	}
}

type ParseFunc func(*Parser) (Node, error)

type ParseRes struct {
	Value Node
	Err   error
}

type ParseIter func(*Parser) iter.Seq[ParseRes]

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
