package aop

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	a.build()
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
	fmt.Println("building" + a.tmpLocation())

	fstcmd := "mkdir -p " + a.tmpLocation()
	sndcmd := `find . -type d -exec mkdir -p "` + a.tmpLocation() + `/{}" \;`

	_, err := exec.Command("bash", "-c", fstcmd).CombinedOutput()
	if err != nil {
		a.flog.Println(err.Error())
	}

	fmt.Println(sndcmd)
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

	fmt.Println(buildstr)
	_, err := exec.Command("bash", "-c", buildstr).CombinedOutput()
	if err != nil {
		a.flog.Println(err.Error())
	}

}

// pointCutMatch returns an aspect if there is a pointcut match on this
// line or returns an empty aspect
func pointCutMatch(a []Aspect, l string) Aspect {
	for i := 0; i < len(a); i++ {

		// look for functions
		if strings.Contains(l, "func "+a[i].pointkut.def) {
			return a[i]
		}

		// look for package/function
		//if strings.Contains(l, "func "+a[i].pointkut.def) {
		//		return a[i]
		//	}

	}

	return Aspect{}
}

// transform reads line by line over each src file and inserts advice
// where appropriate
func (a *Aop) transform() {

	fzs := a.findGoFiles()

	for i := 0; i < len(fzs); i++ {
		curfile := fzs[i]

		file, err := os.Open(curfile)
		if err != nil {
			a.flog.Println(err)
		}
		defer file.Close()

		out := ""

		// poor man's scope
		scope := 0

		cur_aspect := Aspect{}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			l := scanner.Text()

			newAspect := pointCutMatch(a.aspects, l)
			if newAspect.pointkut.def != "" {
				scope += 1

				cur_aspect = newAspect

				// before advice
				if cur_aspect.advize.before != "" {
					out += l + "\n" + cur_aspect.advize.before + "\n"
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

		b, err := f.WriteString(out)
		fmt.Println(b)
		if err != nil {
			a.flog.Println(err)
		}

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
		if strings.Contains(f.Name(), ".goa") {
			aspects = append(aspects, f.Name())
			fmt.Println(f.Name())
		}
	}

	return aspects
}
