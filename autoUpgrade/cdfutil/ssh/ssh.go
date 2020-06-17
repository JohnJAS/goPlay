package ssh

import (
	"bytes"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
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

func CreatSSHClient(node string, userName string, keyPath string, port string) (client *ssh.Client, err error) {
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

	return ssh.Dial("tcp", net.JoinHostPort(node, port), config)

}

func CheckConnection(node string, userName string, keyPath string, port string) (err error) {
	var client *ssh.Client
	client, err = CreatSSHClient(node, userName, keyPath, port)
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

func SSHExecCmd(node string, userName string, keyPath string, port string, cmd string) (err error) {
	var client *ssh.Client
	client, err = CreatSSHClient(node, userName, keyPath, port)
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

func SSHExecCmdReturnResult(node string, userName string, keyPath string, port string, cmd string) (outbuf bytes.Buffer, errbuf bytes.Buffer, err error) {
	var client *ssh.Client
	client, err = CreatSSHClient(node, userName, keyPath, port)
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

func CopyFileLocal2Remote(conn *ssh.Client, srcfile string, desfile string) (err error) {

	// create new SFTP client
	var c *sftp.Client
	c, err = sftp.NewClient(conn)
	if err != nil {
		return
	}
	defer c.Close()

	var sf *os.File
	sf, err = os.Open(srcfile)
	if err != nil {
		return
	}
	defer sf.Close()

	err = c.MkdirAll(filepath.ToSlash(filepath.Dir(desfile)))
	if err != nil {
		return
	}

	var df *sftp.File
	df, err = c.Create(desfile)
	if err != nil {
		return
	}
	defer df.Close()

	// Copy the file
	//var r int64
	_, err = df.ReadFrom(sf)
	//log.Println(fmt.Sprintf("Copy file %s : %d",desfile ,r))

	return
}
