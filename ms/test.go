package main

func tranferPatter(pattern string) string{

}

type interval struct {
	from int,
	to int,
}

func is_match(pattern, str string)bool{
	patternSlice := []byte(pattern)

	for i, ch := range patternSlice {
		if ch == '{' {
			key := patternSlice[i-1]
		}
	}



	return false
}

func main(){

}