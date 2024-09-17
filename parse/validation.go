package parse

import "github.com/stygian91/xfer/lex"

var validationParseIters = []ParseIter{NewParseListIter(FuncCall, lex.LSQUARE, lex.RSQUARE, lex.COMMA)}

func Validation(p *Parser) (Node, error) {
	node := Node{Kind: VALIDATION}

	children, err := p.ParseSeq(validationParseIters)
	if err != nil {
		return Node{}, err
	}
	node.Children = children

	return node, nil
}
