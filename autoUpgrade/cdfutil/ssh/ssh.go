package ssh

import (
	"bytes"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"
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

func CheckConnection(node string, userName string, keyPath string) (err error) {
	var client *ssh.Client
	client, err = CreatSSHClient(node, userName, keyPath)
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

func SSHExecCmd(node string, userName string, keyPath string, cmd string) (err error) {
	var client *ssh.Client
	client, err = CreatSSHClient(node, userName, keyPath)
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

	return

}

func SSHExecCmdReturnResult(node string, userName string, keyPath string, cmd string) (outbuf bytes.Buffer, errbuf bytes.Buffer, err error) {
	var client *ssh.Client
	client, err = CreatSSHClient(node, userName, keyPath)
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

	err = session.Run(cmd)

	return

}

func CopyFile(node string, userName string, keyPath string, srcfile string, desfile string) (err error) {
	var conn *ssh.Client

	conn, err = CreatSSHClient(node, userName, keyPath)
	if err != nil {
		return
	}

	// create new SFTP client
	c, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	fmt.Println(srcfile)
	fmt.Println(desfile)

	s, err := os.Open(srcfile)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	d, err := c.Create(desfile)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	// Copy the file
	var r int64
	r, err = d.ReadFrom(s)
	fmt.Sprintln("Read : %d", r)


	return
}
