package main

import (
	"net/http"
)

// panic test
func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("there is no need to panic")
}

func panic2Handler(w http.ResponseWriter, r *http.Request) {
	panic("there is no need to panic")
}

func main() {
	http.HandleFunc("/panic", panicHandler)
	http.HandleFunc("/panic2", panic2Handler)

	http.ListenAndServe(":3000", nil)
}
