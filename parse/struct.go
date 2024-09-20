package parse

import (
	"fmt"

	"github.com/stygian91/xfer/lex"
)

type StructValue struct {
	Export bool
}

var fieldFuncCallParseIters = []ParseIter{ParseFuncToIter(Ident), ParseFuncToIter(TypeName), TryParseIter(Validation, lex.LSQUARE)}

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

	return node, nil
}

var structFieldParseIters = []ParseIter{
	ExpectButSkipIter(lex.STRUCT),
	ParseFuncToIter(Ident),
	NewParseListIter(StructField, lex.LCURLY, lex.RCURLY, lex.SEMICOLON),
}

func Struct(p *Parser) (Node, error) {
	node := Node{Kind: STRUCT, Value: StructValue{}}

	children, err := p.ParseSeq(structFieldParseIters)
	if err != nil {
		return structErr(err)
	}
	node.Children = children

	return node, nil
}
