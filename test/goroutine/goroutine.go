package main

import (
	"fmt"
	"math/rand"
)

func Generator(signal chan int) chan int {
	result := make(chan int)
	go func() {
		Lable:
		for {
			select {
			case result <- rand.Int():
			case <-signal:
				close(result)
				break Lable
			}
		}
	}()

	return result
}

func main() {
	signal := make(chan int)

	c := Generator(signal)

	fmt.Println(<-c)
	fmt.Println(<-c)

	//signal <- 0
	close(signal)

	fmt.Println(<-signal)


	fmt.Println(<-c)
	fmt.Println(<-c)
}
