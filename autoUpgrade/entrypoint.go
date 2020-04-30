package main

import (
	"fmt"
	"log"
	"os"

	. "autoUpgrade/common"
	"github.com/urfave/cli/v2"
)

var TempFolder string
var LogFile string

func init() {
	//TempFolder is the temp folder of autoUpgrade
	if SysType == "windows" {
		TempFolder = "C:\\tmp\\autoUpgrade"
	} else {
		TempFolder = "/tmp/autoUpgrade"
	}
}

func main() {
	fmt.Println("autoUpgrade.sh transfer to autoUpgrade.go...")

	fmt.Println(SysType)
	fmt.Println(TempFolder)
	fmt.Println(UpgradeLog)

	app := &cli.App{
		Name: "greet",
		Usage: "fight the loneliness!",
		Action: func(c *cli.Context) error {
			fmt.Println("Hello friend!")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
