package aop

// Advice has a function to wrap advice around and code for said
// function
type Advice struct {
	before string
	after  string

	adviceTypeId int
}
