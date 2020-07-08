package main

import (
	"go/ast"
	"go/token"
	"regexp"
	"strconv"
	"strings"
)

const validationRulesRegexp string = `validate:(\s*?)"(((\s*?)(.*?)(\s*?))*?)"`

type (
	data struct {
		PackageName string
		Imports     []string
		Functions   []function
	}

	function struct {
		VarName  string
		TypeName string
		Fields   []field
	}

	field struct {
		Name     string
		Type     string
		BaseType string
		Rules    []rule
	}

	rule struct {
		Type   string
		String string
	}
)

func getData(f *ast.File, src []byte) data {
	outData := data{
		PackageName: f.Name.Name,
		Functions:   []function{},
	}
	baseTypes := getBaseTypes(f)
	var outFunctions []function
	for _, decl := range f.Decls {
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.TYPE { // type declaration
			for _, spec := range gd.Specs {
				ts, typeSpec := spec.(*ast.TypeSpec)
				if st, structType := ts.Type.(*ast.StructType); typeSpec && structType && st.Fields.List != nil { // struct declaration with fields
					outFields := getFields(st, src, baseTypes)
					if len(outFields) == 0 {
						continue
					}
					structName := ts.Name.Name
					outFunctions = append(outFunctions, function{
						VarName:  getVarName(structName),
						TypeName: structName,
						Fields:   outFields,
					})
				}
			}
		}
	}
	outData.Functions = outFunctions
	setImports(&outData)
	return outData
}

func getBaseTypes(f *ast.File) map[string]string {
	out := make(map[string]string)
	for _, decl := range f.Decls {
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.TYPE { // type declaration
			for _, spec := range gd.Specs {
				ts, typeSpec := spec.(*ast.TypeSpec)
				if i, ident := ts.Type.(*ast.Ident); typeSpec && ident {
					out[ts.Name.Name] = i.Name
				}
			}
		}
	}
	return out
}

func getFields(st *ast.StructType, src []byte, baseTypes map[string]string) []field {
	out := make([]field, 0, st.Fields.NumFields())
	for _, fld := range st.Fields.List {
		var fieldName string
		if fld.Names != nil {
			fieldName = fld.Names[0].Name
		}
		var fieldType string
		if ident, ok := fld.Type.(*ast.Ident); ok {
			fieldType = ident.Name
		}
		var fieldBaseType string
		if arrayType, ok := fld.Type.(*ast.ArrayType); ok {
			fieldType = "array"
			fieldBaseType = string(src[arrayType.Pos()+1 : arrayType.End()-1])
		}
		t, ok := baseTypes[fieldType]
		if ok {
			fieldBaseType = t
		}
		var fieldTag string
		if fld.Tag != nil {
			fieldTag = fld.Tag.Value
		}

		if fieldName == "" || fieldType == "" || fieldTag == "" {
			continue
		}

		outRules := getRules(fieldTag)
		if len(outRules) == 0 {
			continue
		}

		out = append(out, field{
			Name:     fieldName,
			Type:     fieldType,
			BaseType: fieldBaseType,
			Rules:    outRules,
		})
	}
	return out
}

func getVarName(structName string) string {
	if structName == "" {
		return ""
	}

	var sb strings.Builder
	letters := strings.Split(structName, "")
	for _, letter := range letters {
		if strings.ToLower(letter) != letter {
			sb.WriteString(strings.ToLower(letter))
		}
		if sb.Len() == 3 {
			break
		}
	}

	if sb.Len() == 0 {
		sb.WriteString(letters[0])
	}

	return sb.String()
}

func getRules(tag string) []rule {
	var out []rule

	if tag == "" {
		return out
	}

	tag = strings.TrimSpace(tag)

	re, err := regexp.Compile(validationRulesRegexp)
	if err != nil {
		return out
	}
	match := strings.TrimSpace(re.FindString(tag))

	if match == "" {
		return out
	}

	// sequential trimming to avoid unexpected deletion of necessary characters
	match = strings.Trim(strings.Trim(match, "validte:\r\n\t"), `"`)
	ruleStrs := strings.Split(match, "|")
	for _, str := range ruleStrs {
		str = strings.TrimSpace(str)
		colonPos := strings.Index(str, ":")
		if colonPos != -1 && colonPos < len(str)-1 {
			rType := strings.TrimSpace(str[:colonPos])
			rString := strings.TrimSpace(str[colonPos+1:])
			if rType != "" && rString != "" {
				if !checkIntRule(rType, rString) {
					continue
				}
				out = append(out, rule{
					Type:   rType,
					String: rString,
				})
			}
		}
	}

	return out
}

func checkIntRule(rType, rString string) bool {
	if !(rType == "min" || rType == "max" || rType == "len") {
		return true
	}
	_, err := strconv.Atoi(rString)
	return err == nil
}

func setImports(d *data) {
	imports := []string{"errors"}
	for _, function := range d.Functions {
		for _, field := range function.Fields {
			if field.Type == "string" || field.Type == "[]string" {
				imports = append(imports, getImportsString(field)...)
			}
			if field.Type == "int" || field.Type == "[]int" {
				imports = append(imports, getImportsInt(field)...)
			}
		}
	}
	d.Imports = imports
}

func getImportsString(f field) []string {
	var imports []string
	for _, rule := range f.Rules {
		switch rule.Type {
		case "regexp":
			if !in(imports, "regexp") {
				imports = append(imports, "regexp")
			}
		case "in":
			if !in(imports, "strings") {
				imports = append(imports, "strings")
			}
		}
	}
	return imports
}

func getImportsInt(f field) []string {
	var imports []string
	for _, rule := range f.Rules {
		if rule.Type == "in" {
			if !in(imports, "strings") {
				imports = append(imports, "strings")
			}
			if !in(imports, "strconv") {
				imports = append(imports, "strconv")
			}
		}
	}
	return imports
}

func in(array []string, str string) bool {
	var found bool
	for _, s := range array {
		found = str == s
		if found {
			break
		}
	}
	return found
}
