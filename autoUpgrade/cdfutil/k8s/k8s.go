package k8s

import (
	cdfSSH "autoUpgrade/cdfutil/ssh"
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
)

//GetCurrentVersion get CDF current version
func GetCurrentVersion(node string, userName string, keyPath string) (err error) {
	//cmd := "kubectl get cm base-configmap -n core -o json"
	cmd := "bash -x /root/workspace/file.sh"

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

	cmdReader, err := session.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	err = session.Run(cmd)

	return
}
