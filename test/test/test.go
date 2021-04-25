package main

import (
	"errors"
	"fmt"
	"os"
)

func conv(str string) (int, error) {
	slice := []byte(str)

	// check the first character
	if slice[0] == '+' || slice[0] == '-' {
		slice = slice[1:]
		if len(slice) <= 0 {
			return 0, errors.New("error")
		}
	}

	//sum
	var n int
	for _, ch := range slice {
		ch = ch - '0'
		if ch > '9' {
			return 0, errors.New("error")
		}
		n = n*10 + int(ch)
	}

	if str[0] == '-' {
		n = -n
	}

	return n, nil
}

func reverse(str string) string{
	slice := []byte(str)

	for from, to := 0, len(slice) - 1 ;

	return
}

func main() {
	strArr := []string{
		"10",
		"-100",
		"+1000",
	}

	fmt.Println("*********************")
