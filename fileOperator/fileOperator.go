package main

import (
	"fmt"
	"log"

	"github.com/gobuffalo/packr"
)

func main() {
	//data, err := ioutil.ReadFile("fileOperator/autoUpgrade.json")
	//if err != nil {
	//	fmt.Println("File reading error", err)
	//	return
	//}
	//fmt.Println("Contents of file:", string(data))
	//
	//fptr := flag.String("fpath", "fileOperator/autoUpgrade.json", "file path to read from")
	//flag.Parse()
	//fmt.Println("value of fpath is", *fptr)

	//still have question why binary file wasn't generated
	box := packr.NewBox("../testBox")
	json, err := box.FindString("autoUpgrade.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Contents of file:", json)
}
