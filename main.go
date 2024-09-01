package main

import (
	"flag"
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

	lexer := lex.NewLexer(string(input))
	lexer.Process()
}
