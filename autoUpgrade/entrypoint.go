package main

import (
	"bytes"
	"encoding/json"
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

	cdfK8S "autoUpgrade/cdfutil/k8s"
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

//LogLevel set log level in autoUpgrade log
var LogLevel int

//CDF version before upgrade(won't refresh)
var OrgCurrentVersion string

//CDF current version(refresh after an CDF version upgrade)
var CurrentVersion string

//USER_UPGRADE_PACKS : upgrade packages user provided, they should be placed correctly.
var USER_UPGRADE_PACKS []string

//Node list of target cluster
var NodeList = cdfCommon.NewNodeList([]cdfCommon.Node{}, 0)

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
	os.Args = append(os.Args, "shcCDFRH75vm01-0.hpeswlab.net")
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
			&cli.BoolFlag{
				Name:        "dry-run",
				Value:       false,
				Destination: &DryRun,
				Usage:       "Dry run for autoUpgrade.(Alpha)",
			},
			&cli.StringFlag{
				Name:  "verbose",
				Value: "DUBUG",
				Usage: "Set log level for autoUpgrade.(Alpha)",
			},
		},
		Action: startExec,
	}
	err := app.Run(os.Args)
	if err != nil {
		cdfLog.WriteLog(Logger, cdfCommon.FATAL, LogLevel, err.Error())
	}
}

//autoUpgrade main process
func startExec(c *cli.Context) (err error) {
	if c.Bool("verbose") {
		LogLevel = cdfLog.TransferLogLevel(c.Value("verbose").(string))
		if LogLevel == 0 {
			LogLevel = cdfCommon.DEBUG
			cdfLog.WriteLog(Logger, cdfCommon.WARN, LogLevel, fmt.Sprintf("Unsupportted input log level %s. Log level works in DUBUG mode.", c.Value("verbose")))
		}
	} else {
		LogLevel = cdfCommon.DEBUG
	}

	//main process start
	log.Println("===========================================================================")

	//init upgrade step
	err = initUpgradeStep()
	if err != nil {
		return
	}
	log.Println()

	//check connection to the cluster
	err = checkConnection(cdfCommon.NewNodeList([]cdfCommon.Node{cdfCommon.NewNode(NodeInCluster, "")}, 1))
	if err != nil {
		return
	}
	log.Println()

	//get cluster information
	err = getNodesInfo()
	if err != nil {
		return
	}
	log.Println()

	//get upgrade packages information
	err = getUpgradePacksInfo()
	if err != nil {
		return
	}
	log.Println()

	//check upgrade packages
	err = checkUpgradePacks()
	if err != nil {
		return
	}
	log.Println()

	err = checkParameters()
	if err != nil {
		return
	}
	log.Println()

	err = checkNodesInfo()
	if err != nil {
		return
	}
	log.Println()
	log.Println("===========================================================================")

	log.Println("Start to dynamic upgrade process...")
	return
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

	cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, "Current directory : "+CurrentDir)
	cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, "User input command: "+strings.Join(os.Args, " "))
}

//Determining start upgrade step...
func initUpgradeStep() error {
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "Determining start upgrade step...")
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
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "UPGRADE_STEP : "+result)
		cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, "Previous upgrade step execution results found. Continuing with step "+result)
		return nil
	} else {
		UpgradeStep = 0
		err := cdfOS.WriteFile(upgradeStepFilePath, UpgradeStep)
		return check(err)
	}
}

//check connection to the cluster nodes
func checkConnection(nodes cdfCommon.NodeList) (err error) {

	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "Checking connection to the cluster...")

	ch := make(chan cdfCommon.ConnectionStatus, nodes.Num)

	go func(chnl chan cdfCommon.ConnectionStatus) {
		for _, node := range nodes.List {
			err := cdfSSH.CheckConnection(node.Name, SysUser, KeyPath)
			if err != nil {
				chnl <- cdfCommon.ConnectionStatus{false, fmt.Sprintf("Failed to connect to node %s", node.Name)}
			} else {
				chnl <- cdfCommon.ConnectionStatus{true, fmt.Sprintf("Successfully connected to node %s", node.Name)}
			}

		}
		close(chnl)
	}(ch)

	for result := range ch {
		if result.Connected {
			cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, result.Description)
		} else {
			cdfLog.WriteLog(Logger, cdfCommon.ERROR, LogLevel, result.Description)
			err = errors.New("\nNode(s) unreachable found. Please check your SSH passwordless configuration and try again.")
		}
	}

	return
}

//Getting nodes info...
func getNodesInfo() (err error) {
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "Getting nodes information...")
	// get current cdf verison
	err = getCurrentVersion(false)
	if err != nil {
		return
	}

	err = getOrgVersion()
	if err != nil {
		return
	}

	err = getCurrentNodesInfo()
	if err != nil {
		return
	}

	return
}

func getOrgVersion() (err error) {
	exist, err := cdfOS.PathExists(filepath.Join(TempFolder, "OrgCurrentVersion"))
	if err != nil {
		return err
	}
	if ! exist {
		cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, "File OrgCurrentVersion not found.")
		OrgCurrentVersion = CurrentVersion
		err = cdfOS.WriteFile(filepath.Join(TempFolder, "OrgCurrentVersion"), OrgCurrentVersion)
	} else {
		cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, "File OrgCurrentVersion found.")
		OrgCurrentVersion, err = cdfOS.ReadFile(filepath.Join(TempFolder, "OrgCurrentVersion"))
	}
	cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, "OrgCurrentVersion: "+OrgCurrentVersion)
	return
}

//
func getCurrentNodesInfo() (err error) {
	exist, err := cdfOS.PathExists(filepath.Join(TempFolder, "Nodes"))
	if err != nil {
		return err
	}
	if !exist {
		// get current nodes info
		var stderr bytes.Buffer
		stderr, err = cdfK8S.GetCurrrentNodes(&NodeList, NodeInCluster, SysUser, KeyPath)
		if err != nil {
			cdfLog.WriteLog(Logger, cdfCommon.ERROR, LogLevel, stderr.String())
			return
		}
		cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, fmt.Sprintf("Node: %v", NodeList))
		var b []byte
		b, err = json.Marshal(NodeList)
		if err != nil {
			return
		}
		err = cdfOS.WriteFile(filepath.Join(TempFolder, "Nodes"), string(b))
		return
	} else {
		var content string
		content, err = cdfOS.ReadFile(filepath.Join(TempFolder, "Nodes"))
		if err != nil {
			return
		}
		err = json.Unmarshal([]byte(content), &NodeList)
		if err != nil {
			return
		}
		cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, fmt.Sprintf("Node: %v", NodeList))
	}
	printNodes(NodeList)
	return
}

func printNodes(nodeList cdfCommon.NodeList) {
	var n []string
	var m []string
	var w []string
	for _,node := range nodeList.List {
		n = append(n, node.Name)
		if node.Role == cdfCommon.MASTER {
			m = append(m, node.Name)
		}
		if node.Role == cdfCommon.WORKER {
			w = append(w, node.Name)
		}
	}
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("ALL_NODES       : %s",strings.Join(n," ")))
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("ALL_MASTERS     : %s",strings.Join(m," ")))
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("ALL_WORKERS     : %s",strings.Join(w," ")))
}

//Get current CDF version
func getCurrentVersion(update bool) error {
	exist, err := cdfOS.PathExists(filepath.Join(TempFolder, "CurrentVersion"))
	if err != nil {
		return err
	}
	if !exist || update {
		CurrentVersion, stdout, stderr, err := cdfK8S.GetCurrentVersion(NodeInCluster, SysUser, KeyPath)
		if err != nil {
			cdfLog.WriteLog(Logger, cdfCommon.ERROR, LogLevel, stderr.String())
			return err
		} else {
			cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, stdout.String())
			cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "CurrentVersion: "+CurrentVersion)
			err = cdfOS.WriteFile(filepath.Join(TempFolder, "CurrentVersion"), CurrentVersion)
			if err != nil {
				return err
			}
			return nil
		}
	} else {
		CurrentVersion, err = cdfOS.ReadFile(filepath.Join(TempFolder, "CurrentVersion"))
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "CURRENT_VERSION : "+CurrentVersion)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

//Getting upgrade package(s) information...
func getUpgradePacksInfo() (err error) {
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "Getting upgrade package(s) information...")
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "CURRENT_DIR : "+CurrentDir)

	pattern := []string{
		"version.txt",
		"upgrade.sh",
	}

	USER_UPGRADE_PACKS, err = cdfOS.ListDirWithFilter(cdfOS.ParentDir(CurrentDir),pattern,cdfOS.FilterAND)

	fmt.Println(USER_UPGRADE_PACKS)

	//sort packages

	//create version:path map

	return
}

//Checking upgrade package(s)...
func checkUpgradePacks() (err error) {
	cdfLog.WriteLog(Logger,cdfCommon.INFO,LogLevel,"Checking upgrade package(s)...")
	return
}

//Checking parameters(s)...
func checkParameters() (err error) {
	cdfLog.WriteLog(Logger,cdfCommon.INFO,LogLevel,"Checking parameters(s)...")
	return
}

//Checking nodes info...
func checkNodesInfo() (err error) {
	cdfLog.WriteLog(Logger,cdfCommon.INFO,LogLevel,"Checking nodes info...")
	return
}
