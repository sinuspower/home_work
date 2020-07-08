package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

const outFileNameSuffix string = "_validation_generated"

var (
	ErrCannotReadFile        = errors.New("can not read file")
	ErrCannotParseFile       = errors.New("can not parse file")
	ErrCannotWriteFile       = errors.New("can not write file")
	ErrCannotBuildSource     = errors.New("can not build source code")
	ErrCannotFormatOutSource = errors.New("can not format source")
)

func generate(sourcePath string) error {
	sourceBytes, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrCannotReadFile, err)
	}

	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, sourcePath, sourceBytes, 0)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrCannotParseFile, err)
	}

	templateData := getData(f, sourceBytes)

	tmpl, err := template.ParseFiles("../templates/func.tmpl", "../templates/base.tmpl")
	if err != nil {
		return err
	}

	var bb bytes.Buffer
	err = tmpl.ExecuteTemplate(&bb, "base", templateData)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrCannotBuildSource, err)
	}

	outBytes, err := format.Source(bb.Bytes())
	if err != nil {
		return fmt.Errorf("%s: %w", ErrCannotFormatOutSource, err)
	}

	var extension string
	sepIndex := strings.LastIndex(sourcePath, string(os.PathSeparator))
	dotIndex := strings.LastIndex(sourcePath, ".")
	sourceFileName := sourcePath
	if dotIndex != -1 && dotIndex > sepIndex { // some folder can have name such "foo.bar"
		sourceFileName = sourcePath[:dotIndex]
		extension = sourcePath[dotIndex:]
	}

	outFileName := sourceFileName + outFileNameSuffix + extension
	err = ioutil.WriteFile(outFileName, outBytes, 0600)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrCannotWriteFile, err)
	}

	return nil
}
