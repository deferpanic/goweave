package weave

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
)

// applyAroundAdvice uses code from gofmt to wrap any after advice
// essentially this is the same stuff you could do w/the cmdline tool,
// gofmt
//
// FIXME - mv to CallExpr
//
// looks for call joinpoints && provides around advice capability
//
// this is currently a hack to support deferpanic's http lib
func (w *Weave) applyAroundAdvice(fname string) string {

	stuff := fileAsStr(fname)

	importsNeeded := []string{}

	for i := 0; i < len(w.aspects); i++ {
		aspect := w.aspects[i]
		if aspect.advize.around != "" {
			pk := aspect.pointkut
			around_advice := aspect.advize.around

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, fname, stuff, parser.Mode(0))
			if err != nil {
				w.flog.Println("Failed to parse source: %s", err.Error())
			}

			// needs match groups
			// wildcards of d,s...etc.
			p := parseExpr(pk.def)
			val := parseExpr(around_advice)

			file = rewriteFile2(p, val, file)

			buf := new(bytes.Buffer)
			err = format.Node(buf, fset, file)
			if err != nil {
				w.flog.Println("Failed to format post-replace source: %v", err)
			}

			actual := buf.String()

			w.writeOut(fname, actual)

			stuff = actual

			for t := 0; t < len(aspect.importz); t++ {
				importsNeeded = append(importsNeeded, aspect.importz[t])
			}

		}
	}

	if len(importsNeeded) > 0 {
		// add any imports for this piece of advice
		stuff = w.writeMissingImports(fname, stuff, importsNeeded)
	}

	return stuff
}

// applyExecutionJP applies any advice for execution joinpoints
func (w *Weave) applyExecutionJP(fname string, stuff string) string {

	rout := stuff

	importsNeeded := []string{}

	for i := 0; i < len(w.aspects); i++ {
		aspect := w.aspects[i]
		if !(aspect.pointkut.kind > 0) {
			continue
		}

		pk := aspect.pointkut.def

		before_advice := aspect.advize.before
		after_advice := aspect.advize.after

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

				// begin line
				begin := fset.Position(fn.Body.Lbrace).Line
				after := fset.Position(fn.Body.Rbrace).Line

				// until this is refactored - any lines we add in our
				// advice need to be accounted for w/begin

				if before_advice != "" {
					rout = w.writeAtLine(fname, begin+linecnt, before_advice)
					linecnt += strings.Count(before_advice, "\n") + 1
				}

				if after_advice != "" {
					if fn.Type.Results != nil {
						rout = w.writeAtLine(fname, after+linecnt-2, after_advice)
					} else {
						rout = w.writeAtLine(fname, after+linecnt-1, after_advice)
					}

					linecnt += strings.Count(after_advice, "\n") + 1
				}

				for t := 0; t < len(aspect.importz); t++ {
					importsNeeded = append(importsNeeded, aspect.importz[t])
				}

			}
		}

	}

	if len(importsNeeded) > 0 {
		// add any imports for this piece of advice applyExecutionJP
		rout = w.writeMissingImports(fname, rout, importsNeeded)
	}

	return rout
}
