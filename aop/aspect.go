package aop

import (
	"errors"
	"io/ioutil"
	"strings"
)

// aspect contains advice, pointcuts and any imports needed
type Aspect struct {
	advize   Advice
	pointkut Pointcut
	importz  []string
}

// grab_aspects looks for an aspect file for each file
// this seems lame and contrary to what we want...
// I'd go as far as to say that we want this to be cross-pkg w/in root?
//
// maybe the rule should be - aspects are valid for anything in a
// project root?
func (a *Aop) loadAspects() {

	fz := a.findAspects()
	for i := 0; i < len(fz); i++ {

		buf, err := ioutil.ReadFile(fz[i])
		if err != nil {
			a.flog.Println(err)
		}
		s := string(buf)

		a.parseAspectFile(s)
	}
}

// parsePointCut parses out a pointcut from shit
func (a *Aop) parsePointCut(body string) (Pointcut, error) {
	pc := strings.Split(body, "pointcut:")

	if len(pc) > 1 {
		rpc := strings.Split(pc[1], "\n")[0]
		t := strings.TrimSpace(rpc)

		return Pointcut{
			def:      t,
			funktion: t,
		}, nil
	} else {
		return Pointcut{}, errors.New("invalid pointcut")
	}
}

// parseImports returns an array of imports for the corresponding advice
func (a *Aop) parseImports(body string) []string {
	impbrace := strings.Split(body, "imports (")

	if len(impbrace) > 1 {
		end := strings.Split(impbrace[1], ")")[0]
		t := strings.TrimSpace(end)
		return strings.Split(t, "\n")
	} else {
		return []string{}
	}
}

// containsBefore returns true if the body has before advice
func (a *Aop) containsBefore(body string) bool {
	if strings.Contains(body, "before: {") {
		return true
	} else {
		return false
	}
}

// containsAfter returns true if the body has after advice
func (a *Aop) containsAfter(body string) bool {
	if strings.Contains(body, "after: {") {
		return true
	} else {
		return false
	}
}

// rightBraceCnt returns the number of right braces in a string
func (a *Aop) rightBraceCnt(body string) int {
	return strings.Count(body, "}")
}

// parseAdvice returns advice about this aspect
func (a *Aop) parseAdvice(body string) Advice {
	advize := strings.Split(body, "advice:")[1]

	a4t := ""
	b4t := ""

	bbrace := strings.Split(advize, "before: {")
	if len(bbrace) > 1 {
		// fixme
		if a.containsAfter(bbrace[1]) {
			b4 := strings.Split(bbrace[1], "}")[0]
			b4t = strings.TrimSpace(b4)
			// ...
		} else {
			cnt := a.rightBraceCnt(bbrace[1])
			// have at most 3 right braces
			// 3 - 3 = 0
			// 4 - 3 = 1
			b4 := strings.SplitAfter(bbrace[1], "}")
			rb := ""
			if cnt == 3 {
				rb = b4[0]
				rb = rb[:len(rb)-1]
			} else {
				for i := 0; i < cnt-3; i++ {
					rb += strings.TrimSpace(b4[i])
				}
			}
			b4t = strings.TrimSpace(rb)
		}
	}

	abrace := strings.Split(advize, "after: {")
	if len(abrace) > 1 {
		a4 := strings.Split(abrace[1], "}")[0]
		a4t = strings.TrimSpace(a4)
	}

	return Advice{
		before: b4t,
		after:  a4t,
	}

}

// parseAspectFile loads an individual file containing aspects
func (a *Aop) parseAspectFile(body string) {
	results := []Aspect{}

	aspects := strings.Split(body, "aspect {")

	for i := 1; i < len(aspects); i++ {

		aspect := aspects[i]
		azpect := Aspect{}

		pk, err := a.parsePointCut(aspect)
		if err != nil {
			a.flog.Println(err)
			continue
		} else {
			azpect.pointkut = pk
		}

		azpect.importz = a.parseImports(aspect)

		azpect.advize = a.parseAdvice(aspect)

		results = append(results, azpect)

	}

	a.aspects = results

}
