package weave

import (
	"go/ast"

	"testing"
)

func TestApplyAroundAdvice(t *testing.T) {

	f1 := `package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/panic", panicHandler)
	http.HandleFunc("/panic2", panic2Handler)
}
`

	expected := `package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/panic", dps.HTTPHandlerFunc(panicHandler))
	http.HandleFunc("/panic2", dps.HTTPHandlerFunc(panic2Handler))
}
`
	w := NewWeave()

	w.writeOut("/tmp/blah", f1)

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
	w.aspects = aspects

	after := w.applyAroundAdvice("/tmp/blah")

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyAroundAdvice is not transforming correctly")
	}

}

func TestGoTx(t *testing.T) {
	w := &Weave{}

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

	w.writeOut("/tmp/blah_test_go", f1)

	aspect2 := Aspect{
		advize: Advice{
			before: "fmt.Println(\"there is no need to panic\")",
		},
		pointkut: Pointcut{
			def: "go",
		},
	}

	aspects := []Aspect{}

	aspects = append(aspects, aspect2)

	w.aspects = aspects

	rootpkg := w.rootPkg()

	after, _ := w.processGoRoutines("/tmp/blah_test_go", rootpkg)

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
	go func(){
fmt.Println("there is no need to panic")
blah()
}()

go func(){
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
		t.Error("processGoRoutines is not transforming correctly")
	}

}

func TestApplyExecutionJP(t *testing.T) {

	f1 := `package main

import (
	"fmt"
	"net/http"
)

func hiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi")
}

func hi2Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi2")
}

func main() {
	http.HandleFunc("/hi", hiHandler)
	http.HandleFunc("/hi2", hi2Handler)

	http.ListenAndServe(":3000", nil)
}`

	expected := `package main

import (
	"fmt"
	"net/http"
)

func hiHandler(w http.ResponseWriter, r *http.Request) {
fmt.Println("before call")
	fmt.Fprintf(w, "Hi")
fmt.Println("after call")
}

func hi2Handler(w http.ResponseWriter, r *http.Request) {
fmt.Println("before call")
	fmt.Fprintf(w, "Hi2")
fmt.Println("after call")
}

func main() {
	http.HandleFunc("/hi", hiHandler)
	http.HandleFunc("/hi2", hi2Handler)

	http.ListenAndServe(":3000", nil)
}
`

	w := &Weave{}

	w.writeOut("/tmp/blah", f1)

	aspect := Aspect{
		advize: Advice{
			before: "fmt.Println(\"before call\")",
			after:  "fmt.Println(\"after call\")",
		},
		pointkut: Pointcut{
			def:  "(http.ResponseWriter, *http.Request)",
			kind: 2,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applyExecutionJP("/tmp/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyExecutionJP is not transforming correctly")
	}

}

func TestContainArgs(t *testing.T) {

	fields := []*ast.Field{}

	if containArgs("main()", fields) != true {
		t.Error("picking up args where there should not be")
	}
}

func TestApplyExecutionJPMain(t *testing.T) {
	f1 := `package main

import (
	"fmt"
)

func main() {
	fmt.Println("test of the microphone")
}`

	expected := `package main

import (
	"fmt"
)

func main() {
fmt.Println("before main")
	fmt.Println("test of the microphone")
}
`

	w := &Weave{}

	w.writeOut("/tmp/blah", f1)

	aspect := Aspect{
		advize: Advice{
			before: "fmt.Println(\"before main\")",
		},
		pointkut: Pointcut{
			def:  "main()",
			kind: 2,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applyExecutionJP("/tmp/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyExecutionJP is not transforming correctly")
	}

}

func TestApplyExecutionJPBefore(t *testing.T) {
	f1 := `package main

import (
	"fmt"
)

func beforeBob() {
	fmt.Println("bob")
}

func main() {
	beforeBob()
}`

	expected := `package main

import (
	"fmt"
)

func beforeBob() {
fmt.Println("before bob")
	fmt.Println("bob")
}

func main() {
	beforeBob()
}
`

	w := &Weave{}

	w.writeOut("/tmp/blah", f1)

	aspect := Aspect{
		advize: Advice{
			before: "fmt.Println(\"before bob\")",
		},
		pointkut: Pointcut{
			def:  "beforeBob()",
			kind: 2,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applyExecutionJP("/tmp/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyExecutionJP is not transforming correctly")
	}

}

func TestApplyExecutionJPAfter(t *testing.T) {
	f1 := `package main

import (
	"fmt"
)

func afterAnny() {
	fmt.Println("anny")
}

func main() {
	afterAnny()
}`

	expected := `package main

import (
	"fmt"
)

func afterAnny() {
	fmt.Println("anny")
fmt.Println("after anny")
}

func main() {
	afterAnny()
}
`

	w := &Weave{}

	w.writeOut("/tmp/blah", f1)

	aspect := Aspect{
		advize: Advice{
			after: "fmt.Println(\"after anny\")",
		},
		pointkut: Pointcut{
			def:  "afterAnny()",
			kind: 2,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applyExecutionJP("/tmp/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyExecutionJP is not transforming correctly")
	}

}

func TestApplyExecutionJPAfterwReturn(t *testing.T) {
	f1 := `package main

import (
	"fmt"
)

func afterwReturn() int {
	fmt.Println("anny")
	return 1
}

func main() {
	afterAnny()
}`

	expected := `package main

import (
	"fmt"
)

func afterwReturn() int {
	fmt.Println("anny")
fmt.Println("after anny")
	return 1
}

func main() {
	afterAnny()
}
`

	w := &Weave{}

	w.writeOut("/tmp/blah", f1)

	aspect := Aspect{
		advize: Advice{
			after: "fmt.Println(\"after anny\")",
		},
		pointkut: Pointcut{
			def:  "afterwReturn()",
			kind: 2,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applyExecutionJP("/tmp/blah", f1)

	if after != expected {
		t.Error("##" + after + "##")
		t.Error("##" + expected + "##")
		t.Error("applyExecutionJP is not transforming correctly")
	}

}

func TestApplyExecutionJPAround(t *testing.T) {
	f1 := `package main

import (
	"fmt"
)

func aroundArnie() {
	fmt.Println("arnie")
}

func main() {
	aroundArnie()
}`

	expected := `package main

import (
	"fmt"
)

func aroundArnie() {
fmt.Println("before arnie")
	fmt.Println("arnie")
fmt.Println("after arnie")
}

func main() {
	aroundArnie()
}
`

	w := &Weave{}

	w.writeOut("/tmp/blah", f1)

	aspect := Aspect{
		advize: Advice{
			before: "fmt.Println(\"before arnie\")",
			after:  "fmt.Println(\"after arnie\")",
		},
		pointkut: Pointcut{
			def:  "aroundArnie()",
			kind: 2,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applyExecutionJP("/tmp/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyExecutionJP is not transforming correctly")
	}

}

func TestApplyExecutionJPInnerFors(t *testing.T) {
	f1 := `package main

import (
	"fmt"
)

func innerFors() {
	fmt.Println("inner")
}

func main() {
	innerFors()
}`

	expected := `package main

import (
	"fmt"
)

func innerFors() {
for i:=0; i<10; i++ {
	fmt.Println(i)
}
	fmt.Println("inner")
}

func main() {
	innerFors()
}
`

	w := &Weave{}

	w.writeOut("/tmp/blah", f1)

	aspect := Aspect{
		advize: Advice{
			before: "for i:=0; i<10; i++ {\n\tfmt.Println(i)\n}",
		},
		pointkut: Pointcut{
			def:  "innerFors()",
			kind: 2,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applyExecutionJP("/tmp/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyExecutionJP is not transforming correctly")
	}

}

func TestApplyExecutionJPRetStr(t *testing.T) {
	f1 := `package main

import (
	"fmt"
)

func retStr() {
	fmt.Println("string")
}

func main() {
	retStr()
}`

	expected := `package main

import (
	"fmt"
)

func retStr() {
fmt.Println("before ret")
	fmt.Println("string")
}

func main() {
	retStr()
}
`

	w := &Weave{}

	w.writeOut("/tmp/blah", f1)

	aspect := Aspect{
		advize: Advice{
			before: "fmt.Println(\"before ret\")",
		},
		pointkut: Pointcut{
			def:  "retStr()",
			kind: 2,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applyExecutionJP("/tmp/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyExecutionJP is not transforming correctly")
	}

}

func TestApplyExecutionJPRetBool(t *testing.T) {
	f1 := `package main

import (
	"fmt"
)

func retBool() bool {
	if 1 == 1 {
		return true
	} else {
		return false
	}
}

func main() {
	retBool()
}`

	expected := `package main

import (
	"fmt"
)

func retBool() bool {
fmt.Println("before bool")
	if 1 == 1 {
		return true
	} else {
		return false
	}
}

func main() {
	retBool()
}
`

	w := &Weave{}

	w.writeOut("/tmp/blah", f1)

	aspect := Aspect{
		advize: Advice{
			before: "fmt.Println(\"before bool\")",
		},
		pointkut: Pointcut{
			def:  "retBool()",
			kind: 2,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applyExecutionJP("/tmp/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyExecutionJP is not transforming correctly")
	}

}

func TestApplyCallAdvice(t *testing.T) {
	f1 := `package main

import (
	"fmt"
)

func getStuff(i int) {
	fmt.Println(i)
}

func main() {
	for i:=0; i<10; i++ {
		getStuff(i)
	}
}`

	expected := `package main

import (
	"fmt"
)

func getStuff(i int) {
	fmt.Println(i)
}

func main() {
	for i:=0; i<10; i++ {
fmt.Println("before call")
		getStuff(i)
	}
}
`

	w := &Weave{}

	w.writeOut("/tmp/blah", f1)

	aspect := Aspect{
		advize: Advice{
			before: "fmt.Println(\"before call\")",
		},
		pointkut: Pointcut{
			def:  "getStuff(int i)",
			kind: 1,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applyCallAdvice("/tmp/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyCallAdvice is not transforming correctly")
	}

}

func TestApplyCallAdvicePkgDot(t *testing.T) {

	f1 := `package main

import (
	"fmt"
	"strconv"
)

func blah() {
	fmt.Println(strconv.Itoa(2))
}

func main() {
	fmt.Println(strconv.Itoa(44))

	blah()
}`

	expected := `package main

import (
	"fmt"
	"strconv"
)

func blah() {
fmt.Println("strconv called")
	fmt.Println(strconv.Itoa(2))
}

func main() {
fmt.Println("strconv called")
	fmt.Println(strconv.Itoa(44))

	blah()
}
`

	w := &Weave{}

	w.writeOut("/tmp/blah", f1)

	aspect := Aspect{
		advize: Advice{
			before: "fmt.Println(\"strconv called\")",
		},
		pointkut: Pointcut{
			def:  "strconv.Itoa(int i)",
			kind: 1,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	w.aspects = aspects

	after := w.applyCallAdvice("/tmp/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyCallAdvice is not transforming correctly")
	}

}

func TestPkgRewrite(t *testing.T) {
	f1 := `package main

import (
	"github.com/some/stuff/subpkg"
)

func main() {
}`

	s := "github.com/some/stuff"

	w := NewWeave()
	w.writeOut("/tmp/blah", f1)

	af := w.ParseAST("/tmp/blah")

	pruned := w.pruneImports(af, s)

	if pruned[0].Path.Value != "\"_weave/github.com/some/stuff/subpkg\"" {
		t.Error(pruned[0].Path.Value)
		t.Error("pruneImports not working")
	}

}

func TestApplyGlobalAdvice(t *testing.T) {

	f1 := `package main

import (
	"fmt"
)

func getStuff(i int) {
	fmt.Println(i)
}

func main() {

	for i := 0; i < 10; i++ {
		getStuff(i)
	}
}`

	expected := `package main

import (
	"fmt"
)
var myCnt int =0

func getStuff(i int) {
myCnt += 1
	fmt.Println(i)
}

func main() {

	for i := 0; i < 10; i++ {
		getStuff(i)
	}
fmt.Println(myCnt)
}
`

	w := NewWeave()

	w.writeOut("/tmp/blah", f1)

	aspect := Aspect{
		advize: Advice{
			before: "myCnt += 1",
		},
		pointkut: Pointcut{
			def:  "getStuff(int i)",
			kind: 2,
		},
	}

	aspect2 := Aspect{
		advize: Advice{
			before: "var myCnt int =0",
		},
		pointkut: Pointcut{
			def:  "*",
			kind: 2,
		},
	}

	aspect3 := Aspect{
		advize: Advice{
			after: "fmt.Println(myCnt)",
		},
		pointkut: Pointcut{
			def:  "main()",
			kind: 2,
		},
	}

	aspects := []Aspect{}
	aspects = append(aspects, aspect)
	aspects = append(aspects, aspect2)
	aspects = append(aspects, aspect3)

	w.aspects = aspects

	fp := "/tmp/blah"
	stuff := w.applyGlobalAdvice(fp, f1)
	w.writeOut(fp, stuff)

	stuff = w.applyCallAdvice(fp, stuff)
	w.writeOut(fp, stuff)

	after := w.applyExecutionJP(fp, stuff)
	w.writeOut(fp, stuff)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyGlobalAdvice is not transforming correctly")
	}

}

/*
func TestParsey(t *testing.T) {
	src := `for i:=0; i< 10; i++ {
			fmt.Println(i)
			}`

	fname := "/tmp/blah"

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fname, src, parser.Mode(0))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(file)
}
*/
