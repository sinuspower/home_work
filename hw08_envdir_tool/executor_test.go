package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("empty cmd", func(t *testing.T) {
		cmd := []string{}
		returnCode := RunCmd(cmd, nil)
		require.Equal(t, 0, returnCode)
	})

	t.Run("cmd is nil", func(t *testing.T) {
		returnCode := RunCmd(nil, nil)
		require.Equal(t, 0, returnCode)
	})

	t.Run("return codes", func(t *testing.T) {
		f, err := os.Create("testdata/test")
		require.Nil(t, err)

		_, err = f.Write([]byte("File contents."))
		require.Nil(t, err)
		f.Close()

		cmd := []string{"cat", "testdata/test"}
		returnCode := RunCmd(cmd, nil)
		require.Equal(t, 0, returnCode)

		cmd = []string{"cat", "testdata/notExist"}
		returnCode = RunCmd(cmd, nil)
		require.Equal(t, 1, returnCode)

		err = os.Remove("testdata/test")
		require.Nil(t, err)
	})

	t.Run("command output", func(t *testing.T) {
		f, err := os.Create("testdata/test")
		require.Nil(t, err)

		_, err = f.Write([]byte(`This is a test file contents.
		This contents must be printed into Stdout by "cat" command
		executed from the "RunCmd" function.`))
		require.Nil(t, err)
		f.Close()

		cmd := []string{"cat", "testdata/test"}
		result, err := catchStdout(RunCmd, cmd, nil)
		require.Nil(t, err)

		expected := `This is a test file contents.
		This contents must be printed into Stdout by "cat" command
		executed from the "RunCmd" function.`

		require.Equal(t, string(result), expected)

		err = os.Remove("testdata/test")
		require.Nil(t, err)
	})
}

func catchStdout(runCmd func(cmd []string, env Environment) int, cmd []string, env Environment) (result []byte, err error) {
	realOut := os.Stdout
	defer func() { os.Stdout = realOut }()

	r, fakeOut, err := os.Pipe()
	if err != nil {
		return
	}

	os.Stdout = fakeOut
	_ = runCmd(cmd, env)
	if err = fakeOut.Close(); err != nil {
		return
	}

	result, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}

	err = r.Close()
	return result, err
}
