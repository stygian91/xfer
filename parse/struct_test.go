package parse_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stygian91/xfer/lex"
	p "github.com/stygian91/xfer/parse"
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
		t.Errorf("TestStructParse() lex.process error: %s", err)
		return
	}

	parser := p.NewParser(tokens)
	actual1, err := p.Struct(&parser)
	if err != nil {
		t.Errorf("TestStructParse() parse.Struct() error: %s", err)
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

	if diff := cmp.Diff(expected1, actual1); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}

	actual2, err := p.Struct(&parser)
	if err != nil {
		t.Errorf("TestStructParse() parse.Struct() error: %s", err)
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

	if diff := cmp.Diff(expected2, actual2); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
