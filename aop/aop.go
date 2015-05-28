package aop

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"go/ast"
	"go/format"
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
	a.transform()

	filepath.Walk(a.tmpLocation(), a.VisitFile)

	a.build()

}

func parseExpr(s string) ast.Expr {
	exp, err := parser.ParseExpr(s)
	if err != nil {
		log.Println("Cannot parse expression %s :%s", s, err.Error())
	}
	return exp
}

// txAfter uses code from gofmt to wrap any after advice
// essentially this is the same stuff you could do w/the cmdline tool,
// gofmt
func (a *Aop) txAfter(fname string, lines string) string {

	stuff := lines

	for i := 0; i < len(a.aspects); i++ {
		aspect := a.aspects[i]
		if aspect.advize.around != "" {
			pk := aspect.pointkut
			around_advice := aspect.advize.around

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, fname, lines, parser.Mode(0))
			if err != nil {
				log.Println("Failed to parse source: %s", err.Error())
			}

			// needs match groups
			// wildcards of d,s...etc.
			p := parseExpr(pk.def)
			val := parseExpr(around_advice)

			file = rewriteFile2(p, val, file)

			buf := new(bytes.Buffer)
			err = format.Node(buf, fset, file)
			if err != nil {
				log.Println("Failed to format post-replace source: %v", err)
			}

			actual := buf.String()

			a.writeOut(fname, actual)

			stuff = actual
		}
	}

	return stuff
}

func (a *Aop) VisitFile(fp string, fi os.FileInfo, err error) error {
	matched, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		fmt.Println(err)
		return err
	}

	if matched {
		fmt.Println("looking at " + fp)

		af := a.ParseAST(fp)

		flines := fileLines(fp)
		pruned := pruneImports(af)
		lines := a.deDupeImports(fp, flines, pruned)

		stuff := a.txAfter(fp, lines)

		a.writeOut(fp, stuff)

	}

	return nil
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

// writeOut writes nlines to path
func (a *Aop) writeOut(path string, nlines string) {

	b := []byte(nlines)
	err := ioutil.WriteFile(path, b, 0644)
	if err != nil {
		log.Println(err)
	}
}

// pruneImports
func pruneImports(f *ast.File) []string {

	pruned := []string{}

	for i := 0; i < len(f.Imports); i++ {
		if f.Imports[i].Path != nil {
			if !inthere(f.Imports[i].Path.Value, pruned) {
				pruned = append(pruned, f.Imports[i].Path.Value)
			}
		}
	}

	return pruned
}

func inthere(p string, ray []string) bool {
	for i := 0; i < len(ray); i++ {
		if ray[i] == p {
			return true
		}
	}

	return false
}

// Parse parses the ast for this file and returns a ParsedFile
func (a *Aop) ParseAST(fname string) *ast.File {
	var err error

	fset := token.NewFileSet()
	af, err := parser.ParseFile(fset, fname, nil, 0)
	if err != nil {
		panic(err)
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
		fmt.Println(err)
	}

	return af
}

// buildDir determines what the root build dir is
func (a *Aop) buildDir() string {
	out, err := exec.Command("bash", "-c", "pwd").CombinedOutput()
	if err != nil {
		a.flog.Println(err.Error())
	}

	return strings.TrimSpace(string(out))
}

// binName returns the expected bin name
func (a *Aop) binName() string {
	s := a.buildDir()
	stuff := strings.Split(s, "/")
	return stuff[len(stuff)-1]
}

// prep prepares any tmp. build dirs
func (a *Aop) prep() {

	fstcmd := "mkdir -p " + a.tmpLocation()
	sndcmd := `find . -type d -exec mkdir -p "` + a.tmpLocation() + `/{}" \;`

	_, err := exec.Command("bash", "-c", fstcmd).CombinedOutput()
	if err != nil {
		a.flog.Println(err.Error())
	}

	_, err = exec.Command("bash", "-c", sndcmd).CombinedOutput()
	if err != nil {
		a.flog.Println(err.Error())
	}

}

// whichgo determines provides the full go path to the current go build
// tool
func (a *Aop) whichGo() string {
	out, err := exec.Command("bash", "-c", "which go").CombinedOutput()
	if err != nil {
		a.flog.Println(err.Error())
	}

	return strings.TrimSpace(string(out))
}

// tmpLocation returns the tmp build dir
func (a *Aop) tmpLocation() string {
	return "/tmp" + a.buildDir()
}

// build does the actual compilation
// right nowe we piggy back off of 6g/8g
func (a *Aop) build() {
	buildstr := "cd " + a.tmpLocation() + " && " + a.whichGo() + " build && cp " +
		a.binName() + " " + a.buildDir() + "/."

	o, err := exec.Command("bash", "-c", buildstr).CombinedOutput()
	if err != nil {
		a.flog.Println(string(o))
	}

}

// pointCutMatch returns an aspect if there is a pointcut match on this
// line or returns an empty aspect
func pointCutMatch(a []Aspect, l string) Aspect {
	for i := 0; i < len(a); i++ {

		// look for exact functions
		if strings.Contains(l, "func "+a[i].pointkut.def) {
			return a[i]
		}

		// look for go-routines
		if strings.Contains(l, "go ") && ("go" == a[i].pointkut.def) {
			return a[i]
		}

	}

	return Aspect{}
}

// returns true if this is a multi-line go routine
func multiLineGo(l string) bool {
	if strings.Contains(l, "go func(") {
		return true
	}

	return false
}

// returns true if this is a single line go routine
func singleLineGo(l string) bool {
	var singlelinego = regexp.MustCompile(`go\s.*\(.*\)`)
	if singlelinego.MatchString(l) {
		return true
	}

	return false
}

// returns a slice of lines
func fileLines(path string) []string {
	stuff := []string{}

	file, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		stuff = append(stuff, scanner.Text())
	}

	return stuff
}

// transform reads line by line over each src file and inserts advice
// where appropriate
//
// only inserts before/advice after
func (a *Aop) transform() {

	fzs := a.findGoFiles()

	rootpkg := a.rootPkg()

	for i := 0; i < len(fzs); i++ {
		curfile := fzs[i]
		importsNeeded := []string{}

		file, err := os.Open(curfile)
		if err != nil {
			a.flog.Println(err)
		}
		defer file.Close()

		out := ""

		// poor man's scope
		scope := 0

		// poor man's import parsing
		inImport := false

		cur_aspect := Aspect{}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			l := scanner.Text()

			// fix me - we can get these from the AST
			if a.importBlock(l) || inImport {
				inImport = true

				if strings.Contains(l, "\"") {

					if strings.Contains(l, rootpkg) {
						l = a.rewriteImport(l, rootpkg)
					}
				}
			}

			// close us out of import block if we are done
			if inImport {
				if strings.Contains(l, ")") {
					inImport = false
				}
			}

			newAspect := pointCutMatch(a.aspects, l)
			if newAspect.pointkut.def != "" {
				scope += 1

				// insert any imports we need to
				for x := 0; x < len(newAspect.importz); x++ {
					importsNeeded = append(importsNeeded, newAspect.importz[x])
				}

				cur_aspect = newAspect

				if cur_aspect.advize.before != "" {

					// go before advice
					if cur_aspect.pointkut.def == "go" {
						if multiLineGo(l) {

							// keep grabbing lines until we are back to
							// existing scope?
							stuff := ""
							nscope := 1
							for i := 0; ; i++ {
								scanner.Scan()
								l2 := scanner.Text()

								fmt.Println(l2)

								if strings.Contains(l2, "{") {
									nscope += 1
								}

								if strings.Contains(l2, "}") {
									nscope -= 1
								}

								if nscope == 0 {
									break
								}

								stuff += l2 + "\n"

							}

							out += "go func(){\n" + cur_aspect.advize.before + "\n" + stuff +
								"\n" + "}()\n"

						} else if singleLineGo(l) {

							// hack - ASTize me
							r := regexp.MustCompile("go\\s(.*)\\((.*)\\)")

							newstr := r.ReplaceAllString(l, "go func(){\n"+
								cur_aspect.advize.before+"\n$1($2)\n"+"}()")

							out += newstr + "\n"

						}
					} else {
						// normal before advice
						out += l + "\n" + cur_aspect.advize.before + "\n"
					}

					continue
				}

			}

			// dat scope
			if strings.Contains(l, "}") || strings.Contains(l, "return") {

				scope -= 1

				out += cur_aspect.advize.after + "\n"
			}

			out += l + "\n"

		}

		if err := scanner.Err(); err != nil {
			a.flog.Println(err)
		}

		f, err := os.Create(a.tmpLocation() + "/" + curfile)
		if err != nil {
			a.flog.Println(err)
		}

		defer f.Close()

		out = a.addMissingImports(importsNeeded, out)

		b, err := f.WriteString(out)
		fmt.Println(b)
		if err != nil {
			a.flog.Println(err)
		}

	}
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
func (a *Aop) rewriteImport(l string, rp string) string {
	return strings.Replace(l, rp, ".", -1)
}

// importBlock detects if we are in an import statement or block
func (a *Aop) importBlock(l string) bool {
	if strings.Contains(l, "import") {
		return true
	} else {
		return false
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

// rootPkg returns the root package of a go build
// this is needed to determine whether or not sub-pkg imports need to be
// re-written - which is basically any project w/more than one folder
func (a *Aop) rootPkg() string {
	out, err := exec.Command("bash", "-c", "go list").CombinedOutput()
	if err != nil {
		a.flog.Println(err.Error())
	}

	return strings.TrimSpace(string(out))
}

// findAspects finds all aspects for this project
func (a *Aop) findAspects() []string {
	aspects := []string{}

	files, _ := ioutil.ReadDir("./")
	for _, f := range files {
		if strings.Contains(f.Name(), ".weave") {
			aspects = append(aspects, f.Name())
			fmt.Println(f.Name())
		}
	}

	return aspects
}
