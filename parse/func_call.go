package parse

import (
	"fmt"

	"github.com/stygian91/xfer/lex"
)

var funcCallStartKinds = []lex.TokenKind{lex.IDENT, lex.LPAREN}
var funcCallLoopKinds = []lex.TokenKind{lex.INT, lex.FLOAT, lex.TRUE, lex.FALSE, lex.STRING, lex.IDENT, lex.RPAREN}

func FuncCall(p *Parser) (Node, error) {
	node := Node{Kind: FUNC_CALL}

	startTokens, err := p.ExpectSeq(funcCallStartKinds)
	if err != nil {
		return Node{}, fmt.Errorf("Parse function call error: %w", err)
	}

	node.Children = append(node.Children, NewIdent(startTokens[0].Literal))

	for {
		token, err := p.ExpectAny(funcCallLoopKinds)
		if err != nil {
			return Node{}, fmt.Errorf("Parse function call error: %w", err)
		}

		// TODO: handle commas
		// TODO: create node based on token type

		if token.Kind == lex.RPAREN {
			break
		}
	}

	return node, nil
}
