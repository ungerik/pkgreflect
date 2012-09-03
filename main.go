package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
)

func main() {
	dir := "../go-start/view"
	os.Remove("../go-start/view/_pkg_reflect.go")
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, dir, nil, 0)
	if err != nil {
		panic(err)
	}
	for _, pkg := range pkgs {
		if ast.PackageExports(pkg) {

			// fmt.Printf("%v", pkg.Files["../go-start/view/jquery.go"])

			// ast.Print(fileSet, pkg)
			file, err := os.Create(path.Join(dir, "_pkg_reflect.go"))
			if err != nil {
				panic(err)
			}
			defer file.Close()
			fmt.Fprintln(file, "// Generated file, don't modify!")
			fmt.Fprintln(file, "package", pkg.Name)
			fmt.Fprintln(file, "")
			fmt.Fprintln(file, `import "reflect"`)
			fmt.Fprintln(file, "")

			// Types
			fmt.Fprintln(file, "var Types = map[string]reflect.Type{")
			ast.Inspect(pkg, func(n ast.Node) bool {
				switch s := n.(type) {
				case *ast.TypeSpec:
					fmt.Fprintf(file, "\t\"%s\": reflect.TypeOf((*%s)(nil)).Elem(),\n", s.Name, s.Name)
				}
				return true
			})
			fmt.Fprintln(file, "}")
			fmt.Fprintln(file, "")

			// Functions
			fmt.Fprintln(file, "var Functions = map[string]reflect.Value{")
			ast.Inspect(pkg, func(n ast.Node) bool {
				switch s := n.(type) {
				case *ast.FuncDecl:
					if s.Recv == nil && s.Body != nil {
						fmt.Fprintf(file, "\t\"%s\": reflect.ValueOf(%s),\n", s.Name, s.Name)
					}
				}
				return true
			})
			fmt.Fprintln(file, "}")
			fmt.Fprintln(file, "")

			// Addresses of variables
			fmt.Fprintln(file, "var Variables = map[string]reflect.Value{")
			for _, f := range pkg.Files {
				for name, object := range f.Scope.Objects {
					if object.Kind == ast.Var && ast.IsExported(name) {
						fmt.Fprintf(file, "\t\"%s\": reflect.ValueOf(&%s),\n", name, name)
					}
				}
			}
			fmt.Fprintln(file, "}")
			fmt.Fprintln(file, "")
		}
	}
}
