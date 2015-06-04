package weave

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

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

// processGoRoutines is in the process of being DEPRECATED
// it only provides regex support for go before/after advice
// once we refactor to AST replacing the go routines this function will
// go away
// this does not write to any files - simply manipulates text
func (w *Weave) processGoRoutines(curfile string, rootpkg string) (string, bool) {
	modified := false

	file, err := os.Open(curfile)
	if err != nil {
		w.flog.Println(err)
	}
	defer file.Close()

	out := ""

	// FIXME
	// poor man's scope
	scope := 0

	cur_aspect := Aspect{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()

		newAspect := pointCutMatch(w.aspects, l)
		if newAspect.pointkut.def != "" {
			modified = true

			scope += 1

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
		w.flog.Println(err)
	}

	return out, modified
}
