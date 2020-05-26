package common

import (
	"runtime"
)

//SysType is the value of windows or linux or others
const SysType = runtime.GOOS

//Log
const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
	FATAL = "FATAL"
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
