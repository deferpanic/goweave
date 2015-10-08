package weave

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// mAvars is a list of AbstractVars found
// prob. needs to be refactored
var mAvars []AbstractVar

// AbstractVar contains the value of a var and it's type
// Val should prob. be interface
// refactor refactor refactor
type AbstractVar struct {
	Kind string
	Val  string
}

// applySetJP applies any advice for set joinpoints
// currently expects a channel type
// need some extra help here to be agnostic
// maybe some helper functions that determine what type it is and
// associated meta-data about it
//
// this is currently a TEST - it only works for specific channels at the
// moment
func (w *Weave) applySetJP(fname string, stuff string) string {

	rout := stuff

	importsNeeded := []string{}

	for i := 0; i < len(w.aspects); i++ {

		aspect := w.aspects[i]
		if aspect.pointkut.kind != 5 {
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
				as, ok2 := fn.Body.List[x].(*ast.AssignStmt)
				if !ok2 {
					continue
				}

				// ll-cool-j
				blah := as.Lhs[0].(*ast.Ident).Name

				if pk != blah {
					continue
				}

				// figure out type

				// ll-cool-j once again
				if len(as.Rhs[0].(*ast.CallExpr).Args) == 2 {
					_, k := (as.Rhs[0].(*ast.CallExpr).Args[0]).(*ast.ChanType)
					if !k {
						continue
					}

					r2, k2 := (as.Rhs[0].(*ast.CallExpr).Args[1]).(*ast.BasicLit)
					if !k2 {
						continue
					}

					avar := AbstractVar{
						Kind: "channel",
						Val:  r2.Value,
					}

					begin := fset.Position(as.Pos()).Line - 1
					after := fset.Position(as.End()).Line + 1

					// this needs to do var subsitution like we do for
					// within advice
					//
					// before_advice := formatAdvice(aspect.advize.before, mName)
					// after_advice := formatAdvice(aspect.advize.after, mName)

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

					mAvars = append(mAvars, avar)

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
