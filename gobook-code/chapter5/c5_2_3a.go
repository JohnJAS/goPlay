/*
代码所在章节：5.2.3节
*/

package main

import (
	"fmt"
	"sync"
)

//工作任务
type task struct {
	begin  int
	end    int
	result chan<- int //只写通道,send-only
}

//任务处理:计算begin到end的和
//执行结果写入到结果chan result中
func (t *task) do() {
	sum := 0
	for i := t.begin; i <= t.end; i++ {
		sum += i
	}
	t.result <- sum
}

func main() {
	//工作通道
	taskchan := make(chan task, 10)

	//结果通道
	resultchan := make(chan int, 10)

	//worker信号通道
	wg := &sync.WaitGroup{}

	//初始化task的goroutine,计算100个自然数之和
	go InitTask(taskchan, resultchan, 100)

	//分发任务在NUMBER个goroutine 池
	DistributeTask(taskchan, resultchan, wg)

	//通过结果通道处理结果
	sum := ProcessResult(resultchan)

	fmt.Println("sum=", sum)
}

//初始化待处理task chan
func InitTask(taskchan chan<- task, resultchan chan int, p int) {
	qu := p / 10    //每10个数一个区间
	mod := p % 10   //剩余不满10个数的单独一个区间
	high := qu * 10 //整区间最大值
	for j := 0; j < qu; j++ {
		b := 10*j + 1
		e := 10 * (j + 1)
		tsk := task{
			begin:  b, //10*j + 1
			end:    e, //10*j + 10
			result: resultchan,
		}
		taskchan <- tsk
	}
	if mod != 0 {
		tsk := task{
			begin:  high + 1,
			end:    p,
			result: resultchan,
		}
		taskchan <- tsk
	}

	close(taskchan)
}

//读取task chan 分发到worker goroutine 处理，workers的总的数量是workers
func DistributeTask(taskchan <-chan task, resultchan chan int, wg *sync.WaitGroup) {

	for task := range taskchan {
		wg.Add(1)
		go ProcessTask(task, wg)
	}
	wg.Wait()
	close(resultchan)
}

//工作goroutine处理具体工作，并将处理结构发送到结果chan
func ProcessTask(t task, wg *sync.WaitGroup) {
	t.do()
	wg.Done()

}

// 读取结果通道，汇总结果
func ProcessResult(resultchan chan int) int {
	sum := 0
	for r := range resultchan {
		sum += r
	}
	return sum
}
