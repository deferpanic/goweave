package aop

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
)

// reWriteFile rewrites curfile with out && adds any missing imports
func (a *Aop) reWriteFile(curfile string, out string, importsNeeded []string) {

	f, err := os.Create(a.tmpLocation() + "/" + curfile)
	if err != nil {
		a.flog.Println(err)
	}

	defer f.Close()

	out = a.addMissingImports(importsNeeded, out)

	_, err = f.WriteString(out)
	if err != nil {
		a.flog.Println(err)
	}
}

// returns a slice of lines
func fileLines(path string) []string {
	stuff := []string{}

	file, err := os.Open(path)

	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		stuff = append(stuff, scanner.Text())
	}

	return stuff
}

// writeOut writes nlines to path
func (a *Aop) writeOut(path string, nlines string) {

	b := []byte(nlines)
	err := ioutil.WriteFile(path, b, 0644)
	if err != nil {
		a.flog.Println(err)
	}
}

// insertShit inserts writes to fname lntxt @ iline
func (a *Aop) insertShit(fname string, iline int, lntxt string) string {
	// insert that shit for front concern
	flines := fileLines(fname)

	out := ""
	for i := 0; i < len(flines); i++ {
		if i == iline {
			out += lntxt + "\n"
		}

		out += flines[i] + "\n"
	}

	a.writeOut(fname, out)

	return out
}
