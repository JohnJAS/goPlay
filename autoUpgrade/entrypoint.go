package main

import (
	. "autoUpgrade/common"
	"fmt"
)

var TempFolder string

func init() {
	//TempFolder is the temp folder of autoUpgrade
	if SysType == "windows" {
		TempFolder = "C:\tmp\autoUpgrade"
	} else {
		TempFolder = "/tmp/autoUpgrade"
	}
}

func main() {
	fmt.Println("autoUpgrade.sh transfer to autoUpgrade.go...")

	fmt.Println(SysType)
	fmt.Println(TempFolder)
	fmt.Println(UpgradeLog)
}
