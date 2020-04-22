package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	//data, err := ioutil.ReadFile("fileOperator/autoUpgrade.json")
	//if err != nil {
	//	fmt.Println("File reading error", err)
	//	return
	//}
	//fmt.Println("Contents of file:", string(data))
	//
	fptr := flag.String("fpath", "autoUpgrade.json", "file path to read from")
	flag.Parse()
	data, err := ioutil.ReadFile(*fptr)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println("Contents of file:", string(data))


	//fptr := flag.String("fpath", "autoUpgrade.json", "file path to read from")
	//flag.Parse()
	//
	//f, err := os.Open(*fptr)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer func() {
	//	if err = f.Close(); err != nil {
	//		log.Fatal(err)
	//	}
	//}()
	//s := bufio.NewScanner(f)
	//for s.Scan() {
	//	fmt.Println(s.Text())
	//}
	//err = s.Err()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//********still have question why binary file wasn't generated*************
	// box := packr.NewBox(".")
	// json, err := box.FindString("autoUpgrade.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Contents of file:", json)
}