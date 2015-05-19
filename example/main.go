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

func main() {
	fmt.Println("test of the microphone")

	beforeBob()
	afterSally()
	aroundTom()

	other.FuncMaster()
}
