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

func RemoveRemoteFile(conn *ssh.Client, path string) (err error) {
	// create new SFTP client
	var c *sftp.Client
	c, err = sftp.NewClient(conn)
	if err != nil {
		return
	}
	defer c.Close()

	err = c.Remove(path)

	return
}

func RemoteFolderExist(conn *ssh.Client, path string) (exist bool, err error) {
	// create new SFTP client
	var c *sftp.Client
	c, err = sftp.NewClient(conn)
	if err != nil {
		return
	}
	defer c.Close()

	_, err = c.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func RemoveRemoteFolder(conn *ssh.Client, path string) (err error) {
	// create new SFTP client
	var c *sftp.Client
	c, err = sftp.NewClient(conn)
	if err != nil {
		return
	}
	defer c.Close()

	var exist bool
	exist, err = RemoteFolderExist(conn, path)
	if exist == false {
		return
	}

	var files []os.FileInfo
	files, err = c.ReadDir(path)
	if err != nil {
		return
	}

	defer c.RemoveDirectory(path)
	for _, info := range files {
		subPath := filepath.ToSlash(filepath.Join(path, info.Name()))
		if info.IsDir() {
			RemoveRemoteFolder(conn, subPath)
		} else {
			RemoveRemoteFile(conn, subPath)
		}
	}

	return
}

func CreateRemoteFolder(conn *ssh.Client, path string) (err error) {
	// create new SFTP client
	var c *sftp.Client
	c, err = sftp.NewClient(conn)
	if err != nil {
		return
	}
	defer c.Close()

	err = c.Mkdir(path)

	return
}

func ChownRemote(conn *ssh.Client, path string, uid, gid int) (err error) {
	// create new SFTP client
	var c *sftp.Client
	c, err = sftp.NewClient(conn)
	if err != nil {
		return
	}
	defer c.Close()

	err = c.Chown(path, uid, gid)

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
