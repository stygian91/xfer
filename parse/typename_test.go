package parse_test

import (
	"testing"

	"github.com/stygian91/xfer/lex"
	p "github.com/stygian91/xfer/parse"
	"github.com/stygian91/xfer/test"
)

func TestTypenameParse(t *testing.T) {
	input := `int float string bool mytype`

	l := lex.NewLexer(lex.StrIter2(input))
	tokens, err := l.Process()
	if err != nil {
		t.Error(err)
		return
	}

	parser := p.NewParser(tokens)
	actual := []p.Node{}
	for i := 0; i < 5; i++ {
		node, err := p.TypeName(&parser)
		if err != nil {
			t.Error(err)
		}

		actual = append(actual, node)
	}

	expected := []p.Node{
		{Kind: p.TYPENAME, Children: []p.Node{{Kind: p.INTTYPE}}},
		{Kind: p.TYPENAME, Children: []p.Node{{Kind: p.FLOATTYPE}}},
		{Kind: p.TYPENAME, Children: []p.Node{{Kind: p.STRINGTYPE}}},
		{Kind: p.TYPENAME, Children: []p.Node{{Kind: p.BOOLTYPE}}},
		{Kind: p.TYPENAME, Children: []p.Node{{Kind: p.CUSTOMTYPE, Value: "mytype"}}},
	}

	test.CheckDiff(t, expected, actual)
}
