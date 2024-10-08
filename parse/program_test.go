package parse_test

import (
	"testing"

	i "github.com/stygian91/iter-go"
	"github.com/stygian91/xfer/lex"
	p "github.com/stygian91/xfer/parse"
	"github.com/stygian91/xfer/test"
)

func TestProgramParse(t *testing.T) {
	input := `struct foo {
		x int;
		y float; z mytype;
	}

	struct bar { baz string }
	struct baz {}
	`

	l := lex.NewLexer(i.StrRuneIter2(input))
	tokens, err := l.Process()
	if err != nil {
		t.Error(err)
		return
	}

	parser := p.NewParser(tokens)
	actual, err := p.Program(&parser)
	if err != nil {
		t.Error(err)
		return
	}

	expected := p.Node{
		Kind: p.PROGRAM,
		Children: []p.Node{
			{
				Kind:  p.STRUCT,
				Value: p.StructValue{Export: false},
				Children: []p.Node{
					{Kind: p.IDENT, Value: p.IdentValue{Name: "foo"}},
					{Kind: p.FIELD, Children: []p.Node{
						{Kind: p.IDENT, Value: p.IdentValue{Name: "x"}},
						{Kind: p.TYPENAME, Children: []p.Node{{Kind: p.INTTYPE}}},
					}},
					{Kind: p.FIELD, Children: []p.Node{
						{Kind: p.IDENT, Value: p.IdentValue{Name: "y"}},
						{Kind: p.TYPENAME, Children: []p.Node{{Kind: p.FLOATTYPE}}},
					}},
					{Kind: p.FIELD, Children: []p.Node{
						{Kind: p.IDENT, Value: p.IdentValue{Name: "z"}},
						{Kind: p.TYPENAME, Children: []p.Node{{Kind: p.CUSTOMTYPE, Value: "mytype"}}},
					}},
				},
			},
			{
				Kind:  p.STRUCT,
				Value: p.StructValue{Export: false},
				Children: []p.Node{
					{Kind: p.IDENT, Value: p.IdentValue{Name: "bar"}},
					{Kind: p.FIELD, Children: []p.Node{
						{Kind: p.IDENT, Value: p.IdentValue{Name: "baz"}},
						{Kind: p.TYPENAME, Children: []p.Node{{Kind: p.STRINGTYPE}}},
					}},
				},
			},
			{
				Kind:  p.STRUCT,
				Value: p.StructValue{Export: false},
				Children: []p.Node{
					{Kind: p.IDENT, Value: p.IdentValue{Name: "baz"}},
				},
			},
		},
	}

	test.CheckDiff(t, expected, actual)
}

func TestStructParseValidation(t *testing.T) {
	input := `struct foo {
		x int;
		y float [];
		z bool [val1(), val2(42), val3(42, "str")]
	}
	`

	l := lex.NewLexer(i.StrRuneIter2(input))
	tokens, err := l.Process()
	if err != nil {
		t.Error(err)
		return
	}

	parser := p.NewParser(tokens)
	actual, err := p.Struct(&parser)
	if err != nil {
		t.Error(err)
		return
	}

	expected := p.Node{
		Kind:  p.STRUCT,
		Value: p.StructValue{Export: false},
		Children: []p.Node{
			{Kind: p.IDENT, Value: p.IdentValue{Name: "foo"}},
			{Kind: p.FIELD, Children: []p.Node{
				{Kind: p.IDENT, Value: p.IdentValue{Name: "x"}},
				{Kind: p.TYPENAME, Children: []p.Node{{Kind: p.INTTYPE}}},
			}},
			{Kind: p.FIELD, Children: []p.Node{
				{Kind: p.IDENT, Value: p.IdentValue{Name: "y"}},
				{Kind: p.TYPENAME, Children: []p.Node{{Kind: p.FLOATTYPE}}},
				{Kind: p.VALIDATION, Children: []p.Node{}},
			}},
			{Kind: p.FIELD, Children: []p.Node{
				{Kind: p.IDENT, Value: p.IdentValue{Name: "z"}},
				{Kind: p.TYPENAME, Children: []p.Node{{Kind: p.BOOLTYPE}}},
				{Kind: p.VALIDATION, Children: []p.Node{
					{Kind: p.FUNC_CALL, Children: []p.Node{
						p.NewIdent("val1"),
					}},
					{Kind: p.FUNC_CALL, Children: []p.Node{
						p.NewIdent("val2"),
						{Kind: p.INT, Value: int64(42)},
					}},
					{Kind: p.FUNC_CALL, Children: []p.Node{
						p.NewIdent("val3"),
						{Kind: p.INT, Value: int64(42)},
						{Kind: p.STRING, Value: "str"},
					}},
				}},
			}},
		},
	}

	test.CheckDiff(t, expected, actual)
}

func TestStructParseErrors(t *testing.T) {
	inputs := []string{
		`struct {}`,
		`struct foo {`,
		`struct foo { x; }`,
		`struct foo { x int [asd(); }`,
	}

	for _, input := range inputs {
		l := lex.NewLexer(i.StrRuneIter2(input))
		tokens, err := l.Process()
		if err != nil {
			t.Error(err)
			return
		}

		parser := p.NewParser(tokens)
		_, err = p.Struct(&parser)
		if err == nil {
			t.Errorf("Expected an error when parsing '%s'", input)
		}
	}
}
