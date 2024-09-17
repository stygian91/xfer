package parse

import (
	"fmt"

	"github.com/stygian91/xfer/lex"
)

var funcCallStartKinds = []lex.TokenKind{lex.IDENT, lex.LPAREN}
var funcArgKinds = []lex.TokenKind{lex.INT, lex.FLOAT, lex.TRUE, lex.FALSE, lex.STRING, lex.IDENT}
var funcCallIters = []ParseIter{ParseFuncToIter(Ident), NewParseListIter(FuncArg, lex.LPAREN, lex.RPAREN, lex.COMMA)}

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

	children, err := p.ParseSeq(funcCallIters)
	if err != nil {
		return funcCallErr(err)
	}

	node.Children = children
	return node, nil
}
