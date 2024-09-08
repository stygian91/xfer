package parse

import (
	"fmt"

	"github.com/stygian91/xfer/lex"
)

type StructValue struct {
	Export bool
	// Ident  string
}

type IdentValue struct {
	Name string
}

var fieldLoopTokenKinds = []lex.TokenKind{lex.RCURLY, lex.IDENT}

func wrapErr(err error) error {
	return fmt.Errorf("Parse struct error: %w", err)
}

func Struct(parser *Parser) (Node, error) {
	node := Node{Kind: STRUCT}
	value := StructValue{}

	_, hasExport := parser.Optional(lex.EXPORT)
	value.Export = hasExport

	if _, err := parser.Expect(lex.STRUCT); err != nil {
		return Node{}, wrapErr(err)
	}

	token, err := parser.Expect(lex.IDENT)
	if err != nil {
		return Node{}, wrapErr(err)
	}
	node.Children = append(node.Children, Node{
		Kind:  IDENT,
		Value: IdentValue{Name: token.Literal},
	})

	if _, err := parser.Expect(lex.LCURLY); err != nil {
		return Node{}, wrapErr(err)
	}

	for {
		token, err := parser.ExpectAny(fieldLoopTokenKinds)
		if err != nil {
			return Node{}, wrapErr(err)
		}

		if token.Kind == lex.RCURLY {
			break
		}

		typename, err := TypeName(parser)
		if err != nil {
			return Node{}, wrapErr(err)
		}

		// TODO: validation list parsing

		node.Children = append(node.Children, Node{
			Kind:     FIELD,
			Children: []Node{{Kind: IDENT, Value: IdentValue{Name: token.Literal}}, typename},
		})
	}

	node.Value = value
	return node, nil
}
