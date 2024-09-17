package parse

import (
	"fmt"
	"slices"

	"github.com/stygian91/xfer/lex"
)

type StructValue struct {
	Export bool
}

var structStartKinds = []lex.TokenKind{lex.STRUCT, lex.IDENT}
var structFieldParseIters = []ParseIter{NewParseListIter(StructField, lex.LCURLY, lex.RCURLY, lex.SEMICOLON)}
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

	if t, exists := p.CurrentToken(); exists && t.Kind == lex.LSQUARE {
		validationNode, err := Validation(p)
		if err != nil {
			return structErr(err)
		}

		node.Children = append(node.Children, validationNode)
	}

	return node, nil
}

func Struct(parser *Parser) (Node, error) {
	node := Node{Kind: STRUCT, Value: StructValue{}}

	startTokens, err := parser.ExpectSeq(structStartKinds)
	if err != nil {
		return structErr(err)
	}
	node.Children = append(node.Children, NewIdent(startTokens[1].Literal))

	fields, err := parser.ParseSeq(structFieldParseIters)
	if err != nil {
		return structErr(err)
	}
	node.Children = slices.Concat(node.Children, fields)

	return node, nil
}
