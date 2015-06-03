package aop

import (
	"testing"
)

func TestSetDef(t *testing.T) {

	var pkTests = []struct {
		n    string
		kind int
		def  string
	}{
		{"call(beforeBob)", 1, "beforeBob"},
		{"execute(beforeBob)", 2, "beforeBob"},
		{"execute(FuncWithArgs(iarg int, sarg string))", 2, "FuncWithArgs(iarg int, sarg string)"},
		{"execute(FuncWithArgsAndReturn(iarg int, sarg string) (int, error))", 2, "FuncWithArgsAndReturn(iarg int, sarg string) (int, error)"},
	}

	for _, tt := range pkTests {
		k, d, e := setDef(tt.n)
		if k != tt.kind {
			t.Errorf("wrong kind")
		}

		if d != tt.def {
			t.Error(d)
			t.Errorf("wrong def")
		}

		if e != nil {
			t.Error(e)
		}

	}

}
