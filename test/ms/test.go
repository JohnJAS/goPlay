package main

import (
	"fmt"
)

type interVal struct {
	from int
	to   int
}

func findKey(pattern []rune) (key string) {
	for i, ch := range pattern {
		if ch == '{' {
			//ignore index out of range error here
			key = string(pattern[i-1])
			break
		}
	}
	return key
}

func tranferPattern(pattern []rune, key string) (string, interVal) {
	var index []int
	for i, ch := range pattern {
		if key == string(ch) {
			index = append(index, i)
		}
	}
	for i, v := range index{

	}

	return "", interVal{0, 0}
}

func is_match(pattern, input string) bool {
	patternSlice := []rune(pattern)
	inputSlice := []rune(input)

	key := findKey(patternSlice)
	//fmt.Println(key)
	//if key not found, compare pattern & input directly
	if key == "" {
		if pattern == input {
			return true
		}
		return false
	}

	pattern, interVal := tranferPattern(patternSlice, key)
	pattern = "of{4,7}ice"
	interVal := interVal{4, 7}
	fmt.Println(pattern)
	fmt.Println(interVal)

	count := 0
	for _, ch := range inputSlice {
		if string(ch) == key {
			count++
		}
	}
	if count >= interVal.from && count <= interVal.to {
		return true
	}

	return false
}

func main() {
	fmt.Println(is_match("of{2,4}ff{1,2}ice", "offffice"))
	fmt.Println(is_match("of{2,4}ff{1,2}ice", "offffffffffffice"))
}
