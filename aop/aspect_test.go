package aop

import (
	"testing"
)

func TestParseAspectFile(t *testing.T) {

	f := `
aspect {
  pointcut: main
  imports (
    "fmt"
    "github.com/deferpanic/deferclient"
  )
  advice: {
	before: {
    	fmt.Println("before main")
  	}
	after: {
	 	fmt.Println("after main")
	}
  }
}
`

	aop := &Aop{}
	aop.parseAspectFile(f)

	if len(aop.aspects) != 1 {
		t.Error("didn't parse aspects")
	}

	first := aop.aspects[0]

	if first.pointkut.def != "main" {
		t.Error("didn't set pointcut definition correctly")
	}

	if first.pointkut.funktion != "main" {
		t.Error("didn't set pointcut function correctly")
	}

	if len(first.importz) != 2 {
		t.Error("didn't parse imports correctly")
	}

	if first.advize.before != "fmt.Println(\"before main\")" {
		t.Error("didn't parse advice correctly")
	}

	if first.advize.after != "fmt.Println(\"after main\")" {
		t.Error("didn't parse advice correctly")
	}

}
