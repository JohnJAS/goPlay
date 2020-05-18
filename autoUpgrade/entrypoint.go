package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"

	cdfOSUtil "autoUpgrade/cdfutil/os"
	cdfCommon "autoUpgrade/common"
)

//TempFolder is autoUpgrade temp folder including re-run mark and auto upgrade log
var TempFolder string

//UpgradeLog is log folder of autoUpgrade
var UpgradeLog string

//LogFilePath is autoUpgrade logfile path
var LogFilePath string

//LogFile is autoUpgrade logfile
var LogFile *os.File

//CURRENT_DIR:
var CURRENT_DIR string

//UPGRADE_STEP: upgrade step already run
var UPGRADE_STEP int

//UPG_EXEC_CALL: upgrade exec call count, init 1 for the first call
var UPG_EXEC_CALL int

//NODE_IN_CLUSTER : because the script may not be in the cluster, users must provide a node in the cluster.
var NODE_IN_CLUSTER string

//WORK_DIR : upgrade work dictionary on the nodes in the cluster
var WORK_DIR string

//DRY_RUN: for autoUpgrade dry-run
var DRY_RUN bool

func init() {

	//identify system OS
	if cdfCommon.SysType == "windows" {
		TempFolder = os.Getenv("TEMP")
		if TempFolder == "" {
			log.Fatal("Failed to find system env TEMP, initailization failed.")
		}
		TempFolder = filepath.Join(TempFolder, "autoUpgrade")
	} else {
		TempFolder = "/tmp/autoUpgrade"
	}

	//create log file
	LogFilePath = filepath.Join(TempFolder, "upgradeLog", "autoUpgrade-"+time.Now().UTC().Format(cdfCommon.RFC3339)+".log")
	exist, err := cdfOSUtil.PathExists(LogFilePath)
	if err != nil {
		log.Fatalln(err)
	} else if exist == false {
		err = os.Mkdir("LogFilePath", 0666)
		if err != nil {
			fmt.Println(err)
		}
	}
	LogFile, err = os.Create(LogFilePath)
	defer LogFile.Close()
	if err != nil {
		log.Fatalln(err)
	}

}

func main() {
	app := &cli.App{
		Name:            "autoUpgrade",
		Usage:           "Upgrade CDF automatically.",
		UsageText:       "autoUpgrade [-d|--dir <working_directory>] [-n|--node <any_node_in_cluster>] [-u|--sysuser <system_user>] [-o|--options <input_options>]",
		Description:     "Requires passwordless SSH to be configured to all cluster nodes. If the script is not run on a cluster node, you must have passwordless SSH configured to all cluster nodes. If the script is run on a cluster node, you must have passwordless SSH configured to all cluster nodes including this node. You can learn more about the auto upgrade through the official document.",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "d",
				Aliases:     []string{"dir"},
				Required:    true,
				Destination: &WORK_DIR,
				Usage:       "The working directory to use on all cluster nodes. ENSURE the directory is empty and the file system as enough space. If you are a non-root user on the nodes inside the cluster, make sure you have permission to this directory.(mandatory)",
			},
			&cli.StringFlag{
				Name:        "n",
				Aliases:     []string{"node"},
				Required:    true,
				Destination: &NODE_IN_CLUSTER,
				Usage:       "IP address of any node inside the cluster.(mandatory)",
			},
			&cli.StringFlag{
				Name:    "u",
				Value:   "root",
				Aliases: []string{"sysuser"},
				Usage:   "The user for the SSH connection to the nodes inside the cluster. This user must have the permission to operate on the nodes inside the cluster. The configuration of the user must be done before running this script.(optional)",
			},
			&cli.StringFlag{
				Name:    "o",
				Aliases: []string{"options"},
				Usage:   "Set the options needed for each version of upgrade. For a single version, the rule is like '[upgradeVersion1]:[option1]=[value1],[option2]=[value2]'. Different versions use '|' to distinguish with others, like '[upgradeVersion1]:[option]=[value]|[upgradeVersion2]:[option]=[value]'.(optional)",
			},
			&cli.StringFlag{
				Name:  "dry-run",
				Value: "false",
				Usage: "Dry run for autoUpgrade.",
			},
		},
		Action: startExec,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func startExec(c *cli.Context) error {
	if c.Bool("dry-run") {
		if c.Value("dry-run") == "true" {
			DRY_RUN = true
		}
	}
	fmt.Println(WORK_DIR)
	fmt.Println(NODE_IN_CLUSTER)
	fmt.Println(DRY_RUN)
	fmt.Println(LogFilePath)
	return nil
}
