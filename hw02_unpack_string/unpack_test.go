package hw02_unpack_string //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	input    string
	expected string
	err      error
}

func TestUnpack(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			input:    "abcd",
			expected: "abcd",
		},
		{ // new case
			input:    "a1",
			expected: "a",
		},
		{ // new case
			input:    "1a",
			expected: "a",
		},
		{ // new case
			input:    "5",
			expected: "",
		},
		{ // new case
			input:    "m",
			expected: "m",
		},
		{ // new case
			input:    "m0",
			expected: "",
		},
		{ // new case
			input:    "m1",
			expected: "m",
		},
		{ // new case
			input:    "0m",
			expected: "m",
		},
		{ // new case
			input:    "m9a0d2",
			expected: "mmmmmmmmmdd",
		},
		{ // new case
			input:    "aa",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "45",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "aaa10b",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "",
			expected: "",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpackWithEscape(t *testing.T) {
	// t.Skip() // Remove if task with asterisk completed

	for _, tst := range [...]test{
		{
			input:    `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			input:    `qwe\45`,
			expected: `qwe44444`,
		},
		{
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			input:    `qwe\\\3`,
			expected: `qwe\3`,
		},
		{ // new case
			input:    `\\`,
			expected: `\`,
		},
		{ // new case
			input:    `5\\`,
			expected: `\`,
		},
		{ // new case
			input:    `55\\`,
			expected: ``,
			err:      ErrInvalidString,
		},
		{ // new case
			input:    `\5\5\`,
			expected: `55`,
		},
		{ // new case
			input:    `\55`,
			expected: `55555`,
		},
		{ // new case
			input:    `\\\\\\\\\\`,
			expected: `\\\\\`,
		},
		{ // new case
			input:    `\\\\5\\\\\\`,
			expected: `\\\\\\\\\`,
		},
		{ // new case
			input:    `\\\\\\\5\\\`,
			expected: `\\\5\`,
		},
		{ // new case
			input:    `a4b\55c2d5e\10`,
			expected: `aaaab55555ccddddde`,
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}
