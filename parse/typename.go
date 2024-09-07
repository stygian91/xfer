package parse

import (
	"fmt"

	"github.com/stygian91/xfer/lex"
)

type TypenameValue struct {
	Name string
}

func TypeName(p *Parser) (Node, error) {
	token, err := p.ExpectAny([]lex.TokenKind{lex.BOOLTYPE, lex.INTTYPE, lex.FLOATTYPE, lex.STRINGTYPE, lex.IDENT})

	if err != nil {
		return Node{}, fmt.Errorf("Parse typename error: %w", err)
	}

	node := Node{Kind: TYPENAME}
	var name string

	switch token.Kind {
	case lex.BOOLTYPE:
		name = "bool"
	case lex.INTTYPE:
		name = "int"
	case lex.FLOATTYPE:
		name = "float"
	case lex.STRINGTYPE:
		name = "string"
	case lex.IDENT:
		name = token.Literal
	default:
		return Node{}, fmt.Errorf("Parse typename error: this code should be unreachalble")
	}

	node.Value = TypenameValue{name}
	return node, nil
}
