package parse

import (
	"fmt"

	"github.com/stygian91/xfer/lex"
)

var funcCallStartKinds = []lex.TokenKind{lex.IDENT, lex.LPAREN}
var funcArgKinds = []lex.TokenKind{lex.INT, lex.FLOAT, lex.TRUE, lex.FALSE, lex.STRING, lex.IDENT}

func funcCallErr(err error) (Node, error) {
	return Node{}, fmt.Errorf("Parse function call error: %w", err)
}

func FuncArg(p *Parser) (Node, error) {
	token, err := p.ExpectAny(funcArgKinds)
	if err != nil {
		return funcCallErr(err)
	}

	var value interface{}
	var kind NodeKind

	switch token.Kind {
	case lex.TRUE:
		kind = BOOL
		value = true
	case lex.FALSE:
		kind = BOOL
		value = false
	case lex.INT:
		kind = INT
		value = token.Value
	case lex.FLOAT:
		kind = FLOAT
		value = token.Value
	case lex.STRING:
		kind = STRING
		value = token.Value
	case lex.IDENT:
		kind = IDENT
		value = IdentValue{Name: token.Literal}
	default:
		panic("FuncArg(): this code should be unreachable")
	}

	return Node{Kind: kind, Value: value}, nil
}

func FuncCall(p *Parser) (Node, error) {
	node := Node{Kind: FUNC_CALL}

	startTokens, err := p.ExpectSeq(funcCallStartKinds)
	if err != nil {
		return funcCallErr(err)
	}

	node.Children = append(node.Children, NewIdent(startTokens[0].Literal))

	if _, exists := p.Optional(lex.RPAREN); exists {
		return node, nil
	}

	firstArg, err := FuncArg(p)
	if err != nil {
		return funcCallErr(err)
	}

	node.Children = append(node.Children, firstArg)

	for {
		if _, exists := p.Optional(lex.COMMA); !exists {
			if _, err := p.Expect(lex.RPAREN); err != nil {
				return funcCallErr(err)
			}
			break
		}

		arg, err := FuncArg(p)
		if err != nil {
			return funcCallErr(err)
		}

		node.Children = append(node.Children, arg)
	}

	return node, nil
}
