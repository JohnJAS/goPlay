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

func CreatSSHClient(node string, userName string, keyPath string, password string, port string) (client *ssh.Client, err error) {
	//if no rsakey or password provided, get rsa key from default place
	if keyPath == "" && password == "" {
		usr, _ := user.Current()
		keyPath = filepath.Join(usr.HomeDir, ".ssh", "id_rsa")
	}

	var config *ssh.ClientConfig
	if password != "" {
		// Define the Client Config as :
		config = &ssh.ClientConfig{
			User:            userName,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Auth:            []ssh.AuthMethod{ssh.Password(password)},
		}
	} else {
		// Get rsa key
		var key ssh.Signer
		key, err = getKeyFile(keyPath)
		if err != nil {
			return
		}
		// Define the Client Config as :
		config = &ssh.ClientConfig{
			User:            userName,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(key),
			},
		}
	}

	return ssh.Dial("tcp", net.JoinHostPort(node, port), config)

}

func CheckConnection(node string, userName string, keyPath string, password string, port string) (err error) {
	var client *ssh.Client
	client, err = CreatSSHClient(node, userName, keyPath, password, port)
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

func SSHExecCmd(node string, userName string, keyPath string, password string, port string, cmd string, output bool) (err error) {
	var client *ssh.Client
	client, err = CreatSSHClient(node, userName, keyPath, password, port)
	defer client.Close()
	if err != nil {
		return
	}

	var session *ssh.Session
	session, err = client.NewSession()
	defer session.Close()
	if err != nil {
		return
	}

	if output {
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
	}

	err = session.Run(cmd)

	return

}

func SSHExecCmdReturnResult(node string, userName string, keyPath string, password string, port string, cmd string) (outbuf bytes.Buffer, errbuf bytes.Buffer, err error) {
	var client *ssh.Client
	client, err = CreatSSHClient(node, userName, keyPath, password, port)
	defer client.Close()
	if err != nil {
		return
	}

	var session *ssh.Session
	session, err = client.NewSession()
	defer session.Close()
	if err != nil {
		return
	}

	session.Stdout = &outbuf
	session.Stderr = &errbuf

	err = session.Run(cmd)

	return

}

func CopyFileLocal2Remote(conn *ssh.Client, srcfile string, desfile string, perm os.FileMode) (err error) {

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
	if err != nil {
		return
	}
	//log.Println(fmt.Sprintf("Copy file %s : %d",desfile ,r))

	err = c.Chmod(desfile, perm)
	if err != nil {
		return
	}

	return
}

func RestoreFolderPerm(conn *ssh.Client, desfolder string, perm os.FileMode) (err error) {

	// create new SFTP client
	var c *sftp.Client
	c, err = sftp.NewClient(conn)
	if err != nil {
		return
	}
	defer c.Close()

	err = c.Chmod(desfolder, perm)
	if err != nil {
		return
	}

	return
}
