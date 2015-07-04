package weave

import (
	"testing"
)

func TestApplyWithinJP(t *testing.T) {
	f1 := `package main

import (
	"fmt"
	"time"
)

func slow() {
	time.Sleep(1 * time.Second)
}

func fast() {
}

func everyCall() {
	slow()
	fast()
	slow()
}

func main() {
	everyCall()
}`

	expected := `package main

import (
	"fmt"
	"time"
)

func slow() {
	time.Sleep(1 * time.Second)
}

func fast() {
}

func everyCall() {
ccall("slow", time.Now())
	slow()
ucall("slow", time.Now())
ccall("fast", time.Now())
	fast()
ucall("fast", time.Now())
ccall("slow", time.Now())
	slow()
ucall("slow", time.Now())
}

func main() {
	everyCall()
}
`

	w := &Weave{}

	w.writeOut("/tmp/blah", f1)

	aspect := Aspect{
		advize: Advice{
			before: "ccall(mName, time.Now())",
			after:  "ucall(mName, time.Now())",
		},
		pointkut: Pointcut{
			def:  "everyCall()",
			kind: 3,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applyWithinJP("/tmp/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyWithinJP is not transforming correctly")
	}

}
