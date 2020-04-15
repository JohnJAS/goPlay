package main

import "fmt"

func main() {
	func(a string){
		fmt.Println(a)
	}("test")
}