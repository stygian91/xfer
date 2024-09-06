package lex

type TokenKind int
type TokenSubkind int

type Token struct {
	Byte    uint
	Line    uint
	Col     uint
	Kind    TokenKind
	Subkind TokenSubkind
	Literal string
	Value   interface{}
}

const NILKIND = 0

// tokens
const (
	LPAREN = iota + 1
	RPAREN
	LSQUARE
	RSQUARE
	LCURLY
	RCURLY
	PLUS
	MINUS
	ASTERISK
	SLASH
	DOT
	COMMA
	EQUAL
	LT
	GT
	BANG
	IDENT
	KEYWORD
	INT
	FLOAT
	STRING
)

// keywords
const (
	STRUCT = iota + 1
	ENUM
	IF
	ELSE
	EXPORT
	TRUE
	FALSE
)

var reservedWords = map[string]TokenSubkind{"struct": STRUCT, "enum": ENUM, "if": IF, "else": ELSE, "export": EXPORT, "true": TRUE, "false": FALSE}

var strEscapable = []rune{'"', '\\', 'r', 'n'}
