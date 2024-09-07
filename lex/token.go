package lex

import "strings"

type TokenKind int

type Token struct {
	Byte    uint
	Line    uint
	Col     uint
	Kind    TokenKind
	Literal string
	Value   interface{}
}

func KindString(kind TokenKind) string {
	switch kind {
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LSQUARE:
		return "LSQUARE"
	case RSQUARE:
		return "RSQUARE"
	case LCURLY:
		return "LCURLY"
	case RCURLY:
		return "RCURLY"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case ASTERISK:
		return "ASTERISK"
	case SLASH:
		return "SLASH"
	case DOT:
		return "DOT"
	case COMMA:
		return "COMMA"
	case EQUAL:
		return "EQUAL"
	case LT:
		return "LT"
	case GT:
		return "GT"
	case BANG:
		return "BANG"
	case IDENT:
		return "IDENT"
	case STRUCT:
		return "STRUCT"
	case ENUM:
		return "ENUM"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case EXPORT:
		return "EXPORT"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case STRING:
		return "STRING"
	case BOOLTYPE:
		return "BOOLTYPE"
	case STRINGTYPE:
		return "STRINGTYPE"
	case INTTYPE:
		return "INTTYPE"
	case FLOATTYPE:
		return "FLOATTYPE"
	default:
		return "NILKIND"
	}
}

func KindsString(kinds []TokenKind) string {
	sb := strings.Builder{}

	for i, kind := range kinds {
		sb.WriteString(KindString(kind))
		if i < len(kinds)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
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
	STRUCT
	ENUM
	IF
	ELSE
	EXPORT
	TRUE
	FALSE
	INT
	FLOAT
	STRING
	BOOLTYPE
	STRINGTYPE
	INTTYPE
	FLOATTYPE
)

var reservedWords = map[string]TokenKind{
	"struct": STRUCT,
	"enum":   ENUM,
	"if":     IF,
	"else":   ELSE,
	"export": EXPORT,
	"true":   TRUE,
	"false":  FALSE,
	"bool":   BOOLTYPE,
	"string": STRINGTYPE,
	"int":    INTTYPE,
	"float":  FLOATTYPE,
}

var strEscapable = []rune{'"', '\\', 'r', 'n'}
