package main

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"os"
)

func init() {
	http.HandleFunc("/declarations", func(rw http.ResponseWriter, req *http.Request) {
		m := map[string]interface{}{}
		decls := []map[string]string{}
		fset := token.NewFileSet()
		var err error
		fn := req.FormValue("path")
		var src interface{} = req.FormValue("src")
		if src.(string) == "" {
			src, err = os.Open(fn)
		}
		if err == nil {
			if fn == "" {
				fn = "<stdin>"
			}
			var af *ast.File
			af, err = parser.ParseFile(fset, fn, src, 0)
			if af != nil {
				for _, d := range af.Decls {
					if p := fset.Position(d.Pos()); p.IsValid() {
						switch n := d.(type) {
						case *ast.FuncDecl:
							decls = append(decls, map[string]string{
								"name":     n.Name.Name,
								"kind":     "func",
								"location": p.String(),
								"doc":      n.Doc.Text(),
							})
						case *ast.GenDecl:
							for _, spec := range n.Specs {
								switch gn := spec.(type) {
								case *ast.TypeSpec:
									if vp := fset.Position(gn.Pos()); vp.IsValid() {
										decls = append(decls, map[string]string{
											"name":     gn.Name.Name,
											"kind":     "type",
											"location": vp.String(),
											"doc":      gn.Doc.Text(),
										})
									}
								case *ast.ValueSpec:
									for _, v := range gn.Names {
										if vp := fset.Position(v.Pos()); vp.IsValid() {
											kind := ""
											switch v.Obj.Kind {
											case ast.Typ:
												kind = "type"
											case ast.Fun:
												kind = "func"
											case ast.Con:
												kind = "constant"
											case ast.Var:
												kind = "variable"
											default:
												continue
											}
											decls = append(decls, map[string]string{
												"name":     v.Name,
												"kind":     kind,
												"location": vp.String(),
												"doc":      "",
											})
										}
									}
								}
							}
						}
					}
				}
			}
		}

		m["declarations"] = decls
		if err != nil {
			m["error"] = "Error: `" + err.Error()
		}
		json.NewEncoder(rw).Encode(m)
	})
}
