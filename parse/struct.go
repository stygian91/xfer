package parse

import (
	"fmt"

	"github.com/stygian91/xfer/lex"
)

type StructValue struct {
	Export bool
	Ident  string
}

func Struct(parser *Parser) (Node, error) {
	node := Node{Kind: STRUCT}
	value := StructValue{}

	_, hasExport := parser.Optional(lex.EXPORT)
	value.Export = hasExport

	if _, err := parser.Expect(lex.STRUCT); err != nil {
		return Node{}, err
	}

	token, err := parser.Expect(lex.IDENT)
	if err != nil {
		return Node{}, err
	}
	value.Ident = token.Literal

	if _, err := parser.Expect(lex.LCURLY); err != nil {
		return Node{}, err
	}

	// TODO: struct fields

	if _, err := parser.Expect(lex.RCURLY); err != nil {
		return Node{}, err
	}

	node.Value = value
	return node, nil
}
