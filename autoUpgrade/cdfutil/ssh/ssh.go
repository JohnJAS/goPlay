package ssh

import (
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

func getKeyFile(keyPath string) (key ssh.Signer, err error) {
	if keyPath == "" {
		usr, _ := user.Current()
		keyPath = filepath.Join(usr.HomeDir, ".ssh", "id_rsa")
	}

	var buf []byte
	buf, err = ioutil.ReadFile(keyPath)
	if err != nil {
		return
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return
	}
	return
}

func CheckConnection(node string, userName string, keyPath string) (err error) {
	// Get rsa key
	var key ssh.Signer
	key, err = getKeyFile(keyPath)
	if err != nil {
		return
	}
	// Define the Client Config as :
	config := &ssh.ClientConfig{
		User:            userName,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	host := node
	port := "22"
	var client *ssh.Client
	client, err = ssh.Dial("tcp", net.JoinHostPort(host, port), config)
	if err != nil {
		return
	}

	var session *ssh.Session
	session, err = client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	return
}

func CreatSSHClient(node string, userName string, keyPath string) (client *ssh.Client, err error) {
	// Get rsa key
	var key ssh.Signer
	key, err = getKeyFile(keyPath)
	if err != nil {
		return
	}
	// Define the Client Config as :
	config := &ssh.ClientConfig{
		User:            userName,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	host := node
	port := "22"

	return ssh.Dial("tcp", net.JoinHostPort(host, port), config)

}

func SSHExecCmd(node string, userName string, keyPath string, cmd string) (err error) {
	// Get rsa key
	var key ssh.Signer
	key, err = getKeyFile(keyPath)
	if err != nil {
		return
	}
	// Define the Client Config as :
	config := &ssh.ClientConfig{
		User:            userName,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	host := node
	port := "22"

	client, err := ssh.Dial("tcp", net.JoinHostPort(host, port), config)

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

	return

}
