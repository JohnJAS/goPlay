package k8s

import (
	"bytes"
	"encoding/json"
	"golang.org/x/crypto/ssh"

	cdfSSH "autoUpgrade/cdfutil/ssh"
)

//GetCurrentVersion get CDF current version
func GetCurrentVersion(node string, userName string, keyPath string) (currentVersion string, outbuf bytes.Buffer, errbuf bytes.Buffer, err error) {
	cmd := "kubectl get cm base-configmap -n core -o json"
	//cmd := "/root/workspace/file.sh"

	client, err := cdfSSH.CreatSSHClient(node, userName, keyPath)
	if err != nil {
		return
	}

	var session *ssh.Session
	session, err = client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	session.Stdout = &outbuf
	session.Stderr = &errbuf

	//cmdReader, err := session.StdoutPipe()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//scanner := bufio.NewScanner(cmdReader)
	//go func() {
	//	for scanner.Scan() {
	//		fmt.Println(scanner.Text())
	//	}
	//}()

	err = session.Run(cmd)

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
