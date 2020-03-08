package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
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

	var b strings.Builder
	runes := []rune(in)
	n := len(runes)
	escape := false

	for i := 0; i < n-1; i++ {
		if unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+1]) && !escape { //numbers isn't allowed
			return "", ErrInvalidString
		}
		if runes[i] == runes[i+1] && runes[i] != '\\' && !escape { // duplicated chars isn't allowed, except `\`
			return "", ErrInvalidString
		}
		if escape && unicode.IsDigit(runes[i]) && !unicode.IsDigit(runes[i+1]) { // write single digit
			b.WriteRune(runes[i])
		}
		if escape = runes[i] == '\\' && !escape; escape {
			continue
		}
		if unicode.IsDigit(runes[i+1]) {
			count, _ := strconv.Atoi(string(runes[i+1]))
			if count != 0 {
				b.WriteString(strings.Repeat(string(runes[i]), count))
			}
		} else if !unicode.IsDigit(runes[i]) { // write single char
			b.WriteRune(runes[i])
		}
	}
	if escape { // write last
		b.WriteRune(runes[n-1])
	}
	if !unicode.IsDigit(runes[n-1]) && runes[n-1] != '\\' { // write last
		b.WriteRune(runes[n-1])
	}
	return b.String(), nil
}
