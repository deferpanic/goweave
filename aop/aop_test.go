package aop

import (
	"testing"
)

func TestTxAfter(t *testing.T) {

	f1 := `package main

func main() {
	http.HandleFunc("/panic", panicHandler)
	http.HandleFunc("/panic2", panic2Handler)
}
`

	expected := `package main

func main() {
	http.HandleFunc("/panic", dps.HTTPHandler(panicHandler))
	http.HandleFunc("/panic2", dps.HTTPHandler(panic2Handler))
}
`
	aop := &Aop{}

	aop.writeImports("/tmp/blah", f1)

	aspect := Aspect{
		advize: Advice{
			around: "http.HandleFunc(d, dps.HTTPHandler(s))",
		},
		pointkut: Pointcut{
			def: "http.HandleFunc(d, s)",
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	aop.aspects = aspects

	after := aop.txAfter("/tmp/blah", f1)

	if after != expected {
		t.Error("txAfter is not transforming correctly")
	}

}
