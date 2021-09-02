/***
one producer one consumer
 */

package main

import "fmt"

var messages = []string{
	"The world itself's",
	"just one big hoax.",
	"Spamming each other with our",
	"running commentary of bullshit,",
	"masquerading as insight, our social media",
	"faking as intimacy.",
	"Or is it that we voted for this?",
	"Not with our rigged elections,",
	"but with our things, our property, our money.",
	"I'm not saying anything new.",
	"We all know why we do this,",
	"not because Hunger Games",
	"books make us happy,",
	"but because we wanna be sedated.",
	"Because it's painful not to pretend,",
	"because we're cowards.",
	"- Elliot Alderson",
	"Mr. Robot",
}

func producer(pipeline chan <-string){
	for _, str := range messages {
		pipeline <- str
	}
	close(pipeline)
}

func consumer(pipeline <-chan string, done chan<- struct{}){
	for {
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

	go producer(pipeline)
	go consumer(pipeline,done)
	<- done
}