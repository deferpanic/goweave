package main

import (
	"github.com/deferpanic/gocut/aop"
	"log"
)

const (
	version = "v0.1"
)

// main is the main point of entry for running gocut
func main() {
	log.Println("gocut " + version)

	a := aop.NewAop()
	a.Run()

}
