package weave

import (
	"os"
	"testing"
)

func TestApplySetJP(t *testing.T) {
	f1 := `package main

import "fmt"

func main() {
    ch := make(chan int, 2)
    ch <- 1
    ch <- 2
    fmt.Println(<-ch)
    fmt.Println(<-ch)
}
`

	expected := `package main

import "fmt"

func main() {
fmt.Println("yo joe")
    ch := make(chan int, 2)
    ch <- 1
    ch <- 2
    fmt.Println(<-ch)
    fmt.Println(<-ch)
}
`

	w := &Weave{}

	w.writeOut(os.TempDir()+"/blah", f1)

	aspect := Aspect{
		advize: Advice{
			before: "fmt.Println(\"yo joe\")",
		},
		pointkut: Pointcut{
			def:  "ch",
			kind: 5,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applySetJP(os.TempDir()+"/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applySetJP is not transforming correctly")
	}

}
