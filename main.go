package main

import (
	"github.com/deferpanic/goa/aop"
	"log"
)

const (
	version = "v0.1"
)

// main is the main point of entry for running goa
func main() {
	log.Println("goa " + version)

	a := aop.NewAop()
	a.Run()

}
