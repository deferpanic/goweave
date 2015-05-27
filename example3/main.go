package main

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

	time.Sleep(1 * time.Second)
}
