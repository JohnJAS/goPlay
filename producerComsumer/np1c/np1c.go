/***
multiple producer one consumer
*/
package main

import (
	"fmt"
	"sync"
)

const producerCount int = 4

var messages = [][]string{
	{
		"The world itself's",
		"just one big hoax.",
		"Spamming each other with our",
		"running commentary of bullshit,",
	},
	{
		"but with our things, our property, our money.",
		"I'm not saying anything new.",
		"We all know why we do this,",
		"not because Hunger Games",
		"books make us happy,",
	},
	{
		"masquerading as insight, our social media",
		"faking as intimacy.",
		"Or is it that we voted for this?",
		"Not with our rigged elections,",
	},
	{
		"but because we wanna be sedated.",
		"Because it's painful not to pretend,",
		"because we're cowards.",
		"- Elliot Alderson",
		"Mr. Robot",
	},
}

func producer(pipeline chan <-string,i int, wg *sync.WaitGroup){
	for _, str := range messages[i] {
		pipeline <- str
	}
	wg.Done()
}


func consumer(pipeline <-chan string, done chan<- struct{}){
	for{
		select{
		case str,ok := <- pipeline:
			if ok {
				fmt.Println(str)
			}else{
				done<-struct {}{}
				break
			}
		default:

		}
	}
}

func main(){
	pipeline := make(chan string)
	done := make(chan struct{})

	wg := sync.WaitGroup{}
	for i:=0;i<producerCount;i++{
		wg.Add(1)
		go producer(pipeline,i, &wg)

	}
	go consumer(pipeline,done)

	wg.Wait()
	close(pipeline)



	<- done
}