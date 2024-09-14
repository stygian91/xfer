package parse

import "github.com/stygian91/xfer/lex"

func Validation(p *Parser) (Node, error) {
	_, err := p.Expect(lex.LSQUARE)
	if err != nil {
		return Node{}, err
	}

	node := Node{Kind: VALIDATION}
	if _, exists := p.Optional(lex.RSQUARE); exists {
		return node, nil
	}

	firstCallNode, err := FuncCall(p)
	if err != nil {
		return Node{}, err
	}

	node.Children = append(node.Children, firstCallNode)

	for {
		if _, exists := p.Optional(lex.COMMA); !exists {
			if _, err := p.Expect(lex.RSQUARE); err != nil {
				return Node{}, err
			}
			break
		}

		callNode, err := FuncCall(p)
		if err != nil {
			return Node{}, err
		}

		node.Children = append(node.Children, callNode)
	}

	return node, nil
}
