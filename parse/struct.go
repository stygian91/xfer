package parse

import (
	"fmt"

	"github.com/stygian91/xfer/lex"
)

type StructValue struct {
	Export bool
}

var fieldLoopKinds = []lex.TokenKind{lex.RCURLY, lex.IDENT}
var structKinds = []lex.TokenKind{lex.STRUCT, lex.IDENT}

func structErr(err error) error {
	return fmt.Errorf("Parse struct error: %w", err)
}

func Struct(parser *Parser) (Node, error) {
	node := Node{Kind: STRUCT}
	value := StructValue{}

	startTokens, err := parser.ExpectSeq(structKinds)
	if err != nil {
		return Node{}, structErr(err)
	}

	node.Children = append(node.Children, NewIdent(startTokens[1].Literal))

	if _, err := parser.Expect(lex.LCURLY); err != nil {
		return Node{}, structErr(err)
	}

	for {
		token, err := parser.ExpectAny(fieldLoopKinds)
		if err != nil {
			return Node{}, structErr(err)
		}

		if token.Kind == lex.RCURLY {
			break
		}

		typename, err := TypeName(parser)
		if err != nil {
			return Node{}, structErr(err)
		}

		// TODO: validation list parsing

		node.Children = append(node.Children, Node{
			Kind:     FIELD,
			Children: []Node{NewIdent(token.Literal), typename},
		})
	}

	node.Value = value
	return node, nil
}
