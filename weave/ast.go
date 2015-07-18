package weave

import (
	"go/ast"
	"go/parser"
	"log"
	"strings"
)

// stolen from http://golang.org/src/cmd/fix/fix.go
func isPkgDot(expr ast.Expr, pkg, name string) bool {
	if len(pkg) > 0 {
		sel, ok := expr.(*ast.SelectorExpr)
		return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
	} else {
		return isIdent(expr, name)
	}
}

// stolen from http://golang.org/src/cmd/fix/fix.go
func isPtrPkgDot(t ast.Expr, pkg, name string) bool {
	ptr, ok := t.(*ast.StarExpr)
	return ok && isPkgDot(ptr.X, pkg, name)
}

// stolen from http://golang.org/src/cmd/fix/fix.go
func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}

// parseExpr returns an ast expression from the source s
// stolen from http://golang.org/src/cmd/fix/fix.go
func parseExpr(s string) ast.Expr {
	exp, err := parser.ParseExpr(s)
	if err != nil {
		log.Println("Cannot parse expression %s :%s", s, err.Error())
	}
	return exp
}

// containArgs ensures the function signature is 'correct'
// the following is left un-implemented and needs to be implemented
// FIXME
// 1) simple types - no pkgs
// 2) order of arguments
// 3) no args
// 4) no simple args
func containArgs(pk string, fn *ast.FuncDecl) bool {

	/*
		this function will return a channel through which
		we can get every argument's type of pfn
	*/
	nextarg := func(pfn *ast.FuncDecl, stop chan int) <-chan ast.Expr {
		argTypeChnl := make(chan ast.Expr)
		go func() {
			defer close(argTypeChnl)
			for _, arg := range pfn.Type.Params.List {
				for range arg.Names {
					select {
					case argTypeChnl <- arg.Type:
					case <-stop:
						return
					}
				}
			}
		}()
		return argTypeChnl
	}

	stop := make(chan int)
	defer close(stop)

	argTypeChnl := nextarg(fn, stop)

	//--------------------

	pk = strings.Split(pk, "(")[1]
	pk = strings.Split(pk, ")")[0]

	arglist := strings.Split(pk, ",")

	if (len(arglist) == 1) && (arglist[0] == "") {
		arglist = []string{}
	}

	// Check whether every argument's type is the same
	for _, argtype := range arglist {

		typelist := strings.Split(argtype, ".")
		isptr := (argtype[0] == '*')

		pkg := ""
		name := typelist[0]

		compareArgType := isPkgDot
		if isptr {
			compareArgType = isPtrPkgDot
		}

		// pkg.type or *pkg.type
		if len(typelist) == 2 {
			pkg = typelist[0]
			name = typelist[1]
			if isptr {
				pkg = pkg[1:]
			}
		}

		nextArgType, ok := <-argTypeChnl

		if !ok || !compareArgType(nextArgType, pkg, name) {
			return false
		}
	}

	return true
}
