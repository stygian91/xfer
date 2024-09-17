package parse_test

import (
	"testing"

	"github.com/stygian91/xfer/lex"
	p "github.com/stygian91/xfer/parse"
	"github.com/stygian91/xfer/test"
)

func TestFuncCall(t *testing.T) {
	input := `foo()foo2(myvar) foo3(42, 3.14, "str", false, true, myvar2)`
	input2 := `foo()foo2(myvar) foo3(
	42,
	3.14,
	"str",
	false,
	true,
	myvar2,
)`

	l := lex.NewLexer(lex.StrIter2(input))
	tokens, err := l.Process()
	if err != nil {
		t.Error(err)
		return
	}

	l2 := lex.NewLexer(lex.StrIter2(input2))
	tokens2, err := l2.Process()
	if err != nil {
		t.Error(err)
		return
	}

	parser := p.NewParser(tokens)
	parser2 := p.NewParser(tokens2)
	actual := []p.Node{}
	actual2 := []p.Node{}
	for i := 0; i < 3; i++ {
		node, err := p.FuncCall(&parser)
		if err != nil {
			t.Error(err)
			return
		}
		actual = append(actual, node)

		node2, err := p.FuncCall(&parser2)
		if err != nil {
			t.Error(err)
			return
		}
		actual2 = append(actual2, node2)
	}

	expected := []p.Node{
		{Kind: p.FUNC_CALL, Children: []p.Node{p.NewIdent("foo")}},
		{Kind: p.FUNC_CALL, Children: []p.Node{
			p.NewIdent("foo2"),
			p.NewIdent("myvar"),
		}},
		{Kind: p.FUNC_CALL, Children: []p.Node{
			p.NewIdent("foo3"),
			{Kind: p.INT, Value: int64(42)},
			{Kind: p.FLOAT, Value: 3.14},
			{Kind: p.STRING, Value: "str"},
			{Kind: p.BOOL, Value: false},
			{Kind: p.BOOL, Value: true},
			p.NewIdent("myvar2"),
		}},
	}

	test.CheckDiff(t, expected, actual)
	test.CheckDiff(t, expected, actual2)
}

func TestFuncCallErr(t *testing.T) {
	input := `foo(123 45)`

	l := lex.NewLexer(lex.StrIter2(input))
	tokens, err := l.Process()
	if err != nil {
		t.Error(err)
		return
	}

	parser := p.NewParser(tokens)
	_, err = p.FuncCall(&parser)
	if err == nil {
		t.Errorf("Expected error when parsing '%s'", input)
	}
}
