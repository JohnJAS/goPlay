package main

import (
	"fmt"
	"time"
)

type job struct {
	id int
}

func worker(id int, job <-chan job, results chan<- int) {
	//listening on workerPools
	for j := range job {
		fmt.Println("worker", id, "started  job", j.id)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j.id)
		results <- j.id * 2
	}
}

func main() {

	const jobsCount = 15
	results := make(chan int, jobsCount)

	const workersCount = 5
	workerPools := make(chan job, workersCount)

	//create workers to listen on workerPools
	for w := 1; w <= workersCount; w++ {
		go worker(w, workerPools, results)
	}

	//assign jobs
	for j := 1; j <= jobsCount; j++ {
		//there won't be deadlock here because goroutine workers are listing on workPools
		workerPools <- job{j}
	}
	close(workerPools)

	for a := 1; a <= jobsCount; a++ {
		<-results
	}
}
