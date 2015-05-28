package main

import (
	"fmt"
	"github.com/deferpanic/goweave/example/other"
)

func beforeBob() {
	fmt.Println("before")
}

func afterSally() {
	fmt.Println("after")
}

func aroundTom() {
	fmt.Println("around")
}

func innerFors() {
	fmt.Println("inner")
}

func retstr() string {
	return "string"
}

func retbool() bool {
	if 1 == 1 {
		return true
	} else {
		return false
	}
}

func main() {
	fmt.Println("test of the microphone")

	beforeBob()
	afterSally()
	aroundTom()
	innerFors()

	blah := retstr()
	fmt.Println(blah)
	blahbool := retbool()
	fmt.Println(blahbool)

	other.FuncMaster()
}
