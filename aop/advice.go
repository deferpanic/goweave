package aop

import (
	"strings"
)

// Advice has a function to wrap advice around and code for said
// function
type Advice struct {
	before string
	after  string
	around string

	adviceTypeId int
}

// adviceType returns a map of id to human expression of advice types
func adviceType() map[int]string {
	return map[int]string{
		1: "before",
		2: "after",
		3: "around",
	}
}

// adviceKind returns the map id of human expression of advice type
func adviceKind(l string) int {
	stuff := strings.Split(l, ": ")
	ostuff := strings.Split(stuff[1], " ")

	switch ostuff[0] {
	case "before":
		return 1
	case "after":
		return 2
	case "around":
		return 3
	}

	return -1
}
