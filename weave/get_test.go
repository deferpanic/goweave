package weave

import (
	"os"
	"testing"
)

func TestApplyGetJP(t *testing.T) {
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
    ch := make(chan int, 2)
    ch <- 1
    ch <- 2
fmt.Println("yo joe")
    fmt.Println(<-ch)
fmt.Println("yo joe")
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
			kind: 4,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applyGetJP(os.TempDir()+"/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyGetJP is not transforming correctly")
	}

}

func TestApplySimpleGetJP(t *testing.T) {
	f1 := `package main

import "fmt"

func main() {
	x := "stuff"
	y := 2
    fmt.Println(x)
    fmt.Println(y)
}
`

	expected := `package main

import "fmt"

func main() {
	x := "stuff"
	y := 2
    fmt.Println(x)
fmt.Println("before get y")
    fmt.Println(y)
}
`

	w := &Weave{}

	w.writeOut(os.TempDir()+"/blah", f1)

	aspect := Aspect{
		advize: Advice{
			before: "fmt.Println(\"before get y\")",
		},
		pointkut: Pointcut{
			def:  "y",
			kind: 4,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applyGetJP(os.TempDir()+"/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyGetJP is not transforming correctly")
	}

}
