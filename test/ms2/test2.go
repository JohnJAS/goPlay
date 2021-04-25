package main

import (
	"fmt"
	"regexp"
)

func is_match(pattern, input string) bool {
	validation := regexp.MustCompile(pattern)
	return validation.MatchString(input)
}

func main() {
	//test case
	fmt.Println(is_match("of{2,4}ff{1,2}ice", "offffice"))
	fmt.Println(is_match("of{2,4}ff{1,2}ice", "offffffffffffice"))
	fmt.Println(is_match("of{2,4}sf{1,2}ice", "offffsffice"))
	fmt.Println(is_match("of{2,4}sf{1,2}ice", "offffsfffice"))
	fmt.Println(is_match("of{2,4}sf{1,2}id{1,2}ce", "offffsffiddce"))
	fmt.Println(is_match("offffsffiddce", "offffsffiddce"))
	fmt.Println(is_match("", "offffsffiddce"))

	//result result:
	//true
	//false
	//true
	//false
	//true
	//true
	//true
}