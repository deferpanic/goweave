package aop

import (
	"errors"
	"log"
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

// pointCutKind returns the map id of human expression of pointcut type
func pointCutKind(l string) int {

	matched, err := regexp.MatchString("call(.*)", l)
	if err != nil {
		log.Println(err)
	}

	if matched {
		return 1
	}

	matched, err = regexp.MatchString("execute(.*)", l)
	if err != nil {
		log.Println(err)
	}

	if matched {
		return 2
	}

	return -1
}

// parsePointCut parses out a pointcut from shit
func (a *Aop) parsePointCut(body string) (Pointcut, error) {
	pc := strings.Split(body, "pointcut:")

	if len(pc) > 1 {
		rpc := strings.Split(pc[1], "\n")[0]
		t := strings.TrimSpace(rpc)
		k := pointCutKind(t)

		return Pointcut{
			def:      t,
			funktion: t,
			kind:     k,
		}, nil
	} else {
		return Pointcut{}, errors.New("invalid pointcut")
	}
}

// https://eclipse.org/aspectj/doc/released/progguide/semantics-pointcuts.html#matching

/*
  * explicit method name
    ```go
      "blah"
    ```

  * partial match method name
    ```go
      "b"
    ```

  * function declaration
    ```go
      (w http.ResponseWriter, r *http.Request)
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

// https://eclipse.org/aspectj/doc/released/progguide/semantics-pointcuts.html

// should prob. be a pointcut type
