package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	. "autoUpgrade/common"
	"github.com/urfave/cli/v2"
)

//TempFolder is autoUpgrade temp folder including re-run mark and auto upgrade log
var TempFolder string
//UpgradeLog is log folder of autoUpgrade
var UpgradeLog string
//LogFile is autoUpgrade logfile path
var LogFilePath string

func init() {
	if SysType == "windows" {
		TempFolder = os.Getenv("TEMP1")
		if TempFolder == "" {
			log.Fatal("Failed to find system env TEMP, initailization failed.")
		}
		TempFolder = filepath.Join(TempFolder,"autoUpgrade")
	} else {
		TempFolder = "/tmp/autoUpgrade"
	}

	LogFilePath = filepath.Join(TempFolder,"upgradeLog")
}

func main() {
	fmt.Println("autoUpgrade.sh transfer to autoUpgrade.go...")

	fmt.Println(SysType)
	fmt.Println(TempFolder)
	fmt.Println(LogFilePath)

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
