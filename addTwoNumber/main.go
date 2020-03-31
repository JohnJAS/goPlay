package addtwonumber

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func listLength(l ListNode) int {
	var num int
	for l.Next != nil {
		l = *l.Next
		num++
	}
	num++
	return num
}

func calculate(val1 int, val2 int, more int) (int, int) {
	return (val1 + val2 + more) % 10, (val1 + val2 + more) / 10
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	len1 := listLength(*l1)
	len2 := listLength(*l2)

	if len1 < len2 {
		l1, l2 = l2, l1
		len1, len2 = len2, len1
	}

	more := 0

	head := l1

	for i := 0; i < len1+1; i++ {

		l1.Val, more = calculate(l1.Val, l2.Val, more)

		if l1.Next == nil {
			if more == 1 {
				l1.Next = &ListNode{1, nil}
			}
			break

		} else {
			l1 = l1.Next
		}

		if l2.Next == nil {
			l2 = &ListNode{0, nil}
		} else {
			l2 = l2.Next
		}

	}

	return head
}

func addNode(node *ListNode, i int) *ListNode {
	return &ListNode{i, node}
}

func main() {
	l1 := &ListNode{3, nil}

	l1 = addNode(l1, 8)
	l1 = addNode(l1, 1)

	l2 := &ListNode{1, nil}

	l2 = addNode(l2, 7)

	result := addTwoNumbers(l1, l2)

	length := listLength(*result)

	for i := 0; i < length; i++ {

		fmt.Println(result.Val)

		result = result.Next
	}

}
