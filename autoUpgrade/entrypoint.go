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

var CURRENT_DIR string

//upgrade step already run
var UPGRADE_STEP int = 0
//upgrade exec call count, init 1 for the first call
var UPG_EXEC_CALL int = 1

//Because the script may not be in the cluster, users must provide a node in the cluster. 
var NODE_IN_CLUSTER string
//upgrade work dictionary on the nodes in the cluster
var WORK_DIR string

func init() {
	if SysType == "windows" {
		TempFolder = os.Getenv("TEMP")
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
	app := &cli.App{
		Name: "autoUpgrade",
    	Usage: "Upgrade CDF with one command! You can learn more about the auto upgrade through the official document.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "d",
				Aliases: []string{"dir"},
				Required: true,
				Destination: &WORK_DIR,
				Usage:   "The working directory to use on all cluster nodes.Ensure the directory is empty and the file system as enough space. If you are a non-root user on the nodes inside the cluster, make sure you have permission to this directory.\n",
			},
			&cli.StringFlag{
				Name:    "n",
				Aliases: []string{"node"},
				Required: true,
				Destination: &NODE_IN_CLUSTER,
				Usage:   "IP address of any node inside the cluster.This parameter is mandatory.\n",
			},
			&cli.StringFlag{
				Name:    "u",
				Value:  "root",
				Aliases: []string{"sysuser"},
				Usage:   "The user for the SSH connection to the nodes inside the cluster. This user must have the permission to operate on the nodes inside the cluster. The configuration of the user must be done before running this script. This parameter is optional.",
			},
			&cli.StringFlag{
				Name:    "o",
				Aliases: []string{"options"},
				Usage:   "Set the options needed for each version of upgrade. For a single version, the rule is shown below. [upgradeVersion1]:[option1]=[value1],[option2]=[value2] Different versions use '|' to distinguish with others. [upgradeVersion1]:[option]=[value]|[upgradeVersion2]:[option]=[value] This parameter is optional.\n",
			},
		},
	}

	err := app.Run(os.Args)
	fmt.Println(WORK_DIR)
	fmt.Println(NODE_IN_CLUSTER)
	if err != nil {
		log.Fatal(err)
	}
}
