package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fileName := "C:\\Users\\shengj\\workspace\\bin\\jq"
	fmt.Println(fileName)
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fileInfo.Mode())


//	if err := os.Chmod(fileName, 777); err != nil {
//		log.Fatalln(err)
//	}
//
//	log.Println(fileInfo.Mode())
}
