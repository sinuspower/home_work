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
		require.EqualError(t, err, "can not open directory: open testdata/bad: no such file or directory")
	})

	t.Run("empty directory", func(t *testing.T) {
		err := os.Mkdir("testdata/empty", 0644)
		require.Nil(t, err)

		env, err := ReadDir("testdata/empty")
		require.Nil(t, env)
		require.Nil(t, err)

		err = os.Remove("testdata/empty")
		require.Nil(t, err)
	})

	t.Run("= in file name", func(t *testing.T) {
		f, err := os.Create("testdata/env/BAD=bad")
		require.Nil(t, err)
		f.Close()

		env, err := ReadDir("testdata/env")
		require.Nil(t, env)
		require.Equal(t, ErrWrongFileName, err)

		err = os.Remove("testdata/env/BAD=bad")
		require.Nil(t, err)
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

	t.Run("whitespaces", func(t *testing.T) {
		err := os.Mkdir("testdata/whitespaces", 0700)
		require.Nil(t, err)
		tabs, err := os.Create("testdata/whitespaces/TABS")
		require.Nil(t, err)
		spaces, err := os.Create("testdata/whitespaces/SPACES")
		require.Nil(t, err)
		combined, err := os.Create("testdata/whitespaces/COMBINED")
		require.Nil(t, err)

		_, err = tabs.Write([]byte(`tabs		`))
		require.Nil(t, err)
		tabs.Close()
		_, err = spaces.Write([]byte(`spaces  `))
		require.Nil(t, err)
		spaces.Close()
		_, err = combined.Write([]byte(`combined  	`))
		require.Nil(t, err)
		combined.Close()

		expected := Environment{
			"TABS":     "tabs",
			"SPACES":   "spaces",
			"COMBINED": "combined",
		}

		actual, err := ReadDir("testdata/whitespaces")

		require.Nil(t, err)
		require.Equal(t, expected, actual)

		err = os.RemoveAll("testdata/whitespaces")
		require.Nil(t, err)
	})
}
