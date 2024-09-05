package lex

type TokenKind int
type TokenSubkind int

type Token struct {
	byte    uint
	line    uint
	col     uint
	kind    TokenKind
	subkind TokenSubkind
	literal string
	value   interface{}
}

func (this Token) Byte() uint            { return this.byte }
func (this Token) Line() uint            { return this.line }
func (this Token) Col() uint             { return this.col }
func (this Token) Kind() TokenKind       { return this.kind }
func (this Token) Subkind() TokenSubkind { return this.subkind }
func (this Token) Literal() string       { return this.literal }
func (this Token) Value() interface{}    { return this.value }

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
