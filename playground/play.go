// Example application that uses all of the available API options.
package main

import (
	"fmt"
	"goPlay/esaySpinner"
	"os"
	"time"
)

func main() {
	spinner := esaySpinner.New(esaySpinner.CharSets[0], 100*time.Millisecond, os.Stdout).WithHiddenCursor(true)
	spinner.Start()
	fmt.Println("123123")
	time.Sleep(3 * time.Second)
	fmt.Println("qweqwe")
	spinner.Stop()
}
