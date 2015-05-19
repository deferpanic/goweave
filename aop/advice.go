package aop

// Advice has a function to wrap advice around and code for said
// function
type Advice struct {
	before string
	after  string

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
