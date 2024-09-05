package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/stygian91/xfer/lex"
)

func main() {
	l := log.New(os.Stderr, "", 0)
	flag.Parse()

	inputFilename := flag.Arg(0)
	if len(inputFilename) == 0 {
		l.Println("Input file required")
		os.Exit(1)
	}

	file, err := os.Open(inputFilename)
	if err != nil {
		l.Printf("Error opening input file: %s", err)
		os.Exit(1)
	}

	input, err := io.ReadAll(file)
	if err != nil {
		l.Printf("Error reading input file: %s", err)
		os.Exit(1)
	}

	lexer := lex.NewLexer(lex.StrIter2(string(input)))
	tokens, err := lexer.Process()
	if err != nil {
		l.Printf("Error lexing: %s", err)
		os.Exit(1)
	}

	fmt.Printf("tokens: %+v\n", tokens)
}
