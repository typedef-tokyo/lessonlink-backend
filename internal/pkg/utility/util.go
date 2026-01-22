package utility

import (
	"strings"
	"unicode"
)

func Trimmer(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}
