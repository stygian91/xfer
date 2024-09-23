package parse_test

import (
	"testing"

	i "github.com/stygian91/iter-go"
	"github.com/stygian91/xfer/lex"
	p "github.com/stygian91/xfer/parse"
	"github.com/stygian91/xfer/test"
)

func TestExport(t *testing.T) {
	input := `export struct Foo { x int; }`

	l := lex.NewLexer(i.StrRuneIter2(input))
	tokens, err := l.Process()
	if err != nil {
		t.Error(err)
		return
	}

	parser := p.NewParser(tokens)
	actual, err := p.Export(&parser)
	if err != nil {
		t.Error(err)
		return
	}

	expected := p.Node{
		Kind:  p.STRUCT,
		Value: p.StructValue{Export: true},
		Children: []p.Node{
			{Kind: p.IDENT, Value: p.IdentValue{Name: "Foo"}},
			{Kind: p.FIELD, Children: []p.Node{
				{Kind: p.IDENT, Value: p.IdentValue{Name: "x"}},
				{Kind: p.TYPENAME, Children: []p.Node{{Kind: p.INTTYPE}}},
			}},
		},
	}

	test.CheckDiff(t, expected, actual)
}
