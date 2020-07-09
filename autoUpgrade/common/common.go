package common

import (
	"os"
	"runtime"
)

//SysType is the value of windows or linux or others
const SysType = runtime.GOOS

//Log level
const (
	DEBUG = 1 << iota //
	INFO
	WARN
	ERROR
	FATAL
)

//File Name
const (
	VersionTXT           = "version.txt"
	UpgradeSH            = "upgrade.sh"
	AutoUpgradeJSON      = "autoUpgrade.json"
	AutoUpgradeChildJSON = "autoUpgradeChild.json"
	ACLPROPERTIES        = "acl.properties"
	UpgradePreCheckSH    = "scripts" + string(os.PathSeparator) + "upgradePreCheck.sh"
)

//Action Type
const (
	AllMasters   = "allMasters"
	AllWorkers   = "allWorkers"
	SingleMaster = "singleMaster"
	AllNodes     = "allNodes"
)

//Time Format
const (
	TIMESTAMP = "20060102150405"
)

//Node Role
const (
	MASTER = "master"
	WORKER = "worker"
)

type Node struct {
	Name string `json:Name`
	Role string `json:Role`
}

func NewNode(name string, role string) Node {
	return Node{Name: name, Role: role}
}

//Nodes of k8s cluster
type NodeList struct {
	List []Node `json:List`
	Num  int    `json:Num`
}

func NewNodeList(list []Node, num int) NodeList {
	return NodeList{List: list, Num: num}
}

func (nodeList *NodeList) AddNode(node Node) {
	nodeList.List = append(nodeList.List, node)
	nodeList.Num++
}

//ConnectionStatus is used in SSH connection check
type ConnectionStatus struct {
	Connected   bool
	Description string
}

//CopyStatus is used in SSH copy check
type CopyStatus struct {
	Copied      bool
	Node        string
	Description string
}

//ExecStatus is used in SSH copy check
type ExecStatus struct {
	Executed    bool
	Node        string
	Description string
}

//CleanStatus is used in SSH copy check
type CleanStatus struct {
	Cleaned     bool
	Node        string
	Description string
}
