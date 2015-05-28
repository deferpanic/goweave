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
	http.HandleFunc("/panic", dps.HTTPHandlerFunc(panicHandler))
	http.HandleFunc("/panic2", dps.HTTPHandlerFunc(panic2Handler))
}
`
	aop := &Aop{}

	aop.writeOut("/tmp/blah", f1)

	aspect := Aspect{
		advize: Advice{
			around: "http.HandleFunc(d, dps.HTTPHandlerFunc(s))",
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

func TestGoTx(t *testing.T) {
	a := &Aop{}

	f1 := `package main

import (
	"fmt"
	"time"
)

func stuff() {
	panic("panic")
}

func blah() {
	stuff()
	fmt.Println("never get here")
}

func main() {
	go blah()

	go func() {
		fmt.Println("inline")
		blah()
	}()

	time.Sleep(1 * time.Second)
}`

	a.writeOut("/tmp/blah_test_go", f1)

	aspect1 := Aspect{
		importz: []string{
			"fmt",
			"github.com/deferpanic/deferclient/deferstats",
		},
		advize: Advice{
			before: "dps := deferstats.NewClient(\"v00L0K6CdKjE4QwX5DL1iiODxovAHUfo\")\ngo dps.CaptureStats()",
		},
		pointkut: Pointcut{
			def: "main",
		},
	}

	aspect2 := Aspect{
		advize: Advice{
			before: "defer dps.Persist()\nfmt.Println(\"there is no need to panic\")",
		},
		pointkut: Pointcut{
			def: "go",
		},
	}

	aspects := []Aspect{}

	aspects = append(aspects, aspect1)
	aspects = append(aspects, aspect2)

	a.aspects = aspects

	rootpkg := a.rootPkg()

	after, i := a.txAspects("/tmp/blah_test_go", rootpkg)

	expected :=
		`package main

import (
	"fmt"
	"time"
)

func stuff() {
	panic("panic")

}

func blah() {
	stuff()
	fmt.Println("never get here")

}

func main() {
dps := deferstats.NewClient("v00L0K6CdKjE4QwX5DL1iiODxovAHUfo")
go dps.CaptureStats()
	go func(){
defer dps.Persist()
fmt.Println("there is no need to panic")
blah()
}()

go func(){
defer dps.Persist()
fmt.Println("there is no need to panic")
		fmt.Println("inline")
		blah()

}()

	time.Sleep(1 * time.Second)

}
`
	if after != expected {
		t.Error("\n" + "#" + after + "#")
		t.Error("\n" + "#" + expected + "#")
		t.Error("txAspects is not transforming correctly")
	}

	if len(i) != 2 {
		t.Error("txAspects is not parsing imports correctly")
	}

}
