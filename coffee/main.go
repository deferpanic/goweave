package main

import (
	"fmt"
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
}
