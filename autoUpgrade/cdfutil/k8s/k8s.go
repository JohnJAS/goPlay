package k8s

import (
	cdfSSH "autoUpgrade/cdfutil/ssh"
	"golang.org/x/crypto/ssh"
	"os"
)

//GetCurrentVersion get CDF current version
func GetCurrentVersion(node string, userName string, keyPath string) (err error) {
	cmd := "kubectl get cm base-configmap1 -n core -o json"
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

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

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

	return err
}
