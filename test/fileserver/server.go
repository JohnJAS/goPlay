package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ServePath = "/var/joseph"
)

func main() {
	s := gin.Default()
	//only server the files
	s.Static("/dev0", ServePath)
	//due to use the fs on server, it can get file list by using curl localhost:8099/dev/
	s.StaticFS("/dev", http.Dir(ServePath))
	s.Run(":8099")
}
