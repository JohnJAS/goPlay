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

	"golang.org/x/crypto/ssh"

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

//SysGroup is the user's group of destination k8s cluster
var SysGroup string

//KeyPath is the rsa key path
var KeyPath string

//PassWord is the ssh password
var PassWord string

//Port is SSH port
var Port string

//Debug for autoUpgrade debug
var Debug bool

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

//VersionPackMap 202002:packname
var VersionPackMap = make(map[string]string)

//VersionPathMap 202002:/path/package
var VersionPathMap = make(map[string]string)

//InternalVersionMap 202002:202008
var InternalVersionMap = make(map[string]string)

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

	//get current directory
	CurrentDir, err = os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	TargetVersion, err = cdfOS.ReadFile(filepath.Join(CurrentDir, cdfCommon.VersionTXT))
	if err != nil {
		log.Fatalln(err)
	}

}

func main() {
	//DEBUG_MODE
	os.Args = append(os.Args, "-n")
	os.Args = append(os.Args, "shcCDFRH75vm01-0.hpeswlab.net")
	os.Args = append(os.Args, "-d")
	os.Args = append(os.Args, "/tmp/workspaceInCluster")
	//os.Args = append(os.Args, "--debug")
	os.Args = append(os.Args, "--dry-run")
	startLog()
	defer LogFile.Close()

	app := &cli.App{
		Name:            "autoUpgrade",
		Usage:           "Upgrade CDF automatically.",
		Description:     "Requires passwordless SSH to be configured to all cluster nodes. If the script is not run on a cluster node, you must have passwordless SSH configured to all cluster nodes. If the script is run on a cluster node, you must have passwordless SSH configured to all cluster nodes including this node. You can learn more about the auto upgrade through the official document.",
		UsageText:       "autoUpgrade [-d|--dir <working_directory>] [-n|--node <any_NodeInCluster>]",
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
				Name:        "g",
				Aliases:     []string{"sysgroup"},
				Destination: &SysGroup,
				Usage:       "The user for the SSH connection to the nodes inside the cluster. This user must have the permission to operate on the nodes inside the cluster. The configuration of the user must be done before running this script.(optional)",
			},
			&cli.StringFlag{
				Name:        "i",
				Aliases:     []string{"rsakey"},
				Destination: &KeyPath,
				Usage:       "The RSA key for the SSH connection to the nodes inside the cluster.",
			},
			&cli.StringFlag{
				Name:        "p",
				Aliases:     []string{"password"},
				Destination: &PassWord,
				Usage:       "The password for the SSH connection to the nodes inside the cluster.",
			},
			&cli.StringFlag{
				Name:        "port",
				Value:       "22",
				Destination: &Port,
				Usage:       "The port for the SSH connection to the nodes inside the cluster.",
			},
			&cli.StringFlag{
				Name:    "o",
				Aliases: []string{"options"},
				Usage:   "Set the options needed for each version of upgrade. For a single version, the rule is like '[upgradeVersion1]:[option1]=[value1],[option2]=[value2]'. Different versions use '|' to distinguish with others, like '[upgradeVersion1]:[option]=[value]|[upgradeVersion2]:[option]=[value]'.(optional)",
			},
			&cli.BoolFlag{
				Name:        "debug",
				Value:       false,
				Destination: &Debug,
				Usage:       "Debug mode for autoUpgrade.",
			},
			&cli.BoolFlag{
				Name:        "dry-run",
				Value:       false,
				Destination: &DryRun,
				Usage:       "Dry run for autoUpgrade.",
			},
			&cli.StringFlag{
				Name:  "verbose",
				Value: "DUBUG",
				Usage: "Set log level for autoUpgrade.",
			},
		},
		Action: startExec,
	}
	err := app.Run(os.Args)
	if err != nil {
		cdfLog.WriteLog(Logger, cdfCommon.FATAL, LogLevel, err.Error(), LogFilePath)
		os.Exit(1)
	}
}

//autoUpgrade main process
func startExec(c *cli.Context) (err error) {

	if !c.Bool("g") || !c.Bool("sysgroup") {
		SysGroup = SysUser
		if Debug {
			log.Println("SysUser  : " + SysUser)
			log.Println("SysGroup : " + SysGroup)
		}
	}

	if c.Bool("verbose") {
		LogLevel = cdfLog.TransferLogLevel(c.Value("verbose").(string))
		if LogLevel == 0 {
			LogLevel = cdfCommon.DEBUG
			cdfLog.WriteLog(Logger, cdfCommon.WARN, LogLevel, fmt.Sprintf("Unsupportted input log level %s. Log level works in DUBUG mode.", c.Value("verbose")))
		}
	} else {
		LogLevel = cdfCommon.DEBUG
	}

	if Debug {
		log.Println("RsaKey   : " + KeyPath)
		log.Println("PassWord : " + PassWord)
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
	if err != nil {
		return
	}
	log.Println("=====================================================================================================")
	err = deleteTempFolder()
	defer LogFile.Close()
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "Congratulations! Auto upgrade process is finished successfully!")
	return
}

func deleteTempFolder() (err error) {
	LogFile.Close()

	TempFolderBK := TempFolder + "BK"
	logFileFolder := filepath.Join(TempFolder, "upgradeLog")
	logFileFolderBK := filepath.Join(TempFolderBK, "upgradeLog")

	exist, _ := cdfOS.PathExists(TempFolderBK)
	if exist {
		os.RemoveAll(TempFolderBK)
	}

	os.Rename(TempFolder, TempFolderBK)

	LogFilePath = filepath.Join(logFileFolderBK, strings.TrimPrefix(LogFilePath, logFileFolder))
	LogFile, err = cdfOS.OpenFile(LogFilePath)
	if err != nil {
		log.Fatal(err)
	}

	Logger = log.New(LogFile, "", 0)

	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "Remove temp folder successfully.")

	return
}

func startLog() {
	var err error

	//create log file
	LogFilePath = filepath.Join(TempFolder, "upgradeLog", "autoUpgrade-"+time.Now().UTC().Format(cdfCommon.TIMESTAMP)+".log")
	LogFile, err = cdfOS.CreateFile(LogFilePath)
	if err != nil {
		log.Fatalln(err)
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
		result, err := cdfOS.ReadFile(upgradeStepFilePath)
		if err != nil && err != io.EOF {
			return err
		} else if result == "" {
			return errors.New("Fail to get UpgradeStep")
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
		if err != nil {
			return err
		}
		return nil
	}
}

//check connection to the cluster nodes
func checkConnection(nodes cdfCommon.NodeList) (err error) {

	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, "Checking connection to the cluster...")

	ch := make(chan cdfCommon.ConnectionStatus, nodes.Num)

	for _, nodeObj := range nodes.List {
		go func(node string) {
			err := cdfSSH.CheckConnection(node, SysUser, KeyPath, PassWord, Port)
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
			err = errors.New("Node(s) unreachable found. Please check your SSH passwordless configuration or password and try again")
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
		stderr, err = cdfK8S.GetCurrrentNodes(&NodeList, NodeInCluster, SysUser, KeyPath, PassWord, Port)
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
		CurrentVersion, stdout, stderr, err = cdfK8S.GetCurrentVersion(NodeInCluster, SysUser, KeyPath, PassWord, Port)
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
	targetVersion := transferVersionFormat(TargetVersion, false)

	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("FROM_VERSION          : %s", fromVersion))
	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("TARGET_VERSION        : %s", targetVersion))

	err = calculateUpgradePath(fromVersion, targetVersion)
	if err != nil {
		return err
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

func getSequenceFromUpgradeChain(version string) (sequence int) {
	// -1 means not found
	sequence = -1

	for i, x := range UpgradeChain {
		if version == x {
			sequence = i + 1
			break
		}
	}
	if sequence == -1 {
		// Upgrade chain (generate from json file) don't contain the CDF version where it started,
		// so check if it is at the very beginning.
		// Basiclly, it means the 201811 which is where we start to support autoUpgrade.
		if version == "201811" {
			sequence = 0
		}
	}

	return
}

func verifyFrom2Target(fromVersion, targetVersion string) (err error) {
	startSeq := getSequenceFromUpgradeChain(fromVersion)
	targetSeq := getSequenceFromUpgradeChain(targetVersion)

	if Debug {
		log.Println(fmt.Sprintf("startSeq : %v", startSeq))
		log.Println(fmt.Sprintf("startSeq : %v", targetSeq))

	}

	if startSeq == -1 {
		return fmt.Errorf("Upgrading CDF from %s is not supportted in autoUpgrade route. The earliest supportted version is 201811", fromVersion)
	}

	if targetSeq == -1 {
		return fmt.Errorf("Upgrading CDF from %s to %s is not supportted", fromVersion, targetVersion)
	}

	if targetSeq < startSeq {
		return fmt.Errorf("You can't upgrade currently installed %s CDF with %s package. Please use a newer upgrade package", fromVersion, targetVersion)
	}

	if targetSeq == startSeq {
		return fmt.Errorf("No need to upgrade CDF from %s to %s", fromVersion, targetVersion)
	}

	return
}

//calculate upgrade path
func calculateUpgradePath(fromVersion string, targetVersion string) (err error) {

	err = verifyFrom2Target(fromVersion, targetVersion)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var errRT1, errRT2 error
	wg.Add(2)
	go func() {
		errRT1 = generateUpgradePath(fromVersion, targetVersion, false, &wg)
	}()
	go func() {
		errRT2 = generateUpgradePath(fromVersion, targetVersion, true, &wg)
	}()
	wg.Wait()

	if errRT1 != nil || errRT2 != nil {
		err = fmt.Errorf("%v%v", errRT1, errRT2)
	}

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
	//check node info check mark

	//copy precheck to workspace, chmod&chown, run precheck concurrently
	nodesObj := NodeList.List
	var nodes []string
	for _, nodeObj := range nodesObj {
		nodes = append(nodes, nodeObj.Name)
	}

	ch := make(chan cdfCommon.ExecStatus, NodeList.Num)

	for _, node := range nodes {
		go func(node string) {
			var conn *ssh.Client
			if err == nil {
				conn, err = cdfSSH.CreatSSHClient(node, SysUser, KeyPath, PassWord, Port)
			}

			cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, fmt.Sprintf("Creating work directory on node %s ...", node))
			var cmd string
			if err == nil {
				//cmd = fmt.Sprintf("rm -rf %s/", filepath.ToSlash(WorkDir))
				//cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, node+" : "+cmd)
				//err = cdfSSH.SSHExecCmd(node, SysUser, KeyPath, PassWord, Port, cmd, true)
				err = cdfSSH.RemoveRemoteFolder(conn, filepath.ToSlash(WorkDir))
			}

			if err == nil {
				//cmd = fmt.Sprintf("mkdir -p %s/", filepath.ToSlash(WorkDir))
				//cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, node+" : "+cmd)
				//err = cdfSSH.SSHExecCmd(node, SysUser, KeyPath, PassWord, Port, cmd, true)
				err = cdfSSH.CreateRemoteFolder(conn, filepath.ToSlash(WorkDir))
			}

			if err == nil {
				cmd = fmt.Sprintf("chown %s:%s %s/", SysUser, SysGroup, filepath.ToSlash(WorkDir))
				cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, node+" : "+cmd)
				err = cdfSSH.SSHExecCmd(node, SysUser, KeyPath, PassWord, Port, cmd, true)
			}

			cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, fmt.Sprintf("Copying upgrade precheck script to %s ...", node))

			var srcFile, targetFile string

			if err == nil {
				//copy files with perm
				srcFile = filepath.Join(CurrentDir, cdfCommon.UpgradePreCheckSH)
				targetFile = filepath.ToSlash(filepath.Join(WorkDir, filepath.Base(cdfCommon.UpgradePreCheckSH)))
				err = cdfSSH.CopyFileLocal2Remote(conn, srcFile, targetFile, 700)
			}

			if err == nil {
				//copy files with perm
				srcFile = filepath.Join(CurrentDir, cdfCommon.AutoUpgradeJSON)
				targetFile = filepath.ToSlash(filepath.Join(WorkDir, filepath.Base(cdfCommon.AutoUpgradeJSON)))
				err = cdfSSH.CopyFileLocal2Remote(conn, srcFile, targetFile, 600)
			}

			var outbuf, errbuf bytes.Buffer
			if err == nil {
				cmdPath := filepath.ToSlash(filepath.Join(WorkDir))
				fromVersion := transferVersionFormat(OriginVersion, false)
				targetVersion := transferVersionFormat(TargetVersion, false)
				cmd = fmt.Sprintf("%s/upgradePreCheck.sh -f %s -t %s", cmdPath, fromVersion, targetVersion)
				outbuf, errbuf, err = cdfSSH.SSHExecCmdReturnResult(node, SysUser, KeyPath, PassWord, Port, cmd)
			}

			if Debug {
				log.Println(cmd)
				fmt.Println(outbuf.String())
				fmt.Println(outbuf.String())
				fmt.Println(err)
			}

			if err == nil {
				ch <- cdfCommon.ExecStatus{true, node, fmt.Sprintf("Node: %s precheck passed.", node)}
			} else {
				ch <- cdfCommon.ExecStatus{false, node, fmt.Sprintf("Node: %s precheck failed. Output：%s Error: %s", node, outbuf.String(), errbuf.String())}
			}
		}(node)
	}

	i := 0
	for result := range ch {
		if result.Executed {
			cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, result.Description)
		} else {
			cdfLog.WriteLog(Logger, cdfCommon.ERROR, LogLevel, result.Description)
			err = fmt.Errorf("Not all cluster nodes has passed the upgrade precheck.")
		}
		i++
		if i == NodeList.Num {
			close(ch)
		}
	}
	if err != nil {
		return
	}
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
			return fmt.Errorf("Invaild format of %s under '%s'", cdfCommon.VersionTXT, path)
		}
		version := versionSlice[0] + versionSlice[1]
		VersionPackMap[version] = pack
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
	for i, version := range UpgradePath {
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("******* Starting upgrade CDF to %s *******", version))
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("UPGRADE_ITERATOR : %d", i))
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("UPGRADE_VERSION  : %s", version))
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("UPGRADE_PACKAGE  : %s", VersionPathMap[version]))

		err = prepareClusterWorkSpace(version)
		if err != nil {
			return
		}

		err = autoUpgradeMainProcess(version)
		if err != nil {
			return
		}

		err = getCurrentVersion(true)
		if err != nil {
			return
		}
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("******* Finished upgrade CDF to %s *******", version))
		log.Println()
	}

	err = cleanWorkDirInsideCluster(UpgradePath[len(UpgradePath)-1])
	return
}

//stepExec
func stepExec(mode string, message string, f func(...string) error, version string, args string, order string) (err error) {
	if UpgradeStep >= UpgExecCall {
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("UPGRADE-STEP %d \"%s\" already executed, continue to next step...", UpgExecCall, message))
		UpgExecCall++
		return
	}
	printUpgradeStep(UpgExecCall, message)

	err = f(mode, version, args, order)
	if err != nil {
		return
	}

	err = increaseUpgradeStep(UpgExecCall)

	return
}

//print upgrade step
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

func transferMode(mode string) (nodeList cdfCommon.NodeList) {
	switch mode {
	case cdfCommon.SingleMaster:
		nodeList.AddNode(cdfCommon.NewNode(NodeInCluster, ""))
	case cdfCommon.AllMasters:
		for _, node := range NodeList.List {
			if node.Role == cdfCommon.MASTER {
				nodeList.AddNode(node)
			}
		}
	case cdfCommon.AllWorkers:
		for _, node := range NodeList.List {
			if node.Role == cdfCommon.WORKER {
				nodeList.AddNode(node)
			}
		}
	case cdfCommon.AllNodes:
		nodeList = NodeList
	default:
		panic("Wrong usage of function transferMode in main package.")
	}
	return
}

func prepareClusterWorkSpace(version string) (err error) {
	message := fmt.Sprintf("Copy %s upgrade package to all cluster nodes..", version)
	err = stepExec(cdfCommon.AllNodes, message, copyUpgradePacksToCluster, version, "", "")
	return
}

func autoUpgradeMainProcess(version string) (err error) {

	err = dynamicChildUpgradeProcess(version)
	if err != nil {
		return
	}

	err = dynamicUpgradeProcess(version)
	if err != nil {
		return
	}

	return
}

func dynamicChildUpgradeProcess(version string) (err error) {
	if Debug {
		log.Println(fmt.Sprintf("\n---------------Walking in dynamicChildUpgradeProcess: %s ---------------", version))
	}

	jsonPath := filepath.Join(CurrentDir, cdfCommon.AutoUpgradeChildJSON)
	autoUpgradeJSONObj, err := cdfJson.GetAutoUpgradeJsonObj(jsonPath)
	if err != nil {
		return
	}

	if Debug {
		log.Println(fmt.Sprintf("autoUpgradeChildJsonObj : %v", autoUpgradeJSONObj))
	}

	var internalVersionList []string
	startFlag := false
	finishFlag := false

	if !stringContains(InternalUpgradePath, CurrentVersion) {
		startFlag = true
	}

	for _, tempVersion := range InternalUpgradePath {
		if tempVersion == CurrentVersion {
			startFlag = true
			continue
		}
		if tempVersion == version {
			finishFlag = true
		}
		if startFlag && !finishFlag {
			internalVersionList = append(internalVersionList, tempVersion)
		}
		if finishFlag {
			break
		}
	}

	if len(internalVersionList) == 0 {
		cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, fmt.Sprintf("No internal version found till %s CDF.", version))
		return
	}

	for _, internalVersion := range internalVersionList {
		if Debug {
			log.Println(fmt.Sprintf("internalVersion : %s", internalVersion))
		}
		InternalVersionMap[internalVersion] = version

		runes := []rune(internalVersion)
		internalDotVersion := string(runes[:4]) + "." + string(runes[4:])
		cmdPath := filepath.Join("packages", internalDotVersion)
		if Debug {
			log.Println(fmt.Sprintf("versionWithDot : %s", internalDotVersion))
			log.Println(fmt.Sprintf("cmdPath : %s", cmdPath))
		}

		var releaseJSONObj cdfJson.Release
		releaseJSONObj, err = cdfJson.GetReleaseJsonObj(autoUpgradeJSONObj, internalVersion)
		if err != nil {
			return
		}
		if Debug {
			log.Println(fmt.Sprintf("releaseJSONObj : %v", releaseJSONObj))
		}

		steps := releaseJSONObj.Steps

		for _, step := range steps {

			mode := step.Action
			msg := step.Description
			order := step.Order
			cmd := step.Command

			err = stepExec(mode, msg, upgradeProcess, internalVersion, filepath.Join(WorkDir, VersionPackMap[InternalVersionMap[internalVersion]], cmdPath, cmd), order)
			if err != nil {
				break
			}
		}

	}
	return
}

func dynamicUpgradeProcess(version string) (err error) {
	if Debug {
		log.Println(fmt.Sprintf("\n---------------Walking in dynamicUpgradeProcess: %s ---------------", version))
	}

	jsonPath := filepath.Join(CurrentDir, cdfCommon.AutoUpgradeJSON)
	autoUpgradeJSONObj, err := cdfJson.GetAutoUpgradeJsonObj(jsonPath)
	if err != nil {
		return
	}

	if Debug {
		log.Println(fmt.Sprintf("autoUpgradeJSONObj : %v", autoUpgradeJSONObj))
	}

	if Debug {
		log.Println(fmt.Sprintf("version : %s", version))
	}

	var releaseJSONObj cdfJson.Release
	releaseJSONObj, err = cdfJson.GetReleaseJsonObj(autoUpgradeJSONObj, version)
	if err != nil {
		return
	}
	if Debug {
		log.Println(fmt.Sprintf("releaseJSONObj : %v", releaseJSONObj))
	}

	steps := releaseJSONObj.Steps

	for _, step := range steps {

		mode := step.Action
		msg := step.Description
		order := step.Order
		cmd := step.Command

		err = stepExec(mode, msg, upgradeProcess, version, filepath.Join(WorkDir, VersionPackMap[version], cmd), order)
		if err != nil {
			break
		}
	}

	return
}

func upgradeProcess(args ...string) (err error) {
	var mode, version, cmd, order string

	mode = args[0]
	version = args[1]
	cmd = args[2]
	order = args[3]
	cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, fmt.Sprintf("mode: %s", mode))
	cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, fmt.Sprintf("version: %s", version))
	cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, fmt.Sprintf("cmd: %s", cmd))
	cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, fmt.Sprintf("order: %s", order))

	if Debug {
		log.Println(fmt.Sprintf("mode : %s", mode))
		log.Println(fmt.Sprintf("version : %s", version))
		log.Println(fmt.Sprintf("order : %s", order))
		log.Println(fmt.Sprintf("cmd : %s", cmd))
	}

	var nodes []string
	nodes, err = getExecNode(mode, version, strconv.Itoa(UpgExecCall))

	if err != nil {
		return
	} else if len(nodes) == 0 {
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("Nothing remains in step %d", UpgExecCall))
		return
	}

	cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("NodesList : %v", nodes))

	for _, node := range nodes {
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("Starting upgrade process on %s...", node))
		execCmd := filepath.ToSlash(cmd)
		//cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("origin cmd: %s", cmd))
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("exec cmd: %s", execCmd))

		if !DryRun {
			err = cdfSSH.SSHExecCmd(node, SysUser, KeyPath, PassWord, Port, execCmd, true)
			if err != nil {
				return
			}
		}

		err = recordNode(node, version, strconv.Itoa(UpgExecCall))

	}

	return
}

func copyUpgradePacksToCluster(args ...string) (err error) {
	var mode, version string

	mode = args[0]
	version = args[1]
	cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, fmt.Sprintf("mode: %s", mode))
	cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, fmt.Sprintf("version: %s", version))

	var files []string
	var folders []string

	files, folders, err = cdfOS.FilePathWalk(VersionPathMap[version])
	if err != nil {
		return
	}

	parentDir := cdfOS.ParentDir(CurrentDir)
	filePermissionMap := make(map[string]os.FileMode)
	folderPermissionMap := make(map[string]os.FileMode)

	for _, file := range files {
		info, _ := os.Stat(file)
		if Debug {
			log.Println("File : " + file)
			log.Println(fmt.Sprintf("permission : %o", info.Mode().Perm()))
			log.Println("")
		}
		filePermissionMap[file] = info.Mode().Perm()
	}

	for _, folder := range folders {
		info, _ := os.Stat(folder)
		if Debug {
			log.Println("Folder : " + folder)
			log.Println(fmt.Sprintf("permission : %o", info.Mode().Perm()))
			log.Println("")
		}
		//baseFolder := strings.TrimPrefix(folder, parentDir)
		//targetFolder := filepath.Join(WorkDir, baseFolder)
		//targetFolder = filepath.ToSlash(targetFolder)
		folderPermissionMap[folder] = info.Mode().Perm()
	}

	if Debug {
		log.Println(fmt.Sprintf("filePermissionMap : %v", filePermissionMap))
		log.Println(fmt.Sprintf("folderPermissionMap : %v", folderPermissionMap))
	}

	var nodes []string
	nodes, err = getExecNode(mode, version, strconv.Itoa(UpgExecCall))
	if err != nil {
		return
	} else if len(nodes) == 0 {
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("Nothing remains in step %d", UpgExecCall))
		return
	}

	if DryRun {
		for _, node := range nodes {
			err = recordNode(node, version, strconv.Itoa(UpgExecCall))
		}
		return
	}

	ch := make(chan cdfCommon.CopyStatus, len(nodes))

	for _, node := range nodes {
		go func(node string) {
			cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("Creating work directory on node %s ...", node))
			var cmd string
			var conn *ssh.Client
			if err == nil {
				conn, err = cdfSSH.CreatSSHClient(node, SysUser, KeyPath, PassWord, Port)
			}

			if err == nil {
				//cmd = fmt.Sprintf("rm -rf %s/", filepath.ToSlash(WorkDir))
				//cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, node+" : "+cmd)
				//err = cdfSSH.SSHExecCmd(node, SysUser, KeyPath, PassWord, Port, cmd, true)
				err = cdfSSH.RemoveRemoteFolder(conn, filepath.ToSlash(WorkDir))
			}

			if err == nil {
				//cmd = fmt.Sprintf("mkdir -p %s/", filepath.ToSlash(WorkDir))
				//cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, node+" : "+cmd)
				//err = cdfSSH.SSHExecCmd(node, SysUser, KeyPath, PassWord, Port, cmd, true)
				err = cdfSSH.CreateRemoteFolder(conn, filepath.ToSlash(WorkDir))
			}

			if err == nil {
				cmd = fmt.Sprintf("chown %s:%s %s/", SysUser, SysGroup, filepath.ToSlash(WorkDir))
				cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, node+" : "+cmd)
				err = cdfSSH.SSHExecCmd(node, SysUser, KeyPath, PassWord, Port, cmd, true)
			}

			cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("Copying upgrade package to %s ...", node))
			if err == nil {
				//copy files with perm
				for _, srcFile := range files {
					baseFile := strings.TrimPrefix(srcFile, parentDir)
					targetFile := filepath.ToSlash(filepath.Join(WorkDir, baseFile))
					err = cdfSSH.CopyFileLocal2Remote(conn, srcFile, targetFile, filePermissionMap[srcFile])
					if err != nil {
						break
					}
				}
				//recover folders perms
				for _, srcFolder := range folders {
					baseFolder := strings.TrimPrefix(srcFolder, parentDir)
					targetFolder := filepath.ToSlash(filepath.Join(WorkDir, baseFolder))
					err = cdfSSH.RestoreFolderPerm(conn, targetFolder, folderPermissionMap[srcFolder])
				}
			}

			cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("Applying upgrade package permssions on Node %s ...", node))
			if cdfCommon.SysType == "windows" && err == nil {
				path := filepath.ToSlash(filepath.Join(WorkDir, VersionPackMap[version]))
				cmd = fmt.Sprintf("cd %s ; setfacl --restore=%s", path, cdfCommon.ACLPROPERTIES)
				cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, node+" : "+cmd)
				err = cdfSSH.SSHExecCmd(node, SysUser, KeyPath, PassWord, Port, cmd, true)
			}

			if err == nil {
				cmd = fmt.Sprintf("chown %s:%s %s/", SysUser, SysGroup, filepath.ToSlash(WorkDir))
				cdfLog.WriteLog(Logger, cdfCommon.DEBUG, LogLevel, node+" : "+cmd)
				err = cdfSSH.SSHExecCmd(node, SysUser, KeyPath, PassWord, Port, cmd, true)
			}

			if err == nil {
				ch <- cdfCommon.CopyStatus{true, node, fmt.Sprintf("Node: %s process completed.", node)}
			} else {
				ch <- cdfCommon.CopyStatus{false, node, fmt.Sprintf("Node: %s process Failed.", node)}
			}
		}(node)
	}

	i := 0
	for result := range ch {
		if result.Copied {
			cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, result.Description)
			err = recordNode(result.Node, version, strconv.Itoa(UpgExecCall))
		} else {
			cdfLog.WriteLog(Logger, cdfCommon.ERROR, LogLevel, result.Description)
			err = fmt.Errorf("Failed to create auto-upgrade workspace inside all cluster nodes")
		}
		i++
		if i == len(nodes) {
			close(ch)
		}
	}

	return
}

func recordNode(node string, version string, step string) (err error) {
	var exist bool
	path := filepath.Join(TempFolder, version, step)
	exist, err = cdfOS.PathExists(path)
	if err != nil {
		return
	}
	if !exist {
		var f *os.File
		f, err = cdfOS.CreateFile(path)
		defer f.Close()
		if err != nil {
			return
		}
		nodeRecordMap := make(map[string]string)
		nodeRecordMap[node] = "done"
		var data []byte
		data, err = json.Marshal(nodeRecordMap)
		if err != nil {
			return
		}
		err = cdfOS.WriteFile(path, string(data))
		return
	} else {
		var jsonString string
		nodeRecordMap := make(map[string]string)
		jsonString, err = cdfOS.ReadFile(path)
		if err != nil {
			return
		}
		err = json.Unmarshal([]byte(jsonString), &nodeRecordMap)
		if err != nil {
			return
		}
		nodeRecordMap[node] = "done"
		var data []byte
		data, err = json.Marshal(nodeRecordMap)
		err = cdfOS.WriteFile(path, string(data))
		return
	}
	return
}

func readNodeRecord(version string, step string) (result map[string]string, err error) {
	var exist bool
	exist, err = cdfOS.PathExists(filepath.Join(TempFolder, version, step))
	if err != nil {
		return
	}
	if !exist {
		return
	} else {
		var jsonString string
		jsonString, err = cdfOS.ReadFile(filepath.Join(TempFolder, version, step))
		if err != nil {
			return
		}
		err = json.Unmarshal([]byte(jsonString), &result)
		return
	}
	return
}

func checkNodeRecord(nodeRecordMap map[string]string, node string) (exist bool) {
	_, ok := nodeRecordMap[node]
	if ok {
		return true
	}
	return false
}

func getExecNode(mode string, version string, step string) (nodes []string, err error) {
	nodeList := transferMode(mode)

	var nodeRecordMap map[string]string

	nodeRecordMap, err = readNodeRecord(version, step)
	if err != nil {
		return
	}
	if Debug {
		log.Println(nodeRecordMap)
	}

	for _, nodeObj := range nodeList.List {
		if !checkNodeRecord(nodeRecordMap, nodeObj.Name) {
			nodes = append(nodes, nodeObj.Name)
		}
	}
	if Debug {
		log.Println(nodes)
		log.Println(len(nodes))
	}

	return
}

func stringContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func cleanWorkDirInsideCluster(version string) (err error) {
	message := fmt.Sprintf("Clean autoUpgrade workspace inside all cluster nodes..")
	err = stepExec(cdfCommon.AllNodes, message, cleanWorkDir, version, "", "")
	return
}

func cleanWorkDir(args ...string) (err error) {
	mode := args[0]
	version := args[1]

	var nodes []string
	nodes, err = getExecNode(mode, version, strconv.Itoa(UpgExecCall))
	if err != nil {
		return
	} else if len(nodes) == 0 {
		cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, fmt.Sprintf("Nothing remains in step %d", UpgExecCall))
		return
	}

	ch := make(chan cdfCommon.CleanStatus, len(nodes))

	for _, node := range nodes {
		go func(node string) {
			cmd := fmt.Sprintf("rm -rf %s", WorkDir)
			err = cdfSSH.SSHExecCmd(node, SysUser, KeyPath, PassWord, Port, cmd, true)
			if err == nil {
				ch <- cdfCommon.CleanStatus{true, node, fmt.Sprintf("Node: %s process completed.", node)}
			} else {
				ch <- cdfCommon.CleanStatus{false, node, fmt.Sprintf("Node: %s process Failed.", node)}
			}
		}(node)
	}

	i := 0
	for result := range ch {
		if result.Cleaned {
			cdfLog.WriteLog(Logger, cdfCommon.INFO, LogLevel, result.Description)
			err = recordNode(result.Node, version, strconv.Itoa(UpgExecCall))
		} else {
			cdfLog.WriteLog(Logger, cdfCommon.ERROR, LogLevel, result.Description)
			err = fmt.Errorf("Failed to delete auto-upgrade workspace inside all cluster nodes")
		}
		i++
		if i == len(nodes) {
			close(ch)
		}
	}

	return
}
