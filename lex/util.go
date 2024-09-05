package lex

import (
	stditer "iter"
	"strings"
)

func isValidFirstIdentRune(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isValidIdentRune(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_'
}

func StrIter(input string) stditer.Seq[rune] {
	return func(yield func(rune) bool) {
		for _, char := range input {
			if !yield(char) {
				return
			}
		}
	}
}

func StrIter2(input string) stditer.Seq2[int, rune] {
	return func(yield func(int, rune) bool) {
		for i, char := range input {
			if !yield(i, char) {
				return
			}
		}
	}
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
