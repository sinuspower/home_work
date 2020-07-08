package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseInputFile(t *testing.T) {
	t.Run("file not exists", func(t *testing.T) {
		err := generate("notexists.go")
		require.EqualError(t, err, "can not read file: open notexists.go: no such file or directory")
	})

	t.Run("can not parse file", func(t *testing.T) {
		f, err := os.Create("badfile")
		require.NoError(t, err)

		_, err = f.Write([]byte("bad"))
		require.NoError(t, err)
		f.Close()

		err = generate("badfile")
		require.EqualError(t, err, "can not parse file: badfile:1:1: expected 'package', found bad")

		err = os.Remove("badfile")
		require.NoError(t, err)
	})
}

func TestOutputFileCreation(t *testing.T) {
	err := os.Mkdir("testdata", 0700)
	require.NoError(t, err)

	err = os.Mkdir("foo.bar", 0700)
	require.NoError(t, err)

	bytes, err := ioutil.ReadFile("../models/models.go")
	require.NoError(t, err)

	err = ioutil.WriteFile("testdata/models", bytes, 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile("testdata/models.go", bytes, 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile("testdata/models.txt", bytes, 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile("foo.bar/models", bytes, 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile("foo.bar/models.go", bytes, 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile("foo.bar/models.txt", bytes, 0644)
	require.NoError(t, err)

	t.Run("without extension", func(t *testing.T) {
		err := generate("testdata/models")
		require.NoError(t, err)
		require.FileExists(t, "testdata/models"+outFileNameSuffix)
	})

	t.Run("with go extension", func(t *testing.T) {
		err := generate("testdata/models.go")
		require.NoError(t, err)
		require.FileExists(t, "testdata/models"+outFileNameSuffix+".go")
	})

	t.Run("with not go extension", func(t *testing.T) {
		err := generate("testdata/models.txt")
		require.NoError(t, err)
		require.FileExists(t, "testdata/models"+outFileNameSuffix+".txt")
	})

	t.Run("folder with dot", func(t *testing.T) {
		err := generate("foo.bar/models")
		require.NoError(t, err)
		require.FileExists(t, "foo.bar/models"+outFileNameSuffix)

		err = generate("foo.bar/models.go")
		require.NoError(t, err)
		require.FileExists(t, "foo.bar/models"+outFileNameSuffix+".go")

		err = generate("foo.bar/models.txt")
		require.NoError(t, err)
		require.FileExists(t, "foo.bar/models"+outFileNameSuffix+".txt")
	})
	err = os.RemoveAll("testdata")
	require.NoError(t, err)

	err = os.RemoveAll("foo.bar")
	require.NoError(t, err)
}
