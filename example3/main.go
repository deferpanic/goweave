package main

import (
	"fmt"
)

func sally(i int) int {
	return i * 2
}

func bob() int {
	return 1
}

func main() {
	fmt.Println(bob())
}
