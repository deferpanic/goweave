package weave

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

// Weave is struct runner for aspect transforms
type Weave struct {
	flog    *log.Logger
	aspects []Aspect

	// warn if AST parsing warns you
	// off by default as many times we don't care
	warnAST bool

	// the pkg where we run goweave
	basePkg string

	// build location is where weave our aspects
	buildLocation string

	// verbose debugging output
	verbose bool

	// appliedGlobal is a HACK - the goal is to keep track of whether or
	// not we have kept tracked of advice in a pkg so far - this is
	// currently not working as expected as sub-pkgs are ignored
	appliedGlobal bool
}

// NewWeave instantiates and returns a new aop
func NewWeave() *Weave {

	w := &Weave{
		flog:          log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
		basePkg:       setBase(),
		buildLocation: tmpLocation(),
	}

	return w

}

// Run preps, grabs advice, transforms the src, and builds the code
func (w *Weave) Run() {
	w.prep()
	w.loadAspects()

	// old regex parsing
	// only used for go routines currently
	// soon to be DEPRECATED
	w.transform()

	// applys around advice && evals execution joinpoints
	filepath.Walk(w.buildLocation, w.VisitFile)

	w.build()

}

// VisitFile walks each file and transforms it's
// this is fairly heavy/expensive/pos right now
func (w *Weave) VisitFile(fp string, fi os.FileInfo, err error) error {

	if fi == nil || err != nil {
		return err
	}

	matched, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		w.flog.Println(err)
		return err
	}

	if matched {

		// provides 'around' style advice
		// HACK - should be moved to call advice
		stuff := w.applyAroundAdvice(fp)
		w.writeOut(fp, stuff)

		// any global advice ? stick it in this file
		// FIXME
		stuff = w.applyGlobalAdvice(fp, stuff)
		w.writeOut(fp, stuff)

		// provides advice matching against call join points
		stuff = w.applyCallAdvice(fp, stuff)
		w.writeOut(fp, stuff)

		// provides advice matching against execution join points
		stuff = w.applyExecutionJP(fp, stuff)
		w.writeOut(fp, stuff)

		// provides advice matching against within join points
		stuff = w.applyWithinJP(fp, stuff)
		w.writeOut(fp, stuff)

		// provides advice matching against declaration join points
		stuff = w.applyDeclarationJP(fp, stuff)
		w.writeOut(fp, stuff)

		// provides advice matching against set join points
		stuff = w.applySetJP(fp, stuff)
		w.writeOut(fp, stuff)

		// provides advice matching against get join points
		stuff = w.applyGetJP(fp, stuff)
		w.writeOut(fp, stuff)

		// finally re-work imports on each
		// to ensure files we didn't apply aspects to have correct
		// imports (eg: main.go)
		w.reWorkImports(fp)

	}

	return nil
}

// Parse parses the ast for this file and returns a ParsedFile
func (w *Weave) ParseAST(fname string) *ast.File {
	var err error

	fset := token.NewFileSet()
	af, err := parser.ParseFile(fset, fname, nil, 0)
	if err != nil {
		w.flog.Println(err)
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
		if w.warnAST {
			w.flog.Println(err)
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
func (w *Weave) transform() {

	fzs := w.findGoFiles()

	rootpkg := w.rootPkg()

	for i := 0; i < len(fzs); i++ {
		out, b := w.processGoRoutines(fzs[i], rootpkg)

		// FIXME
		if b {
			w.reWriteFile(fzs[i], out, []string{})
		}
	}

}

// findGoFiles recursively finds all go files in a project
func (w *Weave) findGoFiles() []string {
	res := []string{}

	visit := func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, ".go") {
			res = append(res, path)
		}
		return nil
	}

	err := filepath.Walk(".", visit)
	if err != nil {
		w.flog.Println(err.Error())
	}

	return res
}

// findAspects finds all aspects for this project
func (w *Weave) findAspects() []string {
	aspects := []string{}

	files, _ := ioutil.ReadDir("./")
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".weave") {
			aspects = append(aspects, f.Name())
			log.Println(f.Name())
		}
	}

	return aspects
}
