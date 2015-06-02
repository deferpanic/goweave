package aop

import (
	"go/ast"
	"testing"
)

func TestApplyAroundAdvice(t *testing.T) {

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

	after := aop.applyAroundAdvice("/tmp/blah", f1)

	if after != expected {
		t.Error("applyAroundAdvice is not transforming correctly")
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

	aspect2 := Aspect{
		advize: Advice{
			before: "defer dps.Persist()\nfmt.Println(\"there is no need to panic\")",
		},
		pointkut: Pointcut{
			def: "go",
		},
	}

	aspects := []Aspect{}

	aspects = append(aspects, aspect2)

	a.aspects = aspects

	rootpkg := a.rootPkg()

	after, _ := a.txAspects("/tmp/blah_test_go", rootpkg)

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

}

/*

FIXME FIXME FIXME

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
*/

func TestApplyExecutionJP(t *testing.T) {

	f1 := `package main

import (
	"net/http"
)

// panic test
func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("there is no need to panic")
}

func panic2Handler(w http.ResponseWriter, r *http.Request) {
	panic("there is no need to panic")
}

func main() {
	http.HandleFunc("/panic", panicHandler)
	http.HandleFunc("/panic2", panic2Handler)

	http.ListenAndServe(":3000", nil)
}`

	expected := `package main

import ("net/http"
)

// panic test
func panicHandler(w http.ResponseWriter, r *http.Request) {
fmt.Println("before call")
	panic("there is no need to panic")
fmt.Println("after call")
}

func panic2Handler(w http.ResponseWriter, r *http.Request) {
fmt.Println("before call")
	panic("there is no need to panic")
fmt.Println("after call")
}

func main() {
	http.HandleFunc("/panic", panicHandler)
	http.HandleFunc("/panic2", panic2Handler)

	http.ListenAndServe(":3000", nil)
}
`

	aop := &Aop{}

	aop.writeOut("/tmp/blah", f1)

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
	aop.aspects = aspects

	after := aop.applyExecutionJP("/tmp/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyExecutionJP is not transforming correctly")
	}

}

func TestContainArgs(t *testing.T) {

	fields := []*ast.Field{}

	if containArgs("main()", fields) != true {
		t.Error("fuck")
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

import ("fmt"
)

func main() {
fmt.Println("before main")
	fmt.Println("test of the microphone")
}
`

	aop := &Aop{}

	aop.writeOut("/tmp/blah", f1)

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
	aop.aspects = aspects

	after := aop.applyExecutionJP("/tmp/blah", f1)

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

import ("fmt"
)

func beforeBob() {
fmt.Println("before bob")
	fmt.Println("bob")
}

func main() {
	beforeBob()
}
`

	aop := &Aop{}

	aop.writeOut("/tmp/blah", f1)

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
	aop.aspects = aspects

	after := aop.applyExecutionJP("/tmp/blah", f1)

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

import ("fmt"
)

func afterAnny() {
	fmt.Println("anny")
fmt.Println("after anny")
}

func main() {
	afterAnny()
}
`

	aop := &Aop{}

	aop.writeOut("/tmp/blah", f1)

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
	aop.aspects = aspects

	after := aop.applyExecutionJP("/tmp/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
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

import ("fmt"
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

	aop := &Aop{}

	aop.writeOut("/tmp/blah", f1)

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
	aop.aspects = aspects

	after := aop.applyExecutionJP("/tmp/blah", f1)

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

import ("fmt"
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

	aop := &Aop{}

	aop.writeOut("/tmp/blah", f1)

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
	aop.aspects = aspects

	after := aop.applyExecutionJP("/tmp/blah", f1)

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

import ("fmt"
)

func retStr() {
fmt.Println("before ret")
	fmt.Println("string")
}

func main() {
	retStr()
}
`

	aop := &Aop{}

	aop.writeOut("/tmp/blah", f1)

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
	aop.aspects = aspects

	after := aop.applyExecutionJP("/tmp/blah", f1)

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

import ("fmt"
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

	aop := &Aop{}

	aop.writeOut("/tmp/blah", f1)

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
	aop.aspects = aspects

	after := aop.applyExecutionJP("/tmp/blah", f1)

	if after != expected {
		t.Error(after)
		t.Error(expected)
		t.Error("applyExecutionJP is not transforming correctly")
	}

}
