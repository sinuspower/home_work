package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("directory not exists", func(t *testing.T) {
		env, err := ReadDir("testdata/bad")
		require.Nil(t, env, "env is not nil")
		require.Equal(t, err, ErrCanNotOpenDir)
	})

	t.Run("empty directory", func(t *testing.T) {
		if err := os.Mkdir("testdata/empty", 0644); err != nil {
			t.Fatal(err)
		}

		env, err := ReadDir("testdata/empty")
		require.Nil(t, env)
		require.Nil(t, err)

		if err := os.Remove("testdata/empty"); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("= in file name", func(t *testing.T) {
		f, err := os.Create("testdata/env/BAD=bad")
		if err != nil {
			t.Fatal(err)
		}
		f.Close()

		env, err := ReadDir("testdata/env")
		require.Nil(t, env)
		require.Equal(t, ErrWrongFileName, err)

		if err := os.Remove("testdata/env/BAD=bad"); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("positive case", func(t *testing.T) {
		expected := Environment{
			"BAR": `bar`,
			"FOO": `   foo
with new line`,
			"HELLO": `"hello"`,
			"UNSET": ``,
		}

		actual, err := ReadDir("testdata/env")

		require.Nil(t, err)
		require.Equal(t, expected, actual)
	})
}
