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
		return "", nil
	}

	runes := []rune(in)
	var b strings.Builder
	var esc bool // previous rune has been escaped by '/'
	var prev rune

	for i, cur := range runes {
		if (i == 0 && i != len(runes)-1) || (cur == '\\' && prev != '\\') {
			if esc {
				b.WriteRune(prev)
				esc = false
			} else if !unicode.IsDigit(prev) && prev != 0 {
				b.WriteRune(prev)
			}
			prev = cur
			continue
		}
		if !esc && unicode.IsDigit(prev) && unicode.IsDigit(cur) { // numbers aren't allowed
			return "", ErrInvalidString
		}
		if !esc && prev != '\\' && cur == prev { //duplicated chars aren't allowed
			return "", ErrInvalidString
		}
		if unicode.IsDigit(cur) && (prev != '\\' || esc) && prev != 0 { // repeat
			cnt := int(cur - '0')
			if cnt != 0 {
				b.WriteString(strings.Repeat(string(prev), cnt))
			}
		} else { // write single char
			if esc { // any
				b.WriteRune(prev)
			}
			if !unicode.IsDigit(prev) && prev != '\\' && prev != 0 { // letter
				b.WriteRune(prev)
			}
		}
		if i == len(runes)-1 { // write last
			if prev == '\\' && !esc { // any
				b.WriteRune(cur)
			} else if !unicode.IsDigit(cur) && cur != '\\' { // letter
				b.WriteRune(cur)
			}
		}
		esc = !esc && prev == '\\'
		prev = cur
	}
	return b.String(), nil
}
