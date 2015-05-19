package aop

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

// adviceType returns a map of id to human expression of advice types
func adviceType() map[int]string {
	return map[int]string{
		1: "before",
		2: "after",
		3: "around",
	}
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
	out, err := exec.Command("mkdir", "-p", a.tmpLocation()).CombinedOutput()
	if err != nil {
		a.flog.Println(err.Error())
	}
	fmt.Printf("%s\n", out)
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
		if strings.Contains(l, a[i].pointkut.def) {
			return a[i]
		}
	}

	return Aspect{}
}

// transform reads line by line over each src file and inserts advice
// where appropriate
func (a *Aop) transform() {
	curfile := "main.go"

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

			fmt.Println("pointcut def:\t" + newAspect.pointkut.def)
			fmt.Println("pointcut function:\t" + newAspect.pointkut.funktion)
			fmt.Println("before advice:\t" + newAspect.advize.before)
			fmt.Println("after advice:\t" + newAspect.advize.after)

			cur_aspect = newAspect

			// before advice
			if (cur_aspect.advize.adviceTypeId == 1) ||
				(cur_aspect.advize.adviceTypeId == 3) {
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

// adviceKind returns the map id of human expression of advice type
func adviceKind(l string) int {
	stuff := strings.Split(l, ": ")
	ostuff := strings.Split(stuff[1], " ")

	switch ostuff[0] {
	case "before":
		return 1
	case "after":
		return 2
	case "around":
		return 3
	}

	return -1
}
