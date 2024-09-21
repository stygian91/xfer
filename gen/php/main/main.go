package main

import (
	"fmt"

	"github.com/stygian91/xfer/gen/php"
	"github.com/stygian91/xfer/lex"
	"github.com/stygian91/xfer/parse"
)

func pan(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	input := `
	struct Foo {
		x int;
		y float;
		z Bar;
	}
	`
	it := lex.StrIter2(input)
	lexer := lex.NewLexer(it)
	tokens, err := lexer.Process()
	pan(err)

	parser := parse.NewParser(tokens)
	pStruct, err := parse.Struct(&parser)
	pan(err)

	code, err := php.Struct(pStruct)
	pan(err)
	fmt.Println(code)
}
