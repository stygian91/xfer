package parse_test

import (
	"testing"

	"github.com/stygian91/xfer/lex"
	p "github.com/stygian91/xfer/parse"
	"github.com/stygian91/xfer/test"
)

func TestStructParse(t *testing.T) {
	input := `struct foo {
		x int;
		y float; z mytype;
	}

	struct bar { baz string; }
	`

	l := lex.NewLexer(lex.StrIter2(input))
	tokens, err := l.Process()
	if err != nil {
		t.Error(err)
		return
	}

	parser := p.NewParser(tokens)
	actual1, err := p.Struct(&parser)
	if err != nil {
		t.Error(err)
		return
	}

	expected1 := p.Node{
		Kind:  p.STRUCT,
		Value: p.StructValue{Export: false},
		Children: []p.Node{
			{Kind: p.IDENT, Value: p.IdentValue{Name: "foo"}},
			{Kind: p.FIELD, Children: []p.Node{
				{Kind: p.IDENT, Value: p.IdentValue{Name: "x"}},
				{Kind: p.TYPENAME, Value: p.TypenameValue{Name: "int"}},
			}},
			{Kind: p.FIELD, Children: []p.Node{
				{Kind: p.IDENT, Value: p.IdentValue{Name: "y"}},
				{Kind: p.TYPENAME, Value: p.TypenameValue{Name: "float"}},
			}},
			{Kind: p.FIELD, Children: []p.Node{
				{Kind: p.IDENT, Value: p.IdentValue{Name: "z"}},
				{Kind: p.TYPENAME, Value: p.TypenameValue{Name: "mytype"}},
			}},
		},
	}

	test.CheckDiff(t, expected1, actual1)

	actual2, err := p.Struct(&parser)
	if err != nil {
		t.Error(err)
		return
	}

	expected2 := p.Node{
		Kind:  p.STRUCT,
		Value: p.StructValue{Export: false},
		Children: []p.Node{
			{Kind: p.IDENT, Value: p.IdentValue{Name: "bar"}},
			{Kind: p.FIELD, Children: []p.Node{
				{Kind: p.IDENT, Value: p.IdentValue{Name: "baz"}},
				{Kind: p.TYPENAME, Value: p.TypenameValue{Name: "string"}},
			}},
		},
	}

	test.CheckDiff(t, expected2, actual2)
}
