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

// pointCutType returns a map of id to human expression of pointcut types
func pointCutType() map[int]string {
	return map[int]string{
		1: "call",
		2: "execution",
	}
}

// set def extracts the joinpoint from a pointcut definition
func setDef(t string) (int, string, error) {

	m := `(execute|call)\((.*)\)`
	re, err := regexp.Compile(m)
	if err != nil {
		return 0, "", errors.New("bad regex")
	}

	res := re.FindAllStringSubmatch(t, -1)
	if len(res[0]) == 3 {
		if res[0][1] == "call" {
			return 1, res[0][2], nil
		} else {
			return 2, res[0][2], nil
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

/*

	::::TODO::::

  * partial match method name
    ```go
      "b"
    ```

  * sub-pkg && method name
    ```go
      pkg/blah
    ```

  * struct && method name
    ```go
      struct.b
    ```

  * sub-pkg && struct && method-name
    ```go
      pkg/struct.b
    ```
*/
