package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"text/template"
	"encoding/json"
)

// код писать тут

type handlerConfig struct {
	URL string `json:"url"`
	Method string
	Auth bool
}
type wrapperTplParams struct {
	ApiName     string
	MethodBlock string
	ParamsName string
	HandlerName string
}

var wrapperTpl = template.Must(template.New("wrapperTpl").Parse(`
	func (h *{{.ApiName}}) wrapper{{.HandlerName}}(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		{{.MethodBlock}}

		var params {{.ParamsName}}
		if p, err := parse{{.ParamsName}}(r); err == nil {
			params = p
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := h.{{.HandlerName}}(ctx, params)
		if err == nil {
			// send success response
			w.WriteHeader(http.StatusOK)
			respondJSON(w, res)
			return
		}

		switch err.(type) {
		case ApiError:
			e := err.(ApiError)
			w.WriteHeader(e.HTTPStatus)
			w.Write([]byte(e.Error()))

		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
	`))

var methodBlockTpl = template.Must(template.New("methodBlockTpl").Parse(`
if r.Method != {{.MethodStr}} {
	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}
`))

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := os.Create(os.Args[2])

	fmt.Fprintln(out, `package `+node.Name.Name)
	fmt.Fprintln(out) // empty line
	fmt.Fprintln(out, `import "net/http"`)
	fmt.Fprintln(out, `import "encoding/json"`)
	fmt.Fprintln(out, `import "context"`)
	fmt.Fprintln(out) // empty line

	fmt.Fprintln(out, `
	func respondJSON(w http.ResponseWriter, data interface{}) {
		d, _ := json.Marshal(data)
		w.Write(d)
	}
	`)

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
			f.Recv.List

			if f.Doc == nil {
				fmt.Println("SKIP comment empty: ", f.Name.Name)
				continue
			}

			codegenRaw := ""
			for _, comment := range f.Doc.List {
				fmt.Println("comment: ", comment)
				pref := "// apigen:api "
				if strings.HasPrefix(comment.Text, pref) {
					codegenRaw = strings.Replace(comment.Text, pref, "", -1)
				}
			}

			if codegenRaw == "" {
				fmt.Println("SKIP codegenRaw is empty")
				continue
			}

			var hConf handlerConfig
			_ = json.Unmarshal([]byte(codegenRaw), &hConf)

			err := wrapperTpl.Execute(out, wrapperTplParams{ApiName: "MyApi", HandlerName: "Create", ParamsName: "CreateParams"})
			if err != nil {
				panic(err.Error())
			}
		}
	}
}
