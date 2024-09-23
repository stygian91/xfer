package parse

import (
	"slices"

	"github.com/stygian91/xfer/lex"
)

var programParseIters = []ParseIter{
	ParseAlternativeIters([]ParseAlternative{
		{Token: lex.EXPORT, Parse: Export},
		{Token: lex.STRUCT, Parse: Struct},
		// TODO: function
	}),
}

func Program(p *Parser) (Node, error) {
	programNode := Node{Kind: PROGRAM}

	for {
		if !p.HasMore() {
			break
		}

		nodes, err := p.ParseSeq(programParseIters)
		if err != nil {
			return Node{}, err
		}
		programNode.Children = slices.Concat(programNode.Children, nodes)
	}

	return programNode, nil
}
