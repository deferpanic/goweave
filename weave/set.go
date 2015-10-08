package weave

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// applySetJP applies any advice for set joinpoints
func (w *Weave) applySetJP(fname string, stuff string) string {

	rout := stuff

	importsNeeded := []string{}

	for i := 0; i < len(w.aspects); i++ {

		aspect := w.aspects[i]
		if !aspect.pointkut.isSet() {
			continue
		}

		pk := aspect.pointkut.def

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, fname, rout, parser.Mode(0))
		if err != nil {
			w.flog.Println("Failed to parse source: %s", err.Error())
		}

		linecnt := 0

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			for x := 0; x < len(fn.Body.List); x++ {
				as, ok2 := fn.Body.List[x].(*ast.SendStmt)
				if !ok2 {
					continue
				}

				as3, ok3 := as.Chan.(*ast.Ident)
				if !ok3 {
					continue
				}

				if as3.Name != pk {
					continue
				}

				// figure out type

				begin := fset.Position(as.Pos()).Line - 1
				after := fset.Position(as.End()).Line + 1

				before_advice := aspect.advize.before
				after_advice := aspect.advize.after

				if before_advice != "" {
					rout = w.writeAtLine(fname, begin+linecnt, before_advice)
					linecnt += strings.Count(before_advice, "\n") + 1
				}

				if after_advice != "" {
					rout = w.writeAtLine(fname, after+linecnt-1, after_advice)

					linecnt += strings.Count(after_advice, "\n") + 1
				}

			}
		}

	}

	if len(importsNeeded) > 0 {
		// add any imports for this piece of advice
		rout = w.writeMissingImports(fname, rout, importsNeeded)
	}

	return rout
}
