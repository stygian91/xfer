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

func structErr(err error) (Node, error) {
	return Node{}, fmt.Errorf("Parse struct error: %w", err)
}

func StructField(p *Parser) (Node, error) {
	node := Node{Kind: FIELD}

	children, err := p.ParseSeq([]ParseIter{ParseFuncToIter(Ident), ParseFuncToIter(TypeName)})
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
	node := Node{Kind: STRUCT}
	value := StructValue{}

	startTokens, err := parser.ExpectSeq(structStartKinds)
	if err != nil {
		return structErr(err)
	}

	node.Children = append(node.Children, NewIdent(startTokens[1].Literal))

	fields, err := parser.ParseSeq([]ParseIter{CreateParseList(StructField, lex.LCURLY, lex.RCURLY, lex.SEMICOLON)})
	if err != nil {
		return structErr(err)
	}
	node.Children = slices.Concat(node.Children, fields)

	node.Value = value
	return node, nil
}
