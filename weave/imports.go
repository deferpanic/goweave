// this file needs to be cleaned up hardcore
// we don't need 10 different ways of doing the same thing

package weave

import (
	"fmt"
	"strings"

	"go/ast"
)

type weaveImport struct {
	path string
	nsa  bool
}

// reads file && de-dupes imports
func (w *Weave) reWorkImports(fp string) string {
	flines := fileLines(fp)
	af := w.ParseAST(fp)

	fmt.Println("reworking " + fp)

	pruned := w.pruneImports(af, w.rootPkg())
	lines := w.deDupeImports(fp, flines, pruned)
	w.writeOut(fp, lines)

	return lines
}

// deDupeImports de-dupes imports
// this is txt processing - not from the ast - FIXME
func (w *Weave) deDupeImports(path string, flines []string, pruned []*ast.ImportSpec) string {
	nlines := ""

	inImport := false

	for i := 0; i < len(flines); i++ {

		// see if we want to add any imports to the file
		if strings.Contains(flines[i], "import (") {
			inImport = true

			nlines += flines[i]

			for x := 0; x < len(pruned); x++ {
				if pruned[x].Name != nil {
					nlines += pruned[x].Name.Name + " " + pruned[x].Path.Value + "\n"
				} else {
					nlines += pruned[x].Path.Value + "\n"
				}
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

// writeMissingImports writes body && any missing imports to fp
func (w *Weave) writeMissingImports(fp string, out string, importsNeeded []string) string {

	out = w.addMissingImports(importsNeeded, out)

	w.writeOut(fp, out)

	// de-dupe imports
	return w.reWorkImports(fp)
}

// pruneImports returns a set of import strings de-duped from the ast
func (w *Weave) pruneImports(f *ast.File, rootpkg string) []*ast.ImportSpec {
	pruned := []*ast.ImportSpec{}

	for i := 0; i < len(f.Imports); i++ {
		if f.Imports[i].Path != nil {

			l := f.Imports[i].Path.Value

			if strings.Contains(l, rootpkg) && !strings.Contains(l, "_weave") {
				f.Imports[i].Path.Value = rewriteImport(l, rootpkg)
				fmt.Println("rewrote improt to " + f.Imports[i].Path.Value)
			}

			if !inthere(f.Imports[i].Path.Value, pruned) {
				pruned = append(pruned, f.Imports[i])
			}
		}
	}

	return pruned
}

// inthere returns true if p is part of ray
func inthere(p string, ray []*ast.ImportSpec) bool {
	for i := 0; i < len(ray); i++ {
		if ray[i].Path.Value == p {
			return true
		}
	}

	return false
}

// addMissingImports adds any imports from advice that was found
func (w *Weave) addMissingImports(imports []string, out string) string {

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
	return "\"_weave/" + l[1:len(l)]
}
