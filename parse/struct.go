package parse

import (
	"fmt"

	"github.com/stygian91/xfer/lex"
)

type StructValue struct {
	Export bool
}

var structStartKinds = []lex.TokenKind{lex.STRUCT, lex.IDENT, lex.LCURLY}

func structErr(err error) (Node, error) {
	return Node{}, fmt.Errorf("Parse struct error: %w", err)
}

// TODO: rework with ParseFuncs
func Struct(parser *Parser) (Node, error) {
	node := Node{Kind: STRUCT}
	value := StructValue{}

	startTokens, err := parser.ExpectSeq(structStartKinds)
	if err != nil {
		return structErr(err)
	}

	node.Children = append(node.Children, NewIdent(startTokens[1].Literal))

	for {
		if _, exists := parser.Optional(lex.RCURLY); exists {
			break
		}

		identNode, err := Ident(parser)
		if err != nil {
			return structErr(err)
		}

		typename, err := TypeName(parser)
		if err != nil {
			return structErr(err)
		}

		var validationNode Node
		hasValidationNode := false

		if t, exists := parser.CurrentToken(); exists && t.Kind == lex.LSQUARE {
			validationNode, err = Validation(parser)
			if err != nil {
				return structErr(err)
			}

			hasValidationNode = true
		}

		_, err = parser.Expect(lex.SEMICOLON)
		if err != nil {
			return structErr(err)
		}

		children := []Node{identNode, typename}
		if hasValidationNode {
			children = append(children, validationNode)
		}

		node.Children = append(node.Children, Node{
			Kind:     FIELD,
			Children: children,
		})
	}

	node.Value = value
	return node, nil
}
