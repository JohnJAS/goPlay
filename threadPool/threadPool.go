package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	id       int
	randomno int
}
type Result struct {
	job         Job
	sumofdigits int
}

var jobs = make(chan Job, 10)
var results = make(chan Result, 10)

func digits(number int) int {
	sum := 0
	no := number
	for no != 0 {
		digit := no % 10
		sum += digit
		no /= 10
	}
	time.Sleep(2 * time.Second)
	return sum
}
func worker(wg *sync.WaitGroup) {
	for job := range jobs { //监听jobs channel, 如果jobs没有关闭会阻塞,jobs关闭时推出循环
		output := Result{job, digits(job.randomno)}
		results <- output
	}
	wg.Done() //waitGroup计数器 - 1
}
func createWorkerPool(noOfWorkers int) {
	var wg sync.WaitGroup //新建一个waitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)      //waitGroup计数器 + 1
		go worker(&wg) //启动一个worker routine
	}
	wg.Wait() //等待计数器变为0
	close(results)
}
func allocate(noOfJobs int) {
	for i := 0; i < noOfJobs; i++ {
		randomno := rand.Intn(999)
		job := Job{i, randomno}
		jobs <- job
	}
	close(jobs)
}
func result(done chan bool) {
	for result := range results { //监听results channel, 直到results channel关闭, 推出循环
		fmt.Printf("Job id %d, input random no %d , sum of digits %d\n", result.job.id, result.job.randomno, result.sumofdigits)
	}
	done <- true
}
func main() {
	startTime := time.Now()
	noOfJobs := 100
	go allocate(noOfJobs)
	done := make(chan bool)
	go result(done)
	noOfWorkers := 20
	createWorkerPool(noOfWorkers)
	<-done //阻塞，除非done channel接收到值true, 否则程序不会进行下去
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("total time taken ", diff.Seconds(), "seconds")
}
