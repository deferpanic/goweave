package aop

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go/ast"
	"go/parser"
	"go/token"

	_ "golang.org/x/tools/go/gcimporter"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/types"
)

// Aop is struct runner for aop transforms
type Aop struct {
	flog    *log.Logger
	aspects []Aspect

	// warn if AST parsing warns you
	// off by default as many times we don't care
	warnAST bool
}

// NewAop instantiates and returns a new aop
func NewAop() *Aop {

	aop := &Aop{
		flog: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

	return aop

}

// Run preps, grabs advice, transforms the src, and builds the code
func (a *Aop) Run() {
	a.prep()
	a.loadAspects()

	// old-school regex parsing
	a.transform()

	// applys around advice && evals execution joinpoints
	filepath.Walk(a.tmpLocation(), a.VisitFile)

	a.build()

}

// VisitFile walks each file and transforms it's
// this is fairly heavy/expensive/pos right now
func (a *Aop) VisitFile(fp string, fi os.FileInfo, err error) error {
	matched, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		a.flog.Println(err)
		return err
	}

	if matched {

		af := a.ParseAST(fp)

		flines := fileLines(fp)

		pruned := a.pruneImports(af, a.rootPkg())
		lines := a.deDupeImports(fp, flines, pruned)

		// provides 'around' style advice
		stuff := a.applyAroundAdvice(fp, lines)
		a.writeOut(fp, stuff)

		// provides advice matching against execution join points
		stuff = a.applyExecutionJP(fp, stuff)
		a.writeOut(fp, stuff)

	}

	return nil
}

// Parse parses the ast for this file and returns a ParsedFile
func (a *Aop) ParseAST(fname string) *ast.File {
	var err error

	fset := token.NewFileSet()
	af, err := parser.ParseFile(fset, fname, nil, 0)
	if err != nil {
		a.flog.Println(err)
	}

	loadcfg := loader.Config{}
	loadcfg.CreateFromFilenames(fname)

	info := types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
	}

	var conf types.Config
	_, err = conf.Check(af.Name.Name, fset, []*ast.File{af}, &info)
	if err != nil {
		if a.warnAST {
			a.flog.Println(err)
		}
	}

	return af
}

// pointCutMatch returns an aspect if there is a pointcut match on this
// line or returns an empty aspect
func pointCutMatch(a []Aspect, l string) Aspect {
	for i := 0; i < len(a); i++ {

		// look for go-routines
		if strings.Contains(l, "go ") && ("go" == a[i].pointkut.def) {
			return a[i]
		}

	}

	return Aspect{}
}

// transform reads line by line over each src file and inserts advice
// where appropriate
//
// only inserts before/advice after
func (a *Aop) transform() {

	fzs := a.findGoFiles()

	rootpkg := a.rootPkg()

	for i := 0; i < len(fzs); i++ {
		out := a.txAspects(fzs[i], rootpkg)
		a.reWriteFile(fzs[i], out, []string{})
	}
}

// findGoFiles recursively finds all go files in a project
func (a *Aop) findGoFiles() []string {
	res := []string{}

	visit := func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, ".go") {
			res = append(res, path)
		}
		return nil
	}

	err := filepath.Walk(".", visit)
	if err != nil {
		a.flog.Println(err.Error())
	}

	return res
}

// findAspects finds all aspects for this project
func (a *Aop) findAspects() []string {
	aspects := []string{}

	files, _ := ioutil.ReadDir("./")
	for _, f := range files {
		if strings.Contains(f.Name(), ".weave") {
			aspects = append(aspects, f.Name())
			log.Println(f.Name())
		}
	}

	return aspects
}
