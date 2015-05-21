package main

import (
	"fmt"
	"github.com/deferpanic/goa/example/other"
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

func main() {
	fmt.Println("test of the microphone")

	beforeBob()
	afterSally()
	aroundTom()
	innerFors()

	other.FuncMaster()
}
