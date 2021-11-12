// Example application that uses all of the available API options.
package main

import (
	"goPlay/esaySpinner"
	"os"
	"time"
)

func main() {
	spinner := esaySpinner.New(esaySpinner.CharSets[0], 100*time.Millisecond, os.Stdout)
	spinner.Start()
	time.Sleep(3 * time.Second)
	spinner.Stop()
}
