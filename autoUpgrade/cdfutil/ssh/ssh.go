package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
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

func CheckConnection(node string, userName string, keyPath string) (err error) {
	// Get rsa key
	var key ssh.Signer
	key, err = getKeyFile(keyPath)
	if err != nil {
		return
	}
	// Define the Client Config as :
	config := &ssh.ClientConfig{
		User: userName,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	host := node
	port := "22"
	address := host+":"+port
	var client *ssh.Client
	client, err = ssh.Dial("tcp", address, config)
	if err != nil {
		return
	}

	var session *ssh.Session
	session, err = client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		return errors.New("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())

	return
}
