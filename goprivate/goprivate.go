package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"

	"joseph.com/goprivate/tools"

	privaterepo "github.com/JohnJAS/test"
)

func main(){
	var conn *ssh.Client
	fmt.Println(conn)
	tools.POP()
	privaterepo.Fun()
}