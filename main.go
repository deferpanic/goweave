package main

import (
	"github.com/deferpanic/goweave/weave"
	"log"
)

const (
	version = "v0.1"
)

// main is the main point of entry for running goweave
func main() {
	log.Println("goweave " + version)

	w := weave.NewWeave()
	w.Run()

}
