package parse_test

import (
	"testing"

	"github.com/stygian91/xfer/lex"
	"github.com/stygian91/xfer/parse"
)

func TestStructParse(t *testing.T) {
	input := "struct foo {\n\n}"
	l := lex.NewLexer(lex.StrIter2(input))
	tokens, err := l.Process()
	if err != nil {
		t.Errorf("TestStructParse() lex.process error: %s", err)
		return
	}

	p := parse.NewParser(tokens)
	sNode, err := parse.Struct(&p)
	if err != nil {
		t.Errorf("TestStructParse() parse.Struct() error: %s", err)
		return
	}

	nValue, ok := sNode.Value.(parse.StructValue)
	if !ok {
		t.Errorf("TestStructParse() node value is not StructValue")
		return
	}

	if nValue.Export {
		t.Errorf("TestStructParse() expected StructValue.Expect to be false, got true")
	}

	if nValue.Ident != "foo" {
		t.Errorf("TestStructParse() expected StructValue.Ident to be 'foo', got %s", nValue.Ident)
	}
}
