package aop

import (
	"testing"
)

func TestPointCutKind(t *testing.T) {

	pc := "pointcut: call(bob)"

	k := pointCutKind(pc)

	if k != 1 {
		t.Error("didn't parse pointcut appropriately")
	}

	pc = "pointcut: execute(bob)"

	k = pointCutKind(pc)

	if k != 2 {
		t.Error("didn't parse pointcut appropriately")
	}

}
