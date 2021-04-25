package main

import (
	"fmt"
)

//链表节点
type ListNode struct {
	Value int
	Next  *ListNode
}

//反转链表的实现
func reversrList(head *ListNode) (reverseHead *ListNode) {
	if head == nil {
		return
	}

	for head != nil {

		node := &ListNode{
			head.Value,
			reverseHead,
		}
		reverseHead = node

		head = head.Next
	}

	return
}

func main() {
	head := new(ListNode)
	head.Value = 1
	ln2 := new(ListNode)
	ln2.Value = 2
	ln3 := new(ListNode)
	ln3.Value = 3
	ln4 := new(ListNode)
	ln4.Value = 4
	ln5 := new(ListNode)
	ln5.Value = 5
	head.Next = ln2
	ln2.Next = ln3
	ln3.Next = ln4
	ln4.Next = ln5

	result := reversrList(head)
	for result != nil {
		fmt.Println(result)
		result = result.Next
	}
}
