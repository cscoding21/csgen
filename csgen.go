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

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

// GetFile returns the current file based on the generators env or passed in value
func GetFile(file ...string) string {
	//---if the user passes in the file path, return an absolute version
	if len(file) > 0 && file[0] != "" {
		fp, err := filepath.Abs(file[0])
		if err != nil {
			panic(err)
		}

		return fp
	}

	//---if nothing is passed in, use the generators file name from the env
	pwd, _ := os.Getwd()
	f := os.Getenv("GOFILE")

	return filepath.Join(pwd, f)
}

// GetStructs return a list of all structs in a given file
func GetStructs(filePath string) ([]Struct, error) {
	// https://magodo.github.io/go-ast-tips/
	out := []Struct{}

	f, err := GetAST(filePath)
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
						Fields:   []Field{},
					}

					switch typeSpec.Type.(type) {
					case *ast.StructType:
						structType := typeSpec.Type.(*ast.StructType)
						for _, field := range structType.Fields.List {
							fieldType := types.ExprString(field.Type)
							tagString := ""

							fmt.Println(fieldType)
							fmt.Printf("Struct: name=%s\n", typeSpec.Name.Name)

							if field.Tag != nil {
								tagString = field.Tag.Value
							}

							for _, name := range field.Names {
								s := Field{
									Name:        name.Name,
									Type:        fieldType,
									TagString:   tagString,
									IsPrimitive: IsPrimitive(fieldType),
									IsPointer:   IsRefType(fieldType),
									IsSlice:     IsSlice(fieldType),
									IsPublic:    IsPublic(name.Name),
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

// GetVariables returns a list of all variable definitions in a given file
func GetVariables(filePath string) ([]Field, error) {
	out := []Field{}

	f, err := GetAST(filePath)
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
						thisVar := Field{
							Name:        id.Name,
							Type:        types.ExprString(spec.Type),
							TagString:   "",
							IsPrimitive: IsPrimitive(types.ExprString(spec.Type)),
							IsPointer:   IsRefType(types.ExprString(spec.Type)),
							IsSlice:     IsSlice(types.ExprString(spec.Type)),
							IsPublic:    IsPublic(id.Name),
						}

						out = append(out, thisVar)
					}
				default:
					//fmt.Printf("Unknown token type: %s\n", decl.Tok)
				}
			}
		default:
			//fmt.Printf("Unknown declaration %v\n", decl.Pos())
		}
	}

	return out, nil
}

// GetImports returns all of the imports in a given file
func GetImports(filePath string) ([]string, error) {
	out := []string{}

	f, err := GetAST(filePath)
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
					//fmt.Printf("Unknown token type: %s\n", decl.Tok)
				}
			}
		default:
			//fmt.Printf("Unknown declaration %v\n", decl.Pos())
		}
	}

	return out, nil
}

// GetFunctions returns all of the functions in a given file
func GetFunctions(filePath string) ([]Function, error) {
	out := []Function{}

	f, err := GetAST(filePath)
	if err != nil {
		fmt.Println(err)
		return out, err
	}

	for _, d := range f.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:
			fn := Function{
				Name:      decl.Name.Name,
				Arguments: []Field{},
				Returns:   []Field{},
				IsPublic:  IsPublic(decl.Name.Name),
			}

			out = append(out, fn)
		default:
			//fmt.Printf("Unknown declaration %v\n", decl.Pos())
		}
	}

	return out, nil
}

// GetInterfaces get a list of all declared interfaces in a given file
func GetInterfaces(filePath string) ([]Interface, error) {
	out := []Interface{}

	f, err := GetAST(filePath)
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
						in := Interface{
							Name:     spec.Name.Name,
							Methods:  []Function{},
							IsPublic: IsPublic(spec.Name.Name),
						}

						ast.Inspect(f, func(n ast.Node) bool {
							if fd, ok := n.(*ast.FuncDecl); ok {
								fmt.Printf("Function: %s, parameters:\n", fd.Name)
								for _, param := range fd.Type.Params.List {
									fmt.Printf("  Name: %s\n", param.Names[0])
									fmt.Printf("    ast type          : %T\n", param.Type)
									fmt.Printf("    type desc         : %+v\n", param.Type)
								}
							}
							return true
						})

						//---FOR REFERENCE: http://goast.yuroyoro.net/
						//					https://gist.github.com/ncdc/fef1099f54a655f8fb11daf86f7868b8
						// for _, field := range spec.Type.Methods.List{
						// 	fn := Function{
						// 		Name:      field.Names[0].Name,
						// 		Arguments: []StructField{},
						// 		Returns:   []StructField{},
						// 		IsPublic:  isPublic(field.Names[0].Name),
						// 	}

						// 	in.Methods = append(in.Methods, fn)
						// }

						out = append(out, in)
					}
				default:
					//fmt.Printf("Unknown token type: %s\n", decl.Tok)
				}
			}
		default:
			//fmt.Printf("Unknown declaration %v\n", decl.Pos())
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

	code, err = imports.Process(name, code, nil)
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

	builder.WriteString("// Code generated . DO NOT EDIT.\n")
	builder.WriteString("// ################################## DO NOT EDIT THIS FILE ######################################\n")
	builder.WriteString("// Common Sense Coding (https://github.com/cscoding21/csgen)\n")
	builder.WriteByte('\n')
	builder.WriteString(fmt.Sprintf("// Generate Date: %v\n", time.Now()))
	builder.WriteString(fmt.Sprintf("// Implementation Name: %s\n", name))
	builder.WriteString("// Developer Note: The contents of this file will be recreated each time its generator is called\n")
	builder.WriteByte('\n')
	builder.WriteString("// -----------------------------------------------------------------------------------------------\n")
	builder.WriteByte('\n')
	builder.WriteByte('\n')
	builder.WriteString(fmt.Sprintf("package %s", pkg))
	builder.WriteByte('\n')

	return &builder
}

// NewCSGenBuilderForOneOffFile returns a string buider with a common header for generated files that are indented to be modified
func NewCSGenBuilderForOneOffFile(name string, pkg string) *strings.Builder {
	builder := strings.Builder{}

	builder.WriteString("// --------------------------------- GENERATED FILE : OK TO EDIT ---------------------------------\n")
	builder.WriteString("// Common Sense Coding (https://github.com/cscoding21/csgen)\n")
	builder.WriteByte('\n')
	builder.WriteString(fmt.Sprintf("// Generate Date: %v\n", time.Now()))
	builder.WriteString(fmt.Sprintf("// Implementation Name: %s\n", name))
	builder.WriteString("// Developer Note: This file will only be generated once.  It is intended to be modified.n")
	builder.WriteByte('\n')
	builder.WriteString("// -----------------------------------------------------------------------------------------------\n")
	builder.WriteByte('\n')
	builder.WriteByte('\n')
	builder.WriteString(fmt.Sprintf("package %s", pkg))
	builder.WriteByte('\n')

	return &builder
}

// ProfileNode get details about an unknown node based on its actual type
func ProfileNode(node ast.Node) {
	if node == nil {
		return
	}

	switch x := node.(type) {
	case *ast.CallExpr:
		id, ok := x.Fun.(*ast.Ident)
		if ok {
			fmt.Printf("CallExpr: %v", id.Name)
		}
	}
}

// GetDefaultPackageConfig return a default set of values for loading module packages
func GetDefaultPackageConfig() *packages.Config {
	cfg := &packages.Config{
		Tests: false,
		Mode: packages.NeedSyntax |
			packages.NeedName |
			packages.NeedEmbedFiles |
			packages.NeedFiles |
			packages.NeedTypes |
			packages.NeedModule |
			packages.NeedDeps,
	}

	return cfg
}

// LoadModule use the "packages" package to load codebase items
func LoadModule(cfg *packages.Config) (Module, error) {
	module := Module{}
	cwd, err := os.Getwd()
	if err != nil {
		return module, err
	}

	module.Path = cwd

	pkgs, err := packages.Load(cfg, cwd)
	if err != nil {
		return module, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		//---TODO: figure out if this is an error condition
		return module, fmt.Errorf("no packages loaded")
	}

	for _, pkg := range pkgs {
		outPackage := Package{
			Name:    pkg.Name,
			Path:    pkg.PkgPath,
			Files:   pkg.GoFiles,
			Structs: []Struct{},
		}

		if len(module.Name) == 0 {
			module.Name = pkg.Module.Path
		}

		// qual := types.RelativeTo(pkg.Types)
		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			if obj != nil && obj.Type() != nil && obj.Type().Underlying() != nil {
				// obj is types.Named,
				// obj.Type() is types.TypeName
				// obj.Type().Underlying() exposes the types.Struct which
				// gives access to the desired types.Var.Embedded()
				st, ok := obj.Type().Underlying().(*types.Struct)
				if !ok {
					continue
				}
				fmt.Printf("Struct: %s - %s\n", obj.Name(), st.String())

				outStruct := Struct{
					Name:    obj.Name(),
					Type:    obj.Type().String(),
					Package: pkg.Name,
					Fields:  []Field{},
				}

				for i := 0; i < st.NumFields(); i++ {
					f := st.Field(i)
					if f.Embedded() /* && strings.HasSuffix(f.Type().String(), "log.Logger") */ {
						ist, ok := f.Origin().Type().Underlying().(*types.Struct)
						if !ok {
							continue
						}

						for x := 0; x < ist.NumFields(); x++ {
							thisField := ist.Field(x)
							ft := thisField.Type().String()

							outField := Field{
								Name:        thisField.Name(),
								Type:        ft,
								TagString:   "",
								IsPrimitive: IsPrimitive(ft),
								IsPointer:   IsRefType(ft),
								IsSlice:     IsSlice(ft),
								IsPublic:    IsPublic(thisField.Name()),
							}

							outStruct.EmbeddedFields = append(outStruct.EmbeddedFields, outField)
						}
					} else {
						ft := f.Type().String()
						outField := Field{
							Name:        f.Name(),
							Type:        ft,
							TagString:   "",
							IsPrimitive: IsPrimitive(ft),
							IsPointer:   IsRefType(ft),
							IsSlice:     IsSlice(ft),
							IsPublic:    IsPublic(f.Name()),
						}

						outStruct.Fields = append(outStruct.Fields, outField)
					}
				}

				outPackage.Structs = append(outPackage.Structs, outStruct)
			}
		}

		module.Packages = append(module.Packages, outPackage)
	}

	return module, nil
}
