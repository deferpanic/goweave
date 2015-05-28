package main

import (
	"github.com/deferpanic/goweave/aop"
	"log"
)

const (
	version = "v0.1"
)

// main is the main point of entry for running goweave
func main() {
	log.Println("goweave " + version)

	a := aop.NewAop()
	a.Run()

}
