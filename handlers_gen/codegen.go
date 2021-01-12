package main

import (
	"go/ast"
	"go/token"
	"go/parser"
	// "text/template"
	"os"
	"log"
	"fmt"
	"strings"
)

// код писать тут

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := os.Create(os.Args[2])

	fmt.Fprintln(out, `package `+node.Name.Name)
	fmt.Fprintln(out) // empty line
	fmt.Fprintln(out, `import "encoding/binary"`)
	fmt.Fprintln(out, `import "bytes"`)
	fmt.Fprintln(out) // empty line

	for _, f := range node.Decls {
		// g, ok := f.(*ast.GenDecl)
		// if ok {
		// 	// SPECS_LOOP:
		// 	for _, spec := range g.Specs {
		// 		// получаем тип спецификации
		// 		currType, ok := spec.(*ast.TypeSpec)
		// 		if !ok {
		// 			fmt.Printf("SKIP %#T is not ast.TypeSpec\n", spec)
		// 			continue
		// 		}

		// 		fmt.Println("currType", currType.Name.Name)

		// 		currStruct, ok := currType.Type.(*ast.StructType)
		// 		if !ok {
		// 			fmt.Printf("SKIP %#T is not ast.StructType\n", currStruct)
		// 			continue
		// 		}
		// 	}
		// 	continue
		// }

		f, ok := f.(*ast.FuncDecl)
		if ok {
			if f.Doc == nil {
				fmt.Println("SKIP comment empty: ", f.Name.Name)
				continue
			}

			needCodegen := false
			for _, comment := range f.Doc.List {
				fmt.Println("comment: ", comment)
				needCodegen = needCodegen || strings.HasPrefix(comment.Text, "// apigen:api")
			}
		}
	}
}