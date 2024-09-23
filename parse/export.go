package parse

import (
	"fmt"

	"github.com/stygian91/xfer/lex"
)

var exportParseIters = []ParseIter{
	ExpectButSkipIter(lex.EXPORT),
	ParseAlternativeIters([]ParseAlternative{
		{Token: lex.STRUCT, Parse: Struct},
		// TODO: change this to expect functions as well
	}),
}

func Export(p *Parser) (Node, error) {
	nodes, err := p.ParseSeq(exportParseIters)
	if err != nil {
		return Node{}, err
	}

	if len(nodes) != 1 {
		return Node{}, fmt.Errorf("Unexpected number of nodes in Export(): %d", len(nodes))
	}

	child := nodes[0]

	switch child.Kind {
	case STRUCT:
		v := child.Value.(StructValue)
		v.Export = true
		child.Value = v
	// TODO: add function parsing
	default:
		panic("Export(): this should be unreachable.")
	}

	return child, nil
}
