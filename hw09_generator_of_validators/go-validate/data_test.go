package main

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetVarName(t *testing.T) {
	testCases := []struct {
		in       string
		expected string
	}{
		{"A", "a"},
		{"lower", "l"},
		{"Variable", "v"},
		{"LongVariableName", "lvn"},
		{"VeryLongVariableName", "vlv"},
		{"LONGUPPER", "lon"},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			require.Equal(t, tc.expected, getVarName(tc.in))
		})
	}
}

func TestGetRules(t *testing.T) {
	testCases := []struct {
		in       string
		expected []rule
	}{
		{
			`json:id`,
			nil,
		},
		{
			``,
			nil,
		},
		{
			`validate`,
			nil,
		},
		{
			`validate:"len:"`,
			nil,
		},
		{
			`validate:"in:fg, rhj|len:bad|min:bad|max:bad"`,
			[]rule{
				rule{
					Type:   "in",
					String: "fg, rhj",
				},
			},
		},
		{
			`json:"id" validate:"len:36"`,
			[]rule{
				rule{
					Type:   "len",
					String: "36",
				},
			},
		},
		{
			`validate:"min:18|max:50|in:21,30,45|len:20|regexp:^\\w+@\\w+\\.\\w+$"`,
			[]rule{
				rule{
					Type:   "min",
					String: "18",
				},
				rule{
					Type:   "max",
					String: "50",
				},
				rule{
					Type:   "in",
					String: "21,30,45",
				},
				rule{
					Type:   "len",
					String: "20",
				},
				rule{
					Type:   "regexp",
					String: `^\\w+@\\w+\\.\\w+$`,
				},
			},
		},
		{
			`	validate:
			"min:  18 |      max:50|
			in:21,30,45|len:20   |
			regexp:^\\w+@\\w+\\.\\w+$"`,
			[]rule{
				rule{
					Type:   "min",
					String: "18",
				},
				rule{
					Type:   "max",
					String: "50",
				},
				rule{
					Type:   "in",
					String: "21,30,45",
				},
				rule{
					Type:   "len",
					String: "20",
				},
				rule{
					Type:   "regexp",
					String: `^\\w+@\\w+\\.\\w+$`,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			_ = getRules(tc.in)
			require.Equal(t, tc.expected, getRules(tc.in))
		})
	}
}

func TestIn(t *testing.T) {
	testCases := []struct {
		name     string
		inArray  []string
		inStr    string
		expected bool
	}{
		{"empty array", []string{}, "string", false},
		{"empty array, emply string", []string{}, "", false},
		{"string found", []string{"one", "two", "three"}, "two", true},
		{"string not found", []string{"one", "two", "three"}, "four", false},
		{"empty string", []string{"one", "two", "three"}, "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, in(tc.inArray, tc.inStr))
		})
	}
}

func TestGetBaseTypes(t *testing.T) {
	t.Run("UserRole", func(t *testing.T) {
		sourcePath := "../models/models.go"
		sourceBytes, err := ioutil.ReadFile(sourcePath)
		require.NoError(t, err, "can not open file")

		fileSet := token.NewFileSet()
		f, err := parser.ParseFile(fileSet, sourcePath, sourceBytes, 0)
		require.NoErrorf(t, err, "can not parse file: %w", err)

		expected := map[string]string{"UserRole": "string"}
		require.Equal(t, expected, getBaseTypes(f))
	})
}

func TestGetData(t *testing.T) {
	var expected data = data{
		PackageName: "models",
		Imports:     []string{"errors", "regexp", "strings", "strconv"},
		Functions: []function{
			{
				VarName:  "u",
				TypeName: "User",
				Fields: []field{
					{
						Name:     "ID",
						Type:     "string",
						BaseType: "",
						Rules: []rule{
							{
								Type:   "len",
								String: "36",
							},
						},
					},
					{
						Name:     "Age",
						Type:     "int",
						BaseType: "",
						Rules: []rule{
							{
								Type:   "min",
								String: "18",
							},
							{
								Type:   "max",
								String: "50",
							},
						},
					},
					{
						Name:     "Email",
						Type:     "string",
						BaseType: "",
						Rules: []rule{
							{
								Type:   "regexp",
								String: `^\\w+@\\w+\\.\\w+$`,
							},
						},
					},
					{
						Name:     "Role",
						Type:     "UserRole",
						BaseType: "string",
						Rules: []rule{
							{
								Type:   "in",
								String: "admin,stuff",
							},
						},
					},
					{
						Name:     "Phones",
						Type:     "array",
						BaseType: "string",
						Rules: []rule{
							{
								Type:   "len",
								String: "11",
							},
						},
					},
				},
			},
			{
				VarName:  "a",
				TypeName: "App",
				Fields: []field{
					{
						Name:     "Version",
						Type:     "string",
						BaseType: "",
						Rules: []rule{
							{
								Type:   "len",
								String: "5",
							},
						},
					},
				},
			},
			{
				VarName:  "r",
				TypeName: "Response",
				Fields: []field{
					{
						Name:     "Code",
						Type:     "int",
						BaseType: "",
						Rules: []rule{
							{
								Type:   "in",
								String: "200,404,500",
							},
						},
					},
					{
						Name:     "Desc",
						Type:     "array",
						BaseType: "string",
						Rules: []rule{
							{
								Type:   "in",
								String: "ok, all good, not found, we miss it, server error, bad server",
							},
						},
					},
				},
			},
		},
	}

	t.Run("build data tree", func(t *testing.T) {
		sourcePath := "../models/models.go"
		sourceBytes, err := ioutil.ReadFile(sourcePath)
		require.NoError(t, err, "can not open file")

		fileSet := token.NewFileSet()
		f, err := parser.ParseFile(fileSet, sourcePath, sourceBytes, 0)
		require.NoErrorf(t, err, "can not parse file: %w", err)

		require.Equal(t, expected, getData(f, sourceBytes))
	})
}
