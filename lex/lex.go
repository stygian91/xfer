package lex

type Lexer struct {
	input  string
	cursor uint
}

// tokens
const (
	LPAREN = iota + 1
	RPAREN
	LSQUARE
	RSQUARE
	LCURLY
	RCURLY
	IDENT
	KEYWORD
	INT
	FLOAT
	STRING
	BOOL
)

// keywords
const (
	STRUCT = "struct"
	IF     = "if"
	ELSE   = "else"
	EXPORT = "export"
)

func NewLexer(input string) Lexer {
	return Lexer{
		input:  input,
		cursor: 0,
	}
}

func (this *Lexer) Process() {
}
