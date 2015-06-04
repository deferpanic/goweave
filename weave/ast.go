package weave

import (
	"go/ast"
	"go/parser"
	"log"
	"strings"
)

func isPkgDot(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
}

func isPtrPkgDot(t ast.Expr, pkg, name string) bool {
	ptr, ok := t.(*ast.StarExpr)
	return ok && isPkgDot(ptr.X, pkg, name)
}

func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}

// parseExpr returns an ast expression from the source s
func parseExpr(s string) ast.Expr {
	exp, err := parser.ParseExpr(s)
	if err != nil {
		log.Println("Cannot parse expression %s :%s", s, err.Error())
	}
	return exp
}

// containArgs ensures the function signature is 'correct'
// this currently is erroneous cause it doesn't look for order of
// arguments && it ignores simple strings -- FIXME
// 1) simple types - no pkgs
// 2) order of arguments
// 3) no args
func containArgs(pk string, p []*ast.Field) bool {

	pk = strings.Split(pk, "(")[1]
	pk = strings.Split(pk, ")")[0]

	argz := strings.Split(pk, ",")

	if (len(argz) == 1) && (argz[0] == "") {
		argz = []string{}
	}

	// early bail if mis-matched argc
	if len(argz) != len(p) {
		return false
	}

	xtrue := 0

	// for now we ignore simple args like string, int
	// also - these are un-ordered right now..
	// also - no support for no args
	for i := 0; i < len(argz); i++ {
		s := strings.Split(argz[i], ".")
		pkg := strings.TrimSpace(s[0])
		iname := strings.TrimSpace(s[1])

		if strings.Contains(pkg, "*") {
			pkg = strings.Replace(pkg, "*", "", -1)
			for _, field := range p {
				if isPtrPkgDot(field.Type, pkg, iname) {
					xtrue += 1
				}
			}

		} else {
			for _, field := range p {

				if isPkgDot(field.Type, pkg, iname) {
					xtrue += 1
				}
			}
		}
	}

	if xtrue == len(argz) {
		return true
	}

	return false
}
