package common

import (
	"runtime"
)

//SysType is the value of windows or linux or others
const SysType = runtime.GOOS

//Log level
const (
	DEBUG = 1 << iota	//
	INFO
	WARN
	ERROR
	FATAL
)

//Time Format
const (
	TIMESTAMP = "20060102150405"
)

//Nodes of k8s cluster
type Nodes struct {
	NodeList []string
	Num      int
}

//ConnectionStatus is used in SSH connection check
type ConnectionStatus struct {
	Connected   bool
	Description string
}
