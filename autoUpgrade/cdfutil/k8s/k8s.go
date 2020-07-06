package k8s

import (
	"bytes"
	"encoding/json"
	"strings"

	cdfSSH "autoUpgrade/cdfutil/ssh"
	cdfCommon "autoUpgrade/common"
)

//GetCurrentVersion get CDF current version
func GetCurrentVersion(node string, userName string, keyPath string, password string, port string) (currentVersion string, outbuf bytes.Buffer, errbuf bytes.Buffer, err error) {
	cmd := "kubectl get cm base-configmap -n core -o json"

	outbuf, errbuf, err = cdfSSH.SSHExecCmdReturnResult(node, userName, keyPath, password, port, cmd)
	if err != nil {
		return
	}

	var cmlv1 map[string]json.RawMessage
	err = json.Unmarshal(outbuf.Bytes(), &cmlv1)
	if err != nil {
		return
	}

	var cmlv2 map[string]string
	err = json.Unmarshal(cmlv1["data"], &cmlv2)
	if err != nil {
		return
	}

	currentVersion = cmlv2["PLATFORM_VERSION"]
	return
}

//GetCurrrentNodes get CDF current nodes
func GetCurrrentNodes(nodelist *cdfCommon.NodeList, node string, userName string, keyPath string, password string, port string) (errbuf bytes.Buffer, err error) {
	cmdMaster := "kubectl get nodes -l master=true -o jsonpath='{.items[?(@.kind==\"Node\")].metadata.name}'"
	cmdWorker := "kubectl get nodes -l 'master notin (true)' -o jsonpath='{.items[?(@.kind==\"Node\")].metadata.name}'"

	var outbuf bytes.Buffer
	outbuf, errbuf, err = cdfSSH.SSHExecCmdReturnResult(node, userName, keyPath, password, port, cmdMaster)
	if err != nil {
		return
	}
	for _, node := range strings.Split(outbuf.String(), " ") {
		nodelist.AddNode(cdfCommon.NewNode(node, cdfCommon.MASTER))
	}

	outbuf, errbuf, err = cdfSSH.SSHExecCmdReturnResult(node, userName, keyPath, password, port, cmdWorker)
	if err != nil {
		return
	}
	for _, node := range strings.Split(outbuf.String(), " ") {
		nodelist.AddNode(cdfCommon.NewNode(node, cdfCommon.WORKER))
	}

	return
}
