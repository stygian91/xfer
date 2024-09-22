package lex

import (
	"strings"
)

func isValidFirstIdentRune(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isValidIdentRune(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_'
}

// assumes input is valid (checked by lex.handleString())
func unescape(str string) string {
	var builder strings.Builder
	escNext := false

	for _, c := range str {
		switch {
		case !escNext && c == '\\':
			escNext = true
		case escNext:
			var r rune
			switch c {
			case 'r':
				r = '\r'
			case 'n':
				r = '\n'
			case '\\':
				r = '\\'
			default:
				panic("Unexpected escape sequence")
			}
			builder.WriteRune(r)
			escNext = false
		default:
			builder.WriteRune(c)
		}
	}

	return builder.String()
}
