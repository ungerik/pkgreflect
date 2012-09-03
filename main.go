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
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, dir, nil, 0)
	if err != nil {
		panic(err)
	}
	for _, pkg := range pkgs {
		if ast.PackageExports(pkg) {
			// ast.Print(fileSet, pkg)
			file, err := os.Create(path.Join(dir, "types.go"))
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(file, "// Generated file, don't modify!")
			fmt.Fprintln(file, "package", pkg.Name)
			fmt.Fprintln(file, "")
			fmt.Fprintln(file, `import "reflect"`)
			fmt.Fprintln(file, "")
			fmt.Fprintln(file, "var Types = map[string]reflect.Type{")
			ast.Inspect(pkg, func(n ast.Node) bool {
				switch s := n.(type) {
				case *ast.InterfaceType:
				case *ast.TypeSpec:
					fmt.Fprintf(file, "\t\"%s\": reflect.TypeOf((*%s)(nil)).Elem(),\n", s.Name, s.Name)
				}
				return true
			})
			fmt.Fprintln(file, "}")
		}
	}
}
