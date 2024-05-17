package csgen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"go/ast"
	"go/format"
	"go/types"
)

// GetFile returns the current file based on the generators env or passed in value
func GetFile(file ...string) string {
	pwd, _ := os.Getwd()
	f := os.Getenv("GOFILE")

	if len(file) > 0 && file[0] != "" {
		f = "/" + file[0]
	}

	return filepath.Join(pwd, f)
}

// GetStructs return a list of all structs in a given file
func GetStructs(filePath string) ([]Struct, error) {
	// https://magodo.github.io/go-ast-tips/
	out := []Struct{}

	f, err := getAst(filePath)
	if err != nil {
		fmt.Println(err)
		return out, err
	}

	pkgName := f.Name.Name

	for _, node := range f.Decls {
		switch node.(type) {

		case *ast.GenDecl:
			genDecl := node.(*ast.GenDecl)
			for _, spec := range genDecl.Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					typeSpec := spec.(*ast.TypeSpec)

					outStruct := Struct{
						FilePath: filePath,
						Package:  pkgName,
						Name:     typeSpec.Name.Name,
						Fields:   []StructField{},
					}

					switch typeSpec.Type.(type) {
					case *ast.StructType:
						structType := typeSpec.Type.(*ast.StructType)
						for _, field := range structType.Fields.List {
							fieldType := types.ExprString(field.Type)
							tagString := ""

							if field.Tag != nil {
								tagString = field.Tag.Value
							}

							for _, name := range field.Names {
								s := StructField{
									Name:        name.Name,
									Type:        fieldType,
									TagString:   tagString,
									IsPrimitive: isPrimitive(fieldType),
									IsPointer:   isRefType(fieldType),
									IsSlice:     isSlice(fieldType),
								}

								outStruct.Fields = append(outStruct.Fields, s)
							}
						}

						out = append(out, outStruct)
					}
				}
			}
		}
	}

	return out, nil
}

// GetVariables
func GetVariables(filePath string) ([]StructField, error) {
	out := []StructField{}

	f, err := getAst(filePath)
	if err != nil {
		fmt.Println(err)
		return out, err
	}

	for _, d := range f.Decls {
		switch decl := d.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.ValueSpec:
					for _, id := range spec.Names {
						thisVar := StructField{
							Name:        id.Name,
							Type:        types.ExprString(spec.Type),
							TagString:   "",
							IsPrimitive: isPrimitive(types.ExprString(spec.Type)),
							IsPointer:   isRefType(types.ExprString(spec.Type)),
							IsSlice:     isSlice(types.ExprString(spec.Type)),
						}

						out = append(out, thisVar)
					}
				default:
					fmt.Printf("Unknown token type: %s\n", decl.Tok)
				}
			}
		default:
			fmt.Printf("Unknown declaration %v\n", decl.Pos())
		}
	}

	return out, nil
}

// GetImports returns all of the imports in a given file
func GetImports(filePath string) ([]string, error) {
	out := []string{}

	f, err := getAst(filePath)
	if err != nil {
		fmt.Println(err)
		return out, err
	}

	for _, d := range f.Decls {
		switch decl := d.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.ImportSpec:
					out = append(out, spec.Path.Value)
				default:
					fmt.Printf("Unknown token type: %s\n", decl.Tok)
				}
			}
		default:
			fmt.Printf("Unknown declaration %v\n", decl.Pos())
		}
	}

	return out, nil
}

// GetFunctions returns all of the functions in a given file
func GetFunctions(filePath string) ([]string, error) {
	out := []string{}

	f, err := getAst(filePath)
	if err != nil {
		fmt.Println(err)
		return out, err
	}

	for _, d := range f.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:
			fmt.Printf("Func: %s\n", decl.Name.Name)
			out = append(out, decl.Name.Name)
		default:
			fmt.Printf("Unknown declaration %v\n", decl.Pos())
		}
	}

	return out, nil
}

// GetInterfaces get a list of all declared interfaces in a given file
func GetInterfaces(filePath string) ([]string, error) {
	out := []string{}

	f, err := getAst(filePath)
	if err != nil {
		fmt.Println(err)
		return out, err
	}

	for _, d := range f.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:
			fmt.Println("Func")
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					switch spec.Type.(type) {
					case *ast.InterfaceType:
						out = append(out, spec.Name.String())
					}
				default:
					fmt.Printf("Unknown token type: %s\n", decl.Tok)
				}
			}
		default:
			fmt.Printf("Unknown declaration %v\n", decl.Pos())
		}
	}

	return out, nil
}

// ExecuteTemplate executes a template against a given object and return the output as a string
func ExecuteTemplate[T any](name string, fileTemplate string, om T) string {
	tmpl, err := template.New(name).Parse(fileTemplate)
	if err != nil {
		panic(err)
	}

	var doc bytes.Buffer
	tmpl.Execute(&doc, om)
	if err != nil {
		panic(err)
	}

	return doc.String()
}

// WriteGeneratedGoFile create a text file with the passed in name and contents
func WriteGeneratedGoFile(name string, contents string) error {
	code, err := format.Source([]byte(contents))
	if err != nil {
		return err
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}

	l, err := f.WriteString(string(code))
	if err != nil {
		f.Close()
		return err
	}

	fmt.Printf("%v bytes written successfully into file %s", l, name)
	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

// NewCSGenBuilderForFile returns a string buider with a common header for generated files
func NewCSGenBuilderForFile(name string, pkg string) *strings.Builder {
	builder := strings.Builder{}

	builder.WriteString("// ################################## DO NOT EDIT THIS FILE ######################################\n")
	builder.WriteString("// ### Common Sense Coding\n")
	builder.WriteByte('\n')
	builder.WriteString("// This file contains generated code. DO NOT EDIT.\n")
	builder.WriteString(fmt.Sprintf("// Generate Date: %v\n", time.Now()))
	builder.WriteString(fmt.Sprintf("// Implementation Name: %s\n", name))
	builder.WriteByte('\n')
	builder.WriteString("// -----------------------------------------------------------------------------------------------\n")
	builder.WriteByte('\n')
	builder.WriteByte('\n')
	builder.WriteString(fmt.Sprintf("package %s", pkg))
	builder.WriteByte('\n')

	return &builder
}
