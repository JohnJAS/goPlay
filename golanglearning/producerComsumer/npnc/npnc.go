/***
multiple producer one consumer
*/
package main

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

const producerCount int = 3
const consumerCount int = 3

var workers []*producers

type producers struct {
	myQ  chan string
	quit chan bool
	id   int
}

func comsumer(workerPool chan *producers) {
}

func producer(ch chan int, i int) {
	ch <- i
}

func main() {
	ch := make(chan int)

	for i := 0; i <= 10; i++ {
		go producer(ch, i)
	}
}
