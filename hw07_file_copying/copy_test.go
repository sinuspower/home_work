package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/udhos/equalfile"
)

func TestCopy(t *testing.T) {
	t.Run("file not found", func(t *testing.T) {
		err := Copy("test", "testCopy", 0, 0)
		require.Equal(t, err, ErrUnsupportedFile)
	})

	t.Run("zero offset and limit", func(t *testing.T) {
		inPath := "testdata/input.txt"
		outPath := "testdata/output.txt"

		err := Copy(inPath, outPath, 0, 0)
		require.Nil(t, err, "error copying")

		cmp := equalfile.New(nil, equalfile.Options{}) // compare using single mode
		filesEqual, err := cmp.CompareFile(inPath, outPath)

		require.Nil(t, err, "error comparing files")
		require.True(t, filesEqual, "files are not equal")

		err = os.Remove(outPath)
		require.Nil(t, err, "error removing output file")
	})

	t.Run("limit greather than file size", func(t *testing.T) {
		inPath := "testdata/input.txt"
		inFile, err := os.Open(inPath)
		require.Nil(t, err, "error open input file")
		inFileInfo, err := inFile.Stat()
		require.Nil(t, err, "error getting input file info")
		inSize := inFileInfo.Size()
		inFile.Close()

		outPath := "testdata/output.txt"
		err = Copy(inPath, outPath, 0, inSize+10)
		require.Nil(t, err, "error copying")

		cmp := equalfile.New(nil, equalfile.Options{})
		filesEqual, err := cmp.CompareFile(inPath, outPath)

		require.Nil(t, err, "error comparing files")
		require.True(t, filesEqual, "files are not equal")

		err = os.Remove(outPath)
		require.Nil(t, err, "error removing output file")
	})

	t.Run("offset greather than file size", func(t *testing.T) {
		inPath := "testdata/input.txt"
		inFile, err := os.Open(inPath)
		require.Nil(t, err, "error open input file")
		inFileInfo, err := inFile.Stat()
		require.Nil(t, err, "error getting input file info")
		inSize := inFileInfo.Size()
		inFile.Close()

		outPath := "testdata/output.txt"
		err = Copy(inPath, outPath, inSize+10, 0)
		require.Equal(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("limit less than zero", func(t *testing.T) {
		inPath := "testdata/input.txt"
		outPath := "testdata/output.txt"

		err := Copy(inPath, outPath, 0, -1000000)
		require.Equal(t, err, ErrNegativeParameter)
	})

	t.Run("offset less than zero", func(t *testing.T) {
		inPath := "testdata/input.txt"
		outPath := "testdata/output.txt"

		err := Copy(inPath, outPath, -14, 0)
		require.Equal(t, err, ErrNegativeParameter)
	})

	t.Run("both offset and limit less than zero case", func(t *testing.T) {
		inPath := "testdata/input.txt"
		outPath := "testdata/output.txt"

		err := Copy(inPath, outPath, -14, -87)
		require.Equal(t, err, ErrNegativeParameter)
	})

	t.Run("offset equals file size", func(t *testing.T) {
		inPath := "testdata/input.txt"
		inFile, err := os.Open(inPath)
		require.Nil(t, err, "error open input file")
		inFileInfo, err := inFile.Stat()
		require.Nil(t, err, "error getting input file info")
		inSize := inFileInfo.Size()
		inFile.Close()

		outPath := "testdata/output.txt"
		err = Copy(inPath, outPath, inSize, 10)
		require.Nil(t, err)

		outFile, err := os.Open(outPath)
		require.Nil(t, err, "error open output file")
		outFileInfo, err := outFile.Stat()
		require.Nil(t, err, "error getting output file info")
		outSize := outFileInfo.Size()
		outFile.Close()

		require.Equal(t, int(outSize), 0, "output file size greather than zero")

		err = os.Remove(outPath)
		require.Nil(t, err, "error removing output file")
	})
}
