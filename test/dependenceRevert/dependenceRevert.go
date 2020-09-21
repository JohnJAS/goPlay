/*
   模拟组装2台电脑
   --- 抽象层 ---
   有显卡Card  方法display
   有内存Memory 方法storage
   有处理器CPU   方法calculate

   --- 实现层层 ---
   有 Intel因特尔公司 、产品有(显卡、内存、CPU)
   有 Kingston 公司， 产品有(内存3)
   有 NVIDIA 公司， 产品有(显卡)

   --- 逻辑层 ---
   1. 组装一台Intel系列的电脑，并运行
   2. 组装一台 Intel CPU  Kingston内存 NVIDIA显卡的电脑，并运行
*/
package main

import "fmt"

//------  抽象层 -----
type Card interface{
	Display()
}

type Memory interface {
	Storage()
}

type CPU interface {
	Calculate()
}

type Computer struct {
	cpu CPU
	mem Memory
	card Card
}

func NewComputer(cpu CPU, mem Memory, card Card) *Computer{
	return &Computer{
		cpu:cpu,
		mem:mem,
		card:card,
	}
}

func (this *Computer) DoWork() {
	this.cpu.Calculate()
	this.mem.Storage()
	this.card.Display()
}

//------  实现层 -----
//intel
type IntelCPU struct {
	CPU
}

func (this *IntelCPU) Calculate() {
	fmt.Println("Intel CPU 开始计算了...")
}

type IntelMemory struct {
	Memory
}

func (this *IntelMemory) Storage() {
	fmt.Println("Intel Memory 开始存储了...")
}

type IntelCard struct {
	Card
}

func (this *IntelCard) Display() {
	fmt.Println("Intel Card 开始显示了...")
}

//kingston
type KingstonMemory struct {
	Memory
}

func (this *KingstonMemory) Storage() {
	fmt.Println("Kingston memory storage...")
}

//nvidia
type NvidiaCard struct {
	Card
}

func (this *NvidiaCard) Display() {
	fmt.Println("Nvidia card display...")
}



//------  业务逻辑层 -----
func main() {
	//intel系列的电脑
	com1 := NewComputer(&IntelCPU{}, &IntelMemory{}, &IntelCard{})
	com1.DoWork()

	//杂牌子
	com2 := NewComputer(&IntelCPU{}, &KingstonMemory{}, &NvidiaCard{})
	com2.DoWork()
}