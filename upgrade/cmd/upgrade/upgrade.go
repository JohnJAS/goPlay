package main

import (
	"github.com/rs/zerolog"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"
	cdflog "upgrade/pkg/log"
)

var tempFolder string

var logfilePath string

var logfile *os.File

var upgradeLog *zerolog.Logger

//init logger
func init() {
	//identify system OS
	if runtime.GOOS == "windows" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal("Failed to get current user info, initailization failed.")
		}
		tempFolder = filepath.Join(usr.HomeDir, "tmp")
	} else {
		tempFolder = "/tmp"
	}

	//identify logfilePath
	logfilePath = filepath.Join(tempFolder, "upgrade", "upgrade-"+time.Now().UTC().Format("20060102150405")+".log")
}

func main() {
	upgradeLog = startLog(logfilePath)
	defer logfile.Close()

	//testlog
	upgradeLog.Info().Msgf("upgrade info")
	upgradeLog.Debug().Msgf("upgrade debug")
}

func startLog(path string) *zerolog.Logger {
	os.MkdirAll(filepath.Dir(path), 0644)
	logfile, _ = os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	logger := cdflog.NewZeroLog(logfile, 0)
	return &logger
}
