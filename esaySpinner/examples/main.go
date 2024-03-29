// Example application that uses all of the available API options.
package main

import (
	"fmt"
	"goPlayg/esaySpinner"
	"os"
	"sync"
	"time"
)

func main() {
	spinner := esaySpinner.New(esaySpinner.CharSets[0], &sync.RWMutex{}, 100*time.Millisecond, os.Stdout).WithFinalMSG("\033[32m[OK]\033[0m\n")
	spinner.Start()
	fmt.Printf("123123")
	time.Sleep(3 * time.Second)
	fmt.Printf("asdasd")
	spinner.Stop()
	fmt.Println("New line!")
}
