package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	cdfLog "autoUpgrade/cdfutil/log"
	cdfOS "autoUpgrade/cdfutil/os"
	cdfSSH "autoUpgrade/cdfutil/ssh"
	cdfCommon "autoUpgrade/common"
	"github.com/urfave/cli/v2"
)

//TempFolder is autoUpgrade temp folder including re-run mark and auto upgrade log
var TempFolder string

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

//SysUser is the user of destination k8s cluster
var SysUser string

//KeyPath is the rsa key path
var KeyPath string

//DryRun for autoUpgrade dry-run
var DryRun bool

//Debug autoUpgrade debug mode
var Debug bool

func init() {
	var err error

	//identify system OS
	if cdfCommon.SysType == "windows" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal("Failed to get current user info, initailization failed.")
		}
		TempFolder = filepath.Join(usr.HomeDir, "tmp", "autoUpgrade")
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
	os.Args = append(os.Args, "shcCDFRH75vm02-0.hpeswlab.net")
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
				Name:        "u",
				Value:       "root",
				Aliases:     []string{"sysuser"},
				Destination: &SysUser,
				Usage:       "The user for the SSH connection to the nodes inside the cluster. This user must have the permission to operate on the nodes inside the cluster. The configuration of the user must be done before running this script.(optional)",
			},
			&cli.StringFlag{
				Name:        "i",
				Aliases:     []string{"rsakey"},
				Destination: &KeyPath,
				Usage:       "The RSA key for the SSH connection to the nodes inside the cluster.",
			},
			&cli.StringFlag{
				Name:    "o",
				Aliases: []string{"options"},
				Usage:   "Set the options needed for each version of upgrade. For a single version, the rule is like '[upgradeVersion1]:[option1]=[value1],[option2]=[value2]'. Different versions use '|' to distinguish with others, like '[upgradeVersion1]:[option]=[value]|[upgradeVersion2]:[option]=[value]'.(optional)",
			},
			&cli.StringFlag{
				Name:  "dry-run",
				Value: "false",
				Usage: "Dry run for autoUpgrade.(Developping)",
			},
			&cli.StringFlag{
				Name:  "debug",
				Value: "false",
				Usage: "Debug mode for autoUpgrade.(Developping)",
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

	if c.Bool("debug") {
		if c.Value("debug") == "true" {
			Debug = true
		}
	}

	//main process start
	log.Println("===========================================================================")

	//init upgrade step
	err = initUpgradeStep()
	if err != nil {
		return err
	}
	log.Println()

	//connect to the cluster
	err = checkConnection(cdfCommon.Nodes{
		[]string{NodeInCluster},
		1,
	})
	if err != nil {
		return err
	}
	log.Println()

	//get upgrade packages information
	err = getUpgradePacksInfo()

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
		UpgradeStep, err = strconv.Atoi(result)
		if err != nil {
			return err
		}
		cdfLog.WriteLog(Logger, cdfCommon.INFO, "UpgradeStep: "+result)
		cdfLog.WriteLog(Logger, cdfCommon.DEBUG, "Previous upgrade step execution results found. Continuing with step "+result)
		return nil
	} else {
		UpgradeStep = 0
		err := cdfOS.WriteFile(upgradeStepFilePath, UpgradeStep)
		return check(err)
	}
}

//
func checkConnection(nodes cdfCommon.Nodes) error {
	var err error
	cdfLog.WriteLog(Logger, cdfCommon.INFO, "Checking connection to the cluster...")

	ch := make(chan cdfCommon.ConnectionStatus, nodes.Num)

	go func(chnl chan cdfCommon.ConnectionStatus) {
		for _, node := range nodes.NodeList {
			err := cdfSSH.CheckConnection(node, SysUser, KeyPath)
			if err != nil {
				chnl <- cdfCommon.ConnectionStatus{false, fmt.Sprintf("Failed to connect to node %s", node)}
			} else {
				chnl <- cdfCommon.ConnectionStatus{true, fmt.Sprintf("Successfully connected to node %s", node)}
			}

		}
		close(chnl)
	}(ch)

	for result := range ch {
		if result.Connected {
			cdfLog.WriteLog(Logger, cdfCommon.INFO, result.Description)
		} else {
			cdfLog.WriteLog(Logger, cdfCommon.ERROR, result.Description)
			err = errors.New("\nNode(s) unreachable found. Please check your SSH passwordless configuration and try again.")
		}
	}

	return check(err)
}

//Getting upgrade package(s) information...
func getUpgradePacksInfo() error {
	return nil
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
