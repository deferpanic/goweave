package weave

import (
	"errors"
	"regexp"
	"strings"
)

// Pointcut describe how to apply a particular aspect
type Pointcut struct {
	def      string
	pkg      string
	funktion string
	kind     int
}

// set def extracts the joinpoint from a pointcut definition
func setDef(t string) (int, string, error) {

	m := `(execute|call|within|get|set|declaration)\((.*)\)`
	re, err := regexp.Compile(m)
	if err != nil {
		return 0, "", errors.New("bad regex")
	}

	res := re.FindAllStringSubmatch(t, -1)
	if len(res[0]) == 3 {
		if res[0][1] == "call" {
			return 1, res[0][2], nil
		} else if res[0][1] == "execute" {
			return 2, res[0][2], nil
		} else if res[0][1] == "within" {
			return 3, res[0][2], nil
		} else if res[0][1] == "get" {
			return 4, res[0][2], nil
		} else if res[0][1] == "set" {
			return 5, res[0][2], nil
		} else if res[0][1] == "declaration" {
			return 6, res[0][2], nil
		} else {
			return 0, "", errors.New("bad pointcut")
		}
	} else {
		return 0, "", errors.New("bad pointcut")
	}
}

// parsePointCut parses out a pointcut from an aspect
func (w *Weave) parsePointCut(body string) (Pointcut, error) {
	pc := strings.Split(body, "pointcut:")

	if len(pc) > 1 {
		rpc := strings.Split(pc[1], "\n")[0]
		t := strings.TrimSpace(rpc)

		k, def, err := setDef(t)
		if err != nil {
			return Pointcut{}, err
		}

		return Pointcut{
			def:      def,
			funktion: def,
			kind:     k,
		}, nil
	} else {
		return Pointcut{}, errors.New("invalid pointcut" + body)
	}
}
