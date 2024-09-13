package parse

import "github.com/stygian91/xfer/lex"

type IdentValue struct {
	Name string
}

func NewIdent(name string) Node {
	return Node{Kind: IDENT, Value: IdentValue{Name: name}}
}

func Ident(parser *Parser) (Node, error) {
	token, err := parser.Expect(lex.IDENT)
	if err != nil {
		return Node{}, err
	}

	return NewIdent(token.Literal), nil
}
