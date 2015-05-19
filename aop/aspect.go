package aop

import (
	"io/ioutil"
	"log"
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
func (a *Aop) parsePointCut(body string) Pointcut {
	pc := strings.Split(body, "pointcut:")[1]
	rpc := strings.Split(pc, "\n")[0]
	t := strings.TrimSpace(rpc)

	return Pointcut{
		def:      t,
		funktion: t,
	}

}

// parseImports returns an array of imports for the corresponding advice
func (a *Aop) parseImports(body string) []string {
	imp := strings.Split(body, "imports (")[1]
	end := strings.Split(imp, ")")[0]
	t := strings.TrimSpace(end)
	return strings.Split(t, "\n")
}

// parseAdvice returns advice about this aspect
func (a *Aop) parseAdvice(body string) Advice {
	advize := strings.Split(body, "advice:")[1]

	be := strings.Split(advize, "before: {")[1]
	b4 := strings.Split(be, "}")[0]
	b4t := strings.TrimSpace(b4)

	ae := strings.Split(advize, "after: {")[1]
	a4 := strings.Split(ae, "}")[0]
	a4t := strings.TrimSpace(a4)

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

		azpect.pointkut = a.parsePointCut(aspect)
		azpect.importz = a.parseImports(aspect)

		azpect.advize = a.parseAdvice(aspect)

		results = append(results, azpect)

	}

	a.aspects = results

	log.Println(len(a.aspects))
	/*
				if strings.Contains(l, "advice") {
					blah := strings.Split(l, "execution(\"")[1]
					shiz := strings.Split(blah, "\"")[0]

					aspect := aspect{}
					avize := advice{}
					aspect.avice = avize
					avize.funktion = shiz

					avize.adviceTypeId = adviceKind(l)

					cur_advice = avize
				} else if strings.Contains(l, "}") {
					results = append(results, cur_advice)
				} else if strings.Contains(l, "goaProceed()") {
					donebegin = true
					continue
				} else {

					// before
					if cur_advice.adviceTypeId == 1 {
						cur_advice.before += l + "\n"
					} else if cur_advice.adviceTypeId == 2 {
						cur_advice.after += l + "\n"
					} else {
						if donebegin {
							cur_advice.after += l + "\n"
						} else {
							cur_advice.before += l + "\n"
						}
					}

				}

			}

		}
	*/

}
