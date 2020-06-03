package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]string

var (
	ErrCanNotOpenDir  = errors.New("can not open directory")
	ErrWrongFileName  = errors.New("wrong file name")
	ErrCanNotOpenFile = errors.New("can not open file")
	ErrCanNotReadLine = errors.New("can not read line from file")
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCanNotOpenDir, err)
	}

	env := make(Environment)
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if strings.Contains(name, "=") {
				return nil, ErrWrongFileName
			}
			f, err := os.Open(filepath.Join(dir, name))
			if err != nil {
				return nil, fmt.Errorf("%s: %w", ErrCanNotOpenFile, err)
			}
			value, err := readString(f)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", ErrCanNotReadLine, err)
			}
			env[name] = value
			f.Close()
		}
	}

	if len(env) == 0 {
		env = nil
	}

	return env, nil
}

func readString(f *os.File) (string, error) {
	fInfo, err := f.Stat()
	if err != nil {
		return "", err
	}
	if fInfo.Size() == 0 {
		return "", nil
	}
	reader := bufio.NewReader(f)

	line, err := reader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return "", nil
	}

	result := string(bytes.ReplaceAll(line, []byte{0}, []byte("\n")))
	return strings.TrimRight(result, "\n\t "), nil
}
