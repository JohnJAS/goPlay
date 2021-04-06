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
	s.StaticFS("/dev", http.Dir(ServePath))
	s.Run(":8099")
}
