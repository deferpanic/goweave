package aop

import (
	"strings"

	"go/ast"
)

// reads file && de-dupes imports
func (a *Aop) reWorkImports(fp string) string {
	flines := fileLines(fp)
	af := a.ParseAST(fp)
	pruned := a.pruneImports(af, a.rootPkg())
	lines := a.deDupeImports(fp, flines, pruned)
	a.writeOut(fp, lines)

	return lines
}

// deDupeImports de-dupes imports
func (a *Aop) deDupeImports(path string, flines []string, pruned []string) string {
	nlines := ""

	inImport := false

	for i := 0; i < len(flines); i++ {

		// see if we want to add any imports to the file
		if strings.Contains(flines[i], "import (") {
			inImport = true

			nlines += flines[i]

			for x := 0; x < len(pruned); x++ {
				nlines += pruned[x] + "\n"
			}

			continue
		}

		if inImport {
			if strings.Contains(flines[i], ")") {
				inImport = false

				nlines += ")" + "\n"
			}

			continue
		}

		// write out the original line
		nlines += flines[i] + "\n"

	}

	return nlines
}

func (a *Aop) writeMissingImports(fp string, out string, importsNeeded []string) string {

	out = a.addMissingImports(importsNeeded, out)

	a.writeOut(fp, out)

	// de-dupe imports
	return a.reWorkImports(fp)
}

// pruneImports returns a set of import strings de-duped from the ast
func (a *Aop) pruneImports(f *ast.File, rootpkg string) []string {
	pruned := []string{}

	for i := 0; i < len(f.Imports); i++ {
		if f.Imports[i].Path != nil {

			l := f.Imports[i].Path.Value

			// chk to see if this is sub-pkg - we re-write it to
			// relative path
			if strings.Contains(l, rootpkg) {
				l = rewriteImport(l, rootpkg)
			}

			if !inthere(l, pruned) {
				pruned = append(pruned, l)
			}
		}
	}

	return pruned
}

// addMissingImports adds any imports from advice that was found
func (a *Aop) addMissingImports(imports []string, out string) string {

	if strings.Contains(out, "import (") {
		s := "\n"
		for i := 0; i < len(imports); i++ {
			s += imports[i] + "\n"
		}

		out = strings.Replace(out, "import (", "import ("+s, -1)
	} else {

		s := ""
		for i := 0; i < len(imports); i++ {
			s += "import " + imports[i] + "\n"
		}

		out = strings.Replace(out, "import ", s+"import", -1)
	}

	return out
}

// rewriteImport is intended to rewrite a sub pkg of the base pkg to a
// relative path since we for now cp it to a diff. workspace
func rewriteImport(l string, rp string) string {
	return strings.Replace(l, rp, ".", -1)
}
