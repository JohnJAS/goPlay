package main

import (
	"fmt"
)

func main() {
	a := [...]int{1, 2, 3, 4, 5}

	b := []int{1, 2, 3, 4, 5}

	c := make([]int, 5, 10)

	fmt.Println(len(a), cap(a))
	fmt.Println(len(b), cap(b))
	fmt.Println(len(c), cap(c))
}
