package aop

// advice has a function to wrap advice around and code for said
// function
type advice struct {
	funktion     string
	before       string
	after        string
	adviceTypeId int
}
