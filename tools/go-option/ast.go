// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"github.com/searKing/golang/go/reflect"
)

const (
	TagOption          = "option"
	TagOptionFlagShort = "short" // `option:",short"`
)

func isPublicName(name string) bool {
	for _, c := range name {
		return unicode.IsUpper(c)
	}
	return false
}

// FormatTypeParams turns TypeParamList into its Go representation, such as:
// [T, Y]. Note that it does not print constraints as this is mainly used for
// formatting type params in method receivers.
func FormatTypeParams(tparams *ast.FieldList) string {
	if tparams == nil || len(tparams.List) == 0 {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteByte('[')
	for i := 0; i < len(tparams.List); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		for j := 0; j < len(tparams.List[i].Names); j++ {
			if j > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(tparams.List[i].Names[j].String())
		}
	}
	buf.WriteByte(']')
	return buf.String()
}

// FormatTypeDeclaration turns TypeParamList into its Go representation, such as:
// [T, Y comparable]. Note that it does not print constraints as this is mainly used for
// formatting type params in method receivers.
func FormatTypeDeclaration(tparams *ast.FieldList) (string, error) {
	if tparams == nil || len(tparams.List) == 0 {
		return "", nil
	}

	var buf bytes.Buffer
	buf.WriteByte('[')
	for i := 0; i < len(tparams.List); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		for j := 0; j < len(tparams.List[i].Names); j++ {
			if j > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(tparams.List[i].Names[j].String())
		}
		buf.WriteString(" ")

		switch expr := tparams.List[i].Type.(type) {
		case *ast.Ident:
			buf.WriteString(expr.String())
		default:
			return "", fmt.Errorf("unsupported expression %T", expr)
		}
	}
	buf.WriteByte(']')
	return buf.String(), nil
}

// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	// Token must be in IMPORT, CONST, TYPE, VAR
	if !ok || decl.Tok != token.TYPE {
		// We only care about const|var declarations.
		return true
	}
	// The trimmedStructName of the type of the constants or variables we are declaring.
	// Can change if this is a multi-element declaration.
	typ := ""
	// Loop over the elements of the declaration. Each element is a ValueSpec:
	// a list of names possibly followed by a type, possibly followed by structs.
	// If the type and value are both missing, we carry down the type (and value,
	// but the "go/types" package takes care of that).
	for _, spec := range decl.Specs {
		tspec := spec.(*ast.TypeSpec) // Guaranteed to succeed as this is TYPE.
		typ = tspec.Name.Name
		declaration, err := FormatTypeDeclaration(tspec.TypeParams)
		if err != nil {
			// This is not the type we're looking for.
			continue
		}
		v := Struct{
			StructTypeName:               typ,
			StructTypeGenericDeclaration: declaration,
			StructTypeGenericTypeParams:  FormatTypeParams(tspec.TypeParams),
			StructTypeImport:             f.typeInfo.Import,
		}
		if c := tspec.Comment; f.lineComment && c != nil && len(c.List) == 1 {
			v.trimmedStructTypeName = strings.TrimSpace(c.Text())
		} else {
			v.trimmedStructTypeName = strings.TrimPrefix(v.StructTypeName, f.trimPrefix)
		}
		if typ != f.typeInfo.Name {
			// This is not the type we're looking for.
			continue
		}
		sExpr, ok := tspec.Type.(*ast.StructType)
		if !ok {
			f.structs = append(f.structs, v)
			continue
		}

		for _, field := range sExpr.Fields.List {
			var fieldName string
			var fieldType string
			var filedIsMap bool
			var fieldIsSlice bool
			var fieldSliceElt string
			{
				switch t := field.Type.(type) {
				case *ast.Ident:
					fieldType = t.String()
				case *ast.ArrayType:
					ident, ok := t.Elt.(*ast.Ident)
					if !ok {
						continue
					}
					if t.Len != nil {
						l, ok := t.Len.(*ast.BasicLit)
						if !ok {
							continue
						}
						fieldType = fmt.Sprintf("[%s]%s", l.Value, ident.String())
					} else {
						fieldType = fmt.Sprintf("[]%s", ident.String())
						fieldIsSlice = true
						fieldSliceElt = ident.String()
					}
				case *ast.MapType:
					k, ok := t.Key.(*ast.Ident)
					if !ok {
						continue
					}
					v, ok := t.Value.(*ast.Ident)
					if !ok {
						continue
					}
					fieldType = fmt.Sprintf("map[%s]%s", k.String(), v.String())
					filedIsMap = true
				case *ast.FuncType:
					if t.Params == nil || t.Results == nil {
						continue
					}
					if len(t.Params.List) == 0 || len(t.Results.List) == 0 {
						continue
					}
					fieldType = "func()"
				case *ast.InterfaceType:
					if t.Methods == nil {
						continue
					}
					fieldType = "interface{}"
				case *ast.SelectorExpr:
					if p, ok := t.X.(*ast.Ident); ok {
						typ := strings.Join([]string{p.String(), t.Sel.String()}, ".")
						if typ == "time.Duration" || typ == "time.Time" {
							fieldType = typ
						}
					}
					if fieldType == "" {
						continue
					}
				default:
					continue
				}
			}
			if len(field.Names) != 0 { // pick first exported Name
				for _, field := range field.Names {
					if !*flagSkipPrivateFields || isPublicName(field.Name) {
						fieldName = field.Name
						break
					}
				}
			} else { // anonymous field
				ident, ok := field.Type.(*ast.Ident)
				if !ok {
					continue
				}

				if !*flagSkipAnonymousFields {
					fieldName = ident.Name
				}
			}

			// nothing to process, continue with next line
			if fieldName == "" {
				continue
			}
			if field.Tag == nil {
				field.Tag = &ast.BasicLit{}
			}

			tags, err := reflect.ParseAstStructTag(field.Tag.Value)
			if err != nil {
				panic(err)
			}
			tagOption, _ := tags.Get(TagOption)
			if tagOption.Name == "-" {
				// ignore this field
				continue
			}

			v.Fields = append(v.Fields, StructField{
				FieldName:        fieldName,
				FieldType:        fieldType,
				FieldDocComment:  field.Doc,
				FieldLineComment: field.Comment,
				OptionTag:        tagOption,
				FieldIsSlice:     fieldIsSlice,
				FieldSliceElt:    fieldSliceElt,
				FieldIsMap:       filedIsMap,
			})
		}
		f.structs = append(f.structs, v)
	}
	return false
}
