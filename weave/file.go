package aop

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
)

// reWriteFile rewrites curfile with out && adds any missing imports
func (w *Weave) reWriteFile(curfile string, out string, importsNeeded []string) {

	f, err := os.Create(w.tmpLocation() + "/" + curfile)
	if err != nil {
		w.flog.Println(err)
	}

	defer f.Close()

	out = w.addMissingImports(importsNeeded, out)

	_, err = f.WriteString(out)
	if err != nil {
		w.flog.Println(err)
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
func (w *Weave) writeOut(path string, nlines string) {

	b := []byte(nlines)
	err := ioutil.WriteFile(path, b, 0644)
	if err != nil {
		w.flog.Println(err)
	}
}

// writeAtLine inserts writes to fname lntxt @ iline
func (w *Weave) writeAtLine(fname string, iline int, lntxt string) string {
	flines := fileLines(fname)

	out := ""
	for i := 0; i < len(flines); i++ {
		if i == iline {
			out += lntxt + "\n"
		}

		out += flines[i] + "\n"
	}

	w.writeOut(fname, out)

	return out
}
