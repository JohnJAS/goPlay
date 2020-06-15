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
	"sync"
	"time"

	cdfJson "autoUpgrade/cdfutil/json"
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
var UpgradeStep int = 0

//UpgExecCall upgrade exec call count, init 1 for the first call
var UpgExecCall int = 1

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

//OriginVersion CDF version before upgrade(won't refresh)
var OriginVersion string

//CurrentVersion current CDF version (refresh after an CDF version upgrade)
var CurrentVersion string

//TargetVersion CDF version user want to upgrade (it should be together with the autoUpgrade script)
var TargetVersion string

//UserUpgradePacks : upgrade packages user provided, they should be placed correctly.
var UserUpgradePacks []string

//UpgradeChain : the upgrade path that autoUpgrade supportted
var UpgradeChain []string

//UpgradePath : upgrade path that autoUpgrade will execute group by version
var UpgradePath []string

//InternalUpgradePath : internal upgrade path that autoUpgrade will execute group by internal version
var InternalUpgradePath []string

//NodeList of target cluster
var NodeList = cdfCommon.NewNodeList([]cdfCommon.Node{}, 0)

//VersionPathMap 202002:/path/package
var VersionPathMap = make(map[string]string)

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
	os.Args = append(os.Args, "/tmp/")
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
		cdfLog.WriteLog(Logger, cdfCommon.FATAL, LogLevel, err.Error(), LogFilePath)
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
	log.Println("=====================================================================================================")

	//init upgrade step
	err = initUpgradeStep()
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

	//check connection to all nodes
	err = checkConnection(NodeList)
	if err != nil {
		return
	}
	log.Println()

	//get upgrade path
	err = getUpgradePath()
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

	log.Println("Starting auto upgrade main process...")
	log.Println("=====================================================================================================")
	err = autoUpgrade()
	log.Println("=====================================================================================================")
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "Congratulations! Auto upgrade process is finished successfully!")
	return
}

func execFunc(f func(i ...interface{}) error, i ...interface{}) (err error) {
	if len(i) == 0 {
		err = f()
		log.Println()
		return err
	} else if len(i) == 1 {
		switch i[0].(type) {
		case string:
			err = f(i[0].(string))
			log.Println()
			return err
		case cdfCommon.NodeList:
			err = f(i[0].(cdfCommon.NodeList))
			log.Println()
			return err
		default:
			return errors.New("INTERNAL ERROR : Unknown type within one parameter")
		}
	} else {
		return errors.New("INTERNAL ERROR : Illegal parameter")
	}
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
		err := cdfOS.WriteFile(upgradeStepFilePath, UpgradeStep)
		return check(err)
	}
}

//check connection to the cluster nodes
func checkConnection(nodes cdfCommon.NodeList) (err error) {

	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "Checking connection to the cluster...")

	ch := make(chan cdfCommon.ConnectionStatus, nodes.Num)

	for _, nodeObj := range nodes.List {
		go func(node string) {
			err := cdfSSH.CheckConnection(node, SysUser, KeyPath)
			if err != nil {
				ch <- cdfCommon.ConnectionStatus{false, fmt.Sprintf("Failed to connect to node %s", node)}
			} else {
				ch <- cdfCommon.ConnectionStatus{true, fmt.Sprintf("Successfully connected to node %s", node)}
			}
		}(nodeObj.Name)
	}

	i := 0
	for result := range ch {
		if result.Connected {
			cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, result.Description)
		} else {
			cdfLog.WriteLog(Logger, cdfCommon.ERROR, LogLevel, result.Description)
			err = errors.New("\nNode(s) unreachable found. Please check your SSH passwordless configuration and try again.")
		}
		i++
		if i == nodes.Num {
			close(ch)
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

	// get the origin version before current upgrade
	err = getOrgVersion()
	if err != nil {
		return
	}

	//get current cluster info
	err = getCurrentNodesInfo()
	if err != nil {
		return
	}

	return
}

// get the origin version before current upgrade
func getOrgVersion() (err error) {
	exist, err := cdfOS.PathExists(filepath.Join(TempFolder, "OriginVersion"))
	if err != nil {
		return err
	}
	if !exist {
		cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, "File OriginVersion not found.")
		OriginVersion = CurrentVersion
		err = cdfOS.WriteFile(filepath.Join(TempFolder, "OriginVersion"), OriginVersion)
	} else {
		cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, "File OriginVersion found.")
		OriginVersion, err = cdfOS.ReadFile(filepath.Join(TempFolder, "OriginVersion"))
	}
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "ORIGIN_VERSION  : "+OriginVersion)
	return
}

//get current cluster info
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
	for _, node := range nodeList.List {
		n = append(n, node.Name)
		if node.Role == cdfCommon.MASTER {
			m = append(m, node.Name)
		}
		if node.Role == cdfCommon.WORKER {
			w = append(w, node.Name)
		}
	}
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("ALL_NODES       : %s", strings.Join(n, " ")))
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("ALL_MASTERS     : %s", strings.Join(m, " ")))
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("ALL_WORKERS     : %s", strings.Join(w, " ")))
}

//Get current CDF version
func getCurrentVersion(update bool) error {
	exist, err := cdfOS.PathExists(filepath.Join(TempFolder, "CurrentVersion"))
	if err != nil {
		return err
	}
	if !exist || update {
		var stdout, stderr bytes.Buffer
		CurrentVersion, stdout, stderr, err = cdfK8S.GetCurrentVersion(NodeInCluster, SysUser, KeyPath)
		if err != nil {
			cdfLog.WriteLog(Logger, cdfCommon.ERROR, LogLevel, stderr.String())
			return err
		} else {
			cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, stdout.String())
			cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "CURRENT_VERSION : "+CurrentVersion)
			err = cdfOS.WriteFile(filepath.Join(TempFolder, "CurrentVersion"), CurrentVersion)
			if err != nil {
				return err
			}
		}
	} else {
		CurrentVersion, err = cdfOS.ReadFile(filepath.Join(TempFolder, "CurrentVersion"))
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "CURRENT_VERSION : "+CurrentVersion)
		if err != nil {
			return err
		}
	}
	return nil
}

//Getting upgrade package(s) information...
func getUpgradePacksInfo() (err error) {
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "Getting upgrade package(s) information...")
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "CURRENT_DIR : "+CurrentDir)

	pattern := []string{
		cdfCommon.VersionTXT,
		cdfCommon.UpgradeSH,
	}

	UserUpgradePacks, err = cdfOS.ListDirWithFilter(cdfOS.ParentDir(CurrentDir), pattern, cdfOS.FilterAND)
	if err != nil {
		return
	}
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("USER_UPGRADE_PACKS : %s", strings.Join(UserUpgradePacks, " ")))

	//create version:path map
	err = initVersionPathMap()
	if err != nil {
		return
	}

	return
}

func getUpgradePath() (err error) {
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "Calculating upgrade path...")
	UpgradeChain, err = cdfJson.GetUpgradeChain(filepath.Join(CurrentDir, cdfCommon.AutoUpgradeJSON))
	if err != nil {
		return
	}
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("UPGRADE_CHAIN         : %s", strings.Join(UpgradeChain, " ")))

	fromVersion := transferVersionFormat(OriginVersion, false)
	targetVersion, err := cdfOS.ReadFile(filepath.Join(CurrentDir, cdfCommon.VersionTXT))
	if err != nil {
		return
	}
	targetVersion = transferVersionFormat(targetVersion, false)

	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("FROM_VERSION          : %s", fromVersion))
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("TARGET_VERSION        : %s", targetVersion))

	err = calculateUpgradePath(fromVersion, targetVersion)
	if err != nil {
		return err
	} else if UpgradePath == nil {
		return errors.New(fmt.Sprintf("No need to upgrade CDF from %s to %s", fromVersion, targetVersion))
	}

	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("CORRECT_UPGRADE_PATH  : %s", strings.Join(UpgradePath, " ")))
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("INTERNAL_UPGRADE_PATH : %s", strings.Join(InternalUpgradePath, " ")))

	return
}

//generate upgrade path
func generateUpgradePath(fromVersion string, targetVersion string, internal bool, wg *sync.WaitGroup) (err error) {
	startFlag := false
	finishFlag := false
	var isMajor, isVersionless bool

	for _, tempVersion := range UpgradeChain {
		//start record
		if tempVersion == fromVersion {
			startFlag = true
			continue
		}
		//recording upgrade path
		if startFlag == true && finishFlag == false {
			isMajor, err = cdfJson.GetIfMajor(filepath.Join(CurrentDir, cdfCommon.AutoUpgradeJSON), tempVersion)
			if err != nil {
				break
			}
			if internal {
				if isMajor || tempVersion == targetVersion {
					InternalUpgradePath = append(InternalUpgradePath, tempVersion)
				}
			} else {
				isVersionless, err = cdfJson.GetIfVersionless(filepath.Join(CurrentDir, cdfCommon.AutoUpgradeJSON), tempVersion)
				if err != nil {
					break
				}
				if isMajor && !isVersionless || tempVersion == targetVersion {
					UpgradePath = append(UpgradePath, tempVersion)
				}
			}
			//stop record
			if tempVersion == targetVersion {
				finishFlag = true
			}
			//exit
			if finishFlag == true {
				break
			}
		}
	}
	wg.Done()
	return
}

//calculate upgrade path
func calculateUpgradePath(fromVersion string, targetVersion string) (err error) {
	if fromVersion == targetVersion {
		return nil
	}


	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		errRT := generateUpgradePath(fromVersion, targetVersion, false, &wg)
		if errRT != nil {
			err = errRT
		}
	}()
	go func() {
		errRT := generateUpgradePath(fromVersion, targetVersion, true, &wg)
		if errRT != nil {
			err = errRT
		}
	}()
	wg.Wait()

	return
}

//verify upgrade path
func verifyUpgradePath() (err error) {
	return
}

//Checking upgrade package(s)...
func checkUpgradePacks() (err error) {
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "Checking upgrade package(s)...")
	return
}

//Checking parameters(s)...
func checkParameters() (err error) {
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "Checking parameters(s)...")
	return
}

//Checking nodes info...
func checkNodesInfo() (err error) {
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "Checking nodes info...")
	return
}

func initVersionPathMap() error {
	for _, pack := range UserUpgradePacks {
		path := filepath.Join(cdfOS.ParentDir(CurrentDir), pack)
		fullVersion, err := cdfOS.ReadFile(filepath.Join(path, cdfCommon.VersionTXT))
		if err != nil {
			return err
		}
		versionSlice := strings.Split(fullVersion, ".")
		if len(versionSlice) < 2 {
			return errors.New(fmt.Sprintf("Invaild format of %s under '%s'", cdfCommon.VersionTXT, path))
		}
		version := versionSlice[0] + versionSlice[1]
		VersionPathMap[version] = path
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("PACKAGE_VERSION : %s  PACKAGE_NAME : %s", version, pack))
	}
	return nil
}

//transfer version format
//example: 2020.02.002 to 202002 or 2020.02
func transferVersionFormat(input string, withDot bool) (result string) {
	versionSlice := strings.Split(input, ".")
	if withDot {
		result = versionSlice[0] + "." + versionSlice[1]
	} else {
		result = versionSlice[0] + versionSlice[1]
	}
	return
}

//start to upgrade CDF one version after one version
func autoUpgrade() (err error) {
	var message string
	for i, version := range UpgradePath {
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("** Starting upgrade CDF to %s **", version))
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("UPGRADE_ITERATOR : %d", i))
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("UPGRADE_VERSION  : %s", version))
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("UPGRADE_PACKAGE  : %s", VersionPathMap[version]))

		message = fmt.Sprintf("Copy %s upgrade package to all cluster nodes..", version)
		stepExec(cdfCommon.AllNodes, message, copyUpgradePacksToCluster, version, "")

		getCurrentVersion(true)
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("** Finished upgrade CDF to %s **", version))
		log.Println()
	}

	return
}

//stepExec
func stepExec(mode string, message string, f func(...string) error, version string, args string) (err error) {
	if UpgradeStep >= UpgExecCall {
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("Upgrade step '%d' '%s' already executed, continue to next one.", UpgExecCall, message))
		return
	}
	printUpgradeStep(UpgExecCall, message)

	err = f(mode, version, args)
	if err != nil {
		return
	}

	err = increaseUpgradeStep(UpgExecCall)

	return
}

//
func printUpgradeStep(step int, message string) {
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "")
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("--------- Starting UPGRADE-STEP %d \"%s\" ----------", step, message))
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "")
}

func increaseUpgradeStep(step int) (err error) {
	upgradeStepFilePath := filepath.Join(TempFolder, "UpgradeStep")
	err = cdfOS.WriteFile(upgradeStepFilePath, step)
	if err != nil {
		return
	}
	UpgradeStep = step
	UpgExecCall++

	return
}

func copyUpgradePacksToCluster(args ...string) (err error) {
	if len(args) < 2 {
		return errors.New("Internal Error in function copyUpgradePacksToCluster")
	}
	if len(args) >= 2 {
		mode := args[0]
		version := args[1]
		log.Println(mode)
		log.Println(version)
	}

	cdfSSH.CopyFile(NodeInCluster,SysUser,KeyPath,cdfCommon.AutoUpgradeJSON,filepath.ToSlash(filepath.Join(WorkDir,cdfCommon.AutoUpgradeJSON)))
	return
}

func dynamicUpgradeProcess() (err error) {
	return
}
