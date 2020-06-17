package common

import (
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
	VersionTXT      = "version.txt"
	UpgradeSH       = "upgrade.sh"
	AutoUpgradeJSON = "autoUpgrade.json"
	ACLPROPERTIES   = "acl.properties"
)

//Action Type
const (
	AllMasters   = "AllMasters"
	AllWorkers   = "AllWorkers"
	SingleMaster = "SingleMaster"
	AllNodes     = "AllNodes"
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
	Description string
}
