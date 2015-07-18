package weave

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
)

// applyGlobalAdvice applies any global scoped advice
// if advice has already been placed in this pkg than we skip
// no support for sub-pkgs yet
// FIXME
func (w *Weave) applyGlobalAdvice(fname string, stuff string) string {
	if w.appliedGlobal {
		return stuff
	}

	rout := stuff

	importsNeeded := []string{}

	for i := 0; i < len(w.aspects); i++ {
		aspect := w.aspects[i]
		if aspect.pointkut.def != "*" {
			continue
		}

		before_advice := aspect.advize.before
		after_advice := aspect.advize.after

		w.appliedGlobal = true

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, fname, rout, parser.Mode(0))
		if err != nil {
			w.flog.Println("Failed to parse source: %s", err.Error())
		}

		// endOfImport returns the line after the import statement
		// used to add global advice
		// just a guess here - doubt this is ordered
		ilen := len(file.Imports)
		s := file.Imports[ilen-1].End()
		ei := fset.Position(s).Line + 1

		if before_advice != "" {
			rout = w.writeAtLine(fname, ei, before_advice)
		}

		if after_advice != "" {
			rout = w.writeAtLine(fname, ei, after_advice)
		}

	}

	if len(importsNeeded) > 0 {
		rout = w.writeMissingImports(fname, rout, importsNeeded)
	}

	return rout
}

// applyCallAdvice applies call advice before/after a call
// around advice is currently hacked in via applyAroundAdvice
func (w *Weave) applyCallAdvice(fname string, stuff string) string {

	rout := stuff

	importsNeeded := []string{}

	for i := 0; i < len(w.aspects); i++ {
		aspect := w.aspects[i]
		if aspect.pointkut.kind != 1 {
			continue
		}

		linecnt := 0

		pk := aspect.pointkut.def
		before_advice := aspect.advize.before
		after_advice := aspect.advize.after

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, fname, rout, parser.Mode(0))
		if err != nil {
			w.flog.Println("Failed to parse source: %s", err.Error())
		}

		// look for call expressions - call joinpoints

		// bend keeps track of tailing binaryExprs
		// if a call-expression's end is not > then we use this
		lastbinstart := 0
		lastbinend := 0

		ast.Inspect(file, func(n ast.Node) bool {

			switch stmt := n.(type) {

			case *ast.CompositeLit:
				// log.Printf("found a composite lit\n")
				lbs := fset.Position(stmt.Lbrace).Line
				lbe := fset.Position(stmt.Rbrace).Line

				if lbs > lastbinstart {
					lastbinstart = lbs
				}

				if lbe > lastbinend {
					lastbinend = lbe
				}

			case *ast.BinaryExpr:
				// child nodes traversed in DFS - so we want the max
				// it's ok when newer nodes re-write this out of scope
				lbs := fset.Position(stmt.X.Pos()).Line
				lbe := fset.Position(stmt.Y.Pos()).Line

				if lbs > lastbinstart {
					lastbinstart = lbs
				}

				if lbe > lastbinend {
					lastbinend = lbe
				}

				// log.Printf("last bin start %d, last bin end %d", lastbinstart, lastbinend)

			case *ast.CallExpr:
				var name string
				switch x := stmt.Fun.(type) {
				case *ast.Ident:
					name = x.Name
				case *ast.SelectorExpr:
					name = x.Sel.Name
				default:
					name = "WRONG"
				}

				pk = strings.Split(pk, "(")[0]

				// fixme - hack - we need support for identifying call
				// expressions w/pkgs
				if strings.Contains(pk, ".") {
					pk = strings.Split(pk, ".")[1]
				}

				// are we part of a bigger expression?
				// if so grab our lines so we don't erroneously set them
				if (string(name) != pk) && (len(stmt.Args) > 1) {
					// log.Printf("found comma..")
					// child nodes traversed in DFS - so we want the max
					// it's ok when newer nodes re-write this out of scope
					lbs := fset.Position(stmt.Lparen).Line
					lbe := fset.Position(stmt.Rparen).Line
					// log.Printf("lbs: %d, lbe: %d\n", lbs, lbe)
					if lbs > 0 {
						//lastbinstart {
						lastbinstart = lbs
					}

					if lbe > lastbinend {
						lastbinend = lbe
					}

				}

				if string(name) == pk {

					begin := fset.Position(stmt.Lparen).Line
					end := fset.Position(stmt.Rparen).Line
					// log.Printf("found expression on line %d\n", begin)
					//log.Printf("lbs: %d, lbe: %d\n", lbs, lbe)

					// adjust begin
					if begin > lastbinend {
						// log.Printf("using this funcs start %d", begin)
					} else {
						if lastbinstart < begin {
							begin = lastbinstart
						}
						// log.Printf("using binexps' start %d", begin)
					}

					if end > lastbinend {
						// log.Printf("using this funcs begin %d", begin)
					} else {
						if lastbinstart < begin {
							begin = lastbinstart
						}
						// log.Printf("using binexps' start %d", begin)
					}

					// adjust end
					if lastbinend > end {
						end = lastbinend
					}

					if before_advice != "" {
						rout = w.writeAtLine(fname, begin+linecnt-1, before_advice)
						// log.Println(rout)
						// log.Printf("writing at line %d", begin+linecnt-1)
						linecnt += strings.Count(before_advice, "\n") + 1
					}

					if after_advice != "" {
						rout = w.writeAtLine(fname, end+linecnt, after_advice)
						// log.Println(rout)
						// log.Printf("writing at line %d", end+linecnt)

						linecnt += strings.Count(after_advice, "\n") + 1
					}

					for t := 0; t < len(aspect.importz); t++ {
						importsNeeded = append(importsNeeded, aspect.importz[t])
					}

				}
			}

			return true
		})

	}

	if len(importsNeeded) > 0 {
		rout = w.writeMissingImports(fname, rout, importsNeeded)
	}

	return rout

}

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
		if aspect.pointkut.kind != 2 {
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

			if fn.Name.Name == fpk && containArgs(pk, fn) {

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
						lcnt := strings.Count(after_advice, "\n") + 1
						rout = w.writeAtLine(fname, after+linecnt-1-lcnt, after_advice)
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
