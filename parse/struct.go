package parse

import (
	"fmt"

	"github.com/stygian91/xfer/lex"
)

type StructValue struct {
	Export bool
}

var structFieldParseIters = []ParseIter{ParseFuncToIter(Ident), NewParseListIter(StructField, lex.LCURLY, lex.RCURLY, lex.SEMICOLON)}
var fieldFuncCallParseIters = []ParseIter{ParseFuncToIter(Ident), ParseFuncToIter(TypeName)}

func structErr(err error) (Node, error) {
	return Node{}, fmt.Errorf("Parse struct error: %w", err)
}

func StructField(p *Parser) (Node, error) {
	node := Node{Kind: FIELD}

	children, err := p.ParseSeq(fieldFuncCallParseIters)
	if err != nil {
		return structErr(err)
	}

	node.Children = children

	if p.CurrentTokenIs(lex.LSQUARE) {
		validationNode, err := Validation(p)
		if err != nil {
			return structErr(err)
		}

		node.Children = append(node.Children, validationNode)
	}

	return node, nil
}

func Struct(p *Parser) (Node, error) {
	node := Node{Kind: STRUCT, Value: StructValue{}}

	_, err := p.Expect(lex.STRUCT)
	if err != nil {
		return structErr(err)
	}

	children, err := p.ParseSeq(structFieldParseIters)
	if err != nil {
		return structErr(err)
	}
	node.Children = children

	return node, nil
}
