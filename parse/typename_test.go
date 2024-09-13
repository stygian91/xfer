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
		{Kind: p.TYPENAME, Value: p.TypenameValue{Name: "int"}},
		{Kind: p.TYPENAME, Value: p.TypenameValue{Name: "float"}},
		{Kind: p.TYPENAME, Value: p.TypenameValue{Name: "string"}},
		{Kind: p.TYPENAME, Value: p.TypenameValue{Name: "bool"}},
		{Kind: p.TYPENAME, Value: p.TypenameValue{Name: "mytype"}},
	}

	test.CheckDiff(t, expected, actual)
}
