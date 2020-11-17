package main

import "fmt"

func main() {
	var a = []int{12, 31, 11, 1, 1, 3, 9, 123, 1}

	quickSort(a, 0, len(a)-1)

	fmt.Println(a)
}

func quickSort(a []int, l, r int) {
	ol := l
	or := r

	if l > r || r < l {
		return
	}

	m := a[l]

	for l != r {
		for a[r] >= m && l != r {
			r--
		}
		if l != r {
			a[l] = a[r]
			l++
		}
		for a[l] < m && l != r {
			l++
		}
		if l != r {
			a[r] = a[l]
			r--
		}

	}

	a[l] = m

	quickSort(a, ol, l-1)
	quickSort(a, l+1, or)

}
