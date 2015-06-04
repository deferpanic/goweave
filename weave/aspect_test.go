package aop

import (
	"testing"
)

func TestBefore(t *testing.T) {

	f := `
aspect {
  pointcut: execute(main)
  advice: {
	before: {
    	fmt.Println("before main")
  	}
  }
}
`

	w := &Weave{}
	w.parseAspectFile(f)

	if len(w.aspects) != 1 {
		t.Error("didn't parse aspects")
	}

	first := w.aspects[0]

	if first.pointkut.kind != 2 {
		t.Error("didn't set pointcut definition correctly")
	}

	if first.pointkut.def != "main" {
		t.Error("didn't set pointcut definition correctly")
	}

	if first.pointkut.funktion != "main" {
		t.Error("didn't set pointcut function correctly")
	}

	if len(first.importz) != 0 {
		t.Error("didn't parse imports correctly")
	}

	s := `fmt.Println("before main")`

	if first.advize.before != s {
		t.Error(s)
		t.Error(first.advize.before)
		t.Error("didn't parse advice correctly")
	}

	if first.advize.after != "" {
		t.Error("didn't parse advice correctly")
	}

}

func TestAfter(t *testing.T) {

	f := `
aspect {
  pointcut: execute(main)
  advice: {
	after: {
	 	fmt.Println("after main")
	}
  }
}
`

	w := &Weave{}
	w.parseAspectFile(f)

	if len(w.aspects) != 1 {
		t.Error("didn't parse aspects")
	}

	first := w.aspects[0]

	if first.pointkut.def != "main" {
		t.Error("didn't set pointcut definition correctly")
	}

	if first.pointkut.funktion != "main" {
		t.Error("didn't set pointcut function correctly")
	}

	if len(first.importz) != 0 {
		t.Error("didn't parse imports correctly")
	}

	if first.advize.before != "" {
		t.Error("didn't parse advice correctly")
	}

	if first.advize.after != "fmt.Println(\"after main\")" {
		t.Error("didn't parse advice correctly")
	}

}

func TestParseAspectFile(t *testing.T) {

	f := `
aspect {
  pointcut: execute(main)
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

	w := &Weave{}
	w.parseAspectFile(f)

	if len(w.aspects) != 1 {
		t.Error("didn't parse aspects")
	}

	first := w.aspects[0]

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

func TestNoImports(t *testing.T) {

	f := `
aspect {
  pointcut: execute(main)
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

	w := &Weave{}
	w.parseAspectFile(f)

	if len(w.aspects) != 1 {
		t.Error("didn't parse aspects")
	}

	first := w.aspects[0]

	if first.pointkut.def != "main" {
		t.Error("didn't set pointcut definition correctly")
	}

	if first.pointkut.funktion != "main" {
		t.Error("didn't set pointcut function correctly")
	}

	if len(first.importz) != 0 {
		t.Error("didn't parse imports correctly")
	}

	if first.advize.before != "fmt.Println(\"before main\")" {
		t.Error("didn't parse advice correctly")
	}

	if first.advize.after != "fmt.Println(\"after main\")" {
		t.Error("didn't parse advice correctly")
	}

}

func TestAspectScope(t *testing.T) {

	f := `
aspect {
  pointcut: execute(innerFors)
  advice: {
    before: {
        for i:=0; i<10; i++ {
          fmt.Println(i)
        }
    }
  }
}
`

	w := &Weave{}
	w.parseAspectFile(f)

	if len(w.aspects) != 1 {
		t.Error("didn't parse aspects")
	}

	first := w.aspects[0]

	s := `for i:=0; i<10; i++ {
          fmt.Println(i)
        }`

	if first.advize.before != s {
		t.Error("didn't parse advice correctly")
	}

}
