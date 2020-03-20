package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

// Unpack makes basic unpacking of given string.
// This function allows skipping characters by '\'.
func Unpack(in string) (string, error) {
	if in == "" {
		return in, nil
	}

	runes := []rune(in)

	if len(runes) == 1 && !unicode.IsDigit(runes[0]) && runes[0] != '\\' {
		return in, nil
	}

	var b strings.Builder
	var esc bool
	prev := runes[0]
	for _, cur := range runes[1:] {
		if !esc && ((unicode.IsDigit(prev) && unicode.IsDigit(cur)) ||
			(prev == cur && prev != '\\')) {
			// numbers and duplicated chars aren't allowed
			return "", ErrInvalidString
		}
		if esc && unicode.IsDigit(prev) && !unicode.IsDigit(cur) { // write single digit
			b.WriteRune(prev)
		}
		if esc = prev == '\\' && !esc; esc {
			prev = cur
			continue
		}
		if unicode.IsDigit(cur) {
			cnt := int(cur - '0')
			b.WriteString(strings.Repeat(string(prev), cnt))
		} else if !unicode.IsDigit(prev) { // write single char
			b.WriteRune(prev)
		}
		prev = cur
	}
	if esc || (!unicode.IsDigit(prev) && prev != '\\') { // write last
		b.WriteRune(prev)
	}
	return b.String(), nil
}
