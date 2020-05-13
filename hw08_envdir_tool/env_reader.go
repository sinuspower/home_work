package main

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]string

var (
	ErrCanNotOpenDir  = errors.New("can not open directory")
	ErrWrongFileName  = errors.New("wrong file name")
	ErrCanNotOpenFile = errors.New("can not open file")
	ErrCanNotReadLine = errors.New("can not read line from file")
)

func readString(f *os.File) (string, error) {
	fInfo, err := f.Stat()
	if err != nil {
		return "", err
	}
	if fInfo.Size() == 0 {
		return "", nil
	}
	reader := bufio.NewReader(f)
	result := ""
	isPrefix := true
	for isPrefix {
		line, p, err := reader.ReadLine()
		if err != nil {
			return "", err
		}
		result += string(bytes.ReplaceAll(line, []byte{0}, []byte("\n")))
		isPrefix = p
	}
	return result, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, ErrCanNotOpenDir
	}

	env := make(Environment)
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if strings.Contains(name, "=") {
				return nil, ErrWrongFileName
			}
			f, err := os.Open(dir + string(os.PathSeparator) + name)
			if err != nil {
				return nil, ErrCanNotOpenFile
			}
			value, err := readString(f)
			if err != nil {
				return nil, ErrCanNotReadLine
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
