package main

import (
	"fmt"
	"log"
	"os"

	. "autoUpgrade/common"
	"github.com/urfave/cli/v2"
)

//TempFolder is autoUpgrade temp folder including re-run mark and auto upgrade log
var TempFolder string
//UpgradeLog is log folder of autoUpgrade
var UpgradeLog string

var LogFile string

func init() {
	if SysType == "windows" {
		TempFolder = "C:\\tmp\\autoUpgrade"
	} else {
		TempFolder = "/tmp/autoUpgrade"
	}

	UpgradeLog = "upgradeLog"
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
