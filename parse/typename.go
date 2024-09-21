package parse

import (
	"fmt"

	"github.com/stygian91/xfer/lex"
)

func TypeName(p *Parser) (Node, error) {
	token, err := p.ExpectAny([]lex.TokenKind{lex.BOOLTYPE, lex.INTTYPE, lex.FLOATTYPE, lex.STRINGTYPE, lex.IDENT})

	if err != nil {
		return Node{}, fmt.Errorf("Parse typename error: %w", err)
	}

	node := Node{Kind: TYPENAME}
	var child Node

	switch token.Kind {
	case lex.BOOLTYPE:
		child = Node{Kind: BOOLTYPE}
	case lex.INTTYPE:
		child = Node{Kind: INTTYPE}
	case lex.FLOATTYPE:
		child = Node{Kind: FLOATTYPE}
	case lex.STRINGTYPE:
		child = Node{Kind: STRINGTYPE}
	case lex.IDENT:
		child = Node{Kind: CUSTOMTYPE, Value: token.Literal}
	default:
		return Node{}, fmt.Errorf("Parse typename error: this code should be unreachalble")
	}

	node.Children = append(node.Children, child)
	return node, nil
}
