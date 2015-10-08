package weave

import (
	"os"
	"testing"
)

func TestApplyDeclarationJP(t *testing.T) {
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
			kind: 6,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applyDeclarationJP(os.TempDir()+"/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyDeclarationJP is not transforming correctly")
	}

}
