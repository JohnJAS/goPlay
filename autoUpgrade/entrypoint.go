package main

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	cdfLog "autoUpgrade/cdfutil/log"
	cdfOS "autoUpgrade/cdfutil/os"
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

//Logger instance of log
var Logger *log.Logger

//CurrentDir current directory of autoUpgprade
var CurrentDir string

//UpgradeStep upgrade step already run
var UpgradeStep int

//UpgExecCall upgrade exec call count, init 1 for the first call
var UpgExecCall int

//NodeInCluster because the script may not be in the cluster, users must provide a node in the cluster.
var NodeInCluster string

//WorkDir upgrade work dictionary on the nodes in the cluster
var WorkDir string

//DryRun for autoUpgrade dry-run
var DryRun bool

func init() {
	var err error

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
	LogFilePath = filepath.Join(TempFolder, "upgradeLog", "autoUpgrade-"+time.Now().UTC().Format(cdfCommon.TIMESTAMP)+".log")
	LogFile, err = cdfOS.CreateFile(LogFilePath)
	defer LogFile.Close()
	if err != nil {
		log.Fatalln(err)
	}

	//get current directory
	CurrentDir, err = os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

}

func main() {
	//DEBUG_MODE
	os.Args = append(os.Args, "-n")
	os.Args = append(os.Args, "1.2.3.4")
	os.Args = append(os.Args, "-d")
	os.Args = append(os.Args, "./dir")
	startLog()
	defer LogFile.Close()

	app := &cli.App{
		Name:            "autoUpgrade",
		Usage:           "Upgrade CDF automatically.",
		UsageText:       "autoUpgrade [-d|--dir <working_directory>] [-n|--node <any_NodeInCluster>] [-u|--sysuser <system_user>] [-o|--options <input_options>]",
		Description:     "Requires passwordless SSH to be configured to all cluster nodes. If the script is not run on a cluster node, you must have passwordless SSH configured to all cluster nodes. If the script is run on a cluster node, you must have passwordless SSH configured to all cluster nodes including this node. You can learn more about the auto upgrade through the official document.",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "d",
				Aliases:     []string{"dir"},
				Required:    true,
				Destination: &WorkDir,
				Usage:       "The working directory to use on all cluster nodes. ENSURE the directory is empty and the file system has enough space. If you are a non-root user on the nodes inside the cluster, make sure you have permission to this directory.(mandatory)",
			},
			&cli.StringFlag{
				Name:        "n",
				Aliases:     []string{"node"},
				Required:    true,
				Destination: &NodeInCluster,
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

//autoUpgrade main process
func startExec(c *cli.Context) error {
	var err error

	if c.Bool("dry-run") {
		if c.Value("dry-run") == "true" {
			DryRun = true
		}
	}

	//main process
	log.Println("===========================================================================")
	//init upgrade step
	err = initUpgradeStep()

	if err != nil {
		return err
	}
	return nil
}

func check(err error) error {
	if err != nil {
		return err
	}
	return nil
}

func startLog() {
	var err error
	LogFile, err = cdfOS.OpenFile(LogFilePath)
	if err != nil {
		log.Fatal(err)
	}
	//initialize logger
	Logger = log.New(LogFile, "", 0)

	cdfLog.WriteLog(Logger, cdfCommon.DEBUG, "Current directory : "+CurrentDir)
	cdfLog.WriteLog(Logger, cdfCommon.DEBUG, "User input command: "+strings.Join(os.Args, " "))
}

//Determining start upgrade step...
func initUpgradeStep() error {
	cdfLog.WriteLog(Logger, cdfCommon.INFO, "Determining start upgrade step...")
	upgradeStepFilePath := filepath.Join(TempFolder, "UpgradeStep")
	exist, _ := cdfOS.PathExists(upgradeStepFilePath)
	if exist == true {
		result, err := cdfOS.ReadFile(upgradeStepFilePath, 1024)
		if err != nil && err != io.EOF {
			return err
		} else if result == "" {
			return errors.New("Fail to get UpgradeStep.")
		}
		UpgradeStep, _ = strconv.Atoi(result)
		cdfLog.WriteLog(Logger, cdfCommon.INFO, "UpgradeStep: "+result)
		cdfLog.WriteLog(Logger, cdfCommon.DEBUG, "Previous upgrade step execution results found. Continuing with step "+string(UpgradeStep))
		return nil
	} else {
		UpgradeStep = 0
		err := cdfOS.WriteFile(upgradeStepFilePath, UpgradeStep)
		return check(err)
	}
}

//Getting upgrade package(s) information...
func getUpgradePacksInfo() {

}

//Checking connection to the cluster nodes
func checkConnection() {

}

//Getting nodes info...
func getNodesInfo() {

}

//Checking upgrade package(s)...
func checkUpgradePacks() {

}

//Checking parameters(s)...
func checkParameters() {

}

//Checking nodes info...
func checkNodesInfo() {

}
