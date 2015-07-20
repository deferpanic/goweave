package weave

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// applyWithinJP applies any advice for within joinpoints
// right now this expects both 'before && after'
func (w *Weave) applyWithinJP(fname string, stuff string) string {

	rout := stuff

	importsNeeded := []string{}

	for i := 0; i < len(w.aspects); i++ {
		aspect := w.aspects[i]
		if aspect.pointkut.kind != 3 {
			continue
		}

		pk := aspect.pointkut.def

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, fname, rout, parser.Mode(0))
		if err != nil {
			w.flog.Println("Failed to parse source: %s", err.Error())
		}

		linecnt := 0

		// look for function declarations - ala look for execution
		// joinpoints
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			fpk := strings.Split(pk, "(")[0]

			// if function name missing --> wildcard
			if fpk == "" {
				fpk = fn.Name.Name
			}

			if fn.Name.Name == fpk && containArgs(pk, fn.Type.Params.List) {

				wb := WithinBlock{
					name:          fn.Name.Name,
					fname:         fname,
					stmts:         fn.Body.List,
					linecnt:       linecnt,
					importsNeeded: importsNeeded,
					aspect:        aspect,
					fset:          fset,
				}

				rout, linecnt, importsNeeded = wb.iterateBodyStatements(w)

			}
		}

	}

	if len(importsNeeded) > 0 {
		// add any imports for this piece of advice
		rout = w.writeMissingImports(fname, rout, importsNeeded)
	}

	return rout
}

// WithinBlock contains all the info to perform a withinBlock pointcut
type WithinBlock struct {
	name          string
	fname         string
	stmts         []ast.Stmt
	linecnt       int
	importsNeeded []string
	aspect        Aspect
	fset          *token.FileSet
}

// iterateBodyStatements places within advice when a joinpoint is found
func (wb *WithinBlock) iterateBodyStatements(w *Weave) (string, int, []string) {
	rout := ""

	for i := 0; i < len(wb.stmts); i++ {

		rout = wb.insertInWithin(wb.stmts[i], w)
	}

	return rout, wb.linecnt, wb.importsNeeded
}

// HACK HACK HACK
func grabMethodName(a ast.Stmt) string {
	es := a.(*ast.ExprStmt)

	s, ok := es.X.(*ast.CallExpr)
	if !ok {
		return ""
	} else {
		e, ok := s.Fun.(*ast.Ident)
		if !ok {
			return ""
		} else {
			return e.Name
		}
	}
}

// insertInWithin places before/after advice around a statement
func (wb *WithinBlock) insertInWithin(a ast.Stmt, w *Weave) string {
	rout := ""

	mName := grabMethodName(a)

	// begin line
	begin := wb.fset.Position(a.Pos()).Line - 1
	after := wb.fset.Position(a.End()).Line + 1

	// until this is refactored - any lines we add in our
	// advice need to be accounted for w/begin
	before_advice := formatAdvice(wb.aspect.advize.before, mName)
	after_advice := formatAdvice(wb.aspect.advize.after, mName)

	if before_advice != "" {
		rout = w.writeAtLine(wb.fname, begin+wb.linecnt, before_advice)
		wb.linecnt += strings.Count(before_advice, "\n") + 1
	}

	if after_advice != "" {
		rout = w.writeAtLine(wb.fname, after+wb.linecnt-1, after_advice)

		wb.linecnt += strings.Count(after_advice, "\n") + 1
	}

	for t := 0; t < len(wb.aspect.importz); t++ {
		wb.importsNeeded = append(wb.importsNeeded, wb.aspect.importz[t])
	}

	return rout
}

// formatAdvice subsitutes any reserved keywords
// currently supported is mName
// mName is the currently called method name
func formatAdvice(advice string, mName string) string {
	return strings.Replace(advice, "mName", "\""+mName+"\"", -1)
}
