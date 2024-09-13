package parse

import (
	"fmt"

	"github.com/stygian91/xfer/lex"
)

type StructValue struct {
	Export bool
}

var structStartKinds = []lex.TokenKind{lex.STRUCT, lex.IDENT}

func structErr(err error) error {
	return fmt.Errorf("Parse struct error: %w", err)
}

func Struct(parser *Parser) (Node, error) {
	node := Node{Kind: STRUCT}
	value := StructValue{}

	startTokens, err := parser.ExpectSeq(structStartKinds)
	if err != nil {
		return Node{}, structErr(err)
	}

	node.Children = append(node.Children, Node{
		Kind:  IDENT,
		Value: IdentValue{Name: startTokens[1].Literal},
	})

	if _, err := parser.Expect(lex.LCURLY); err != nil {
		return Node{}, structErr(err)
	}

	for {
		if t, exists := parser.CurrentToken(); exists && t.Kind == lex.RCURLY {
			parser.Eat()
			break
		}

		token, err := parser.Expect(lex.IDENT)
		if err != nil {
			return Node{}, structErr(err)
		}

		typename, err := TypeName(parser)
		if err != nil {
			return Node{}, structErr(err)
		}

		// TODO: validation list parsing

		_, err = parser.Expect(lex.SEMICOLON)
		if err != nil {
			return Node{}, structErr(err)
		}

		node.Children = append(node.Children, Node{
			Kind:     FIELD,
			Children: []Node{{Kind: IDENT, Value: IdentValue{Name: token.Literal}}, typename},
		})
	}

	node.Value = value
	return node, nil
}
