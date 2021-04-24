package main

import (
	"github.com/rs/zerolog"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"
)

//TempFolder is autoUpgrade temp folder including re-run mark and auto upgrade log
var tempFolder string

var logFilePath string

var logfile *os.File

var upgradelog *zerolog.Logger

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

	//create log file
	path := filepath.Join(tempFolder, "upgrade", "upgrade-"+time.Now().UTC().Format("20060102150405")+".log")
	os.MkdirAll(filepath.Dir(path), 0644)
	logfile, _ = os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logFileWriter := zerolog.ConsoleWriter{Out: logfile, TimeFormat: time.RFC3339}

	multi := zerolog.MultiLevelWriter(consoleWriter, logFileWriter)

	logger := zerolog.New(multi).With().Timestamp().Logger()

	upgradelog = &logger
}

func main() {
	upgradelog.Info().Msgf("upgrade info")
	upgradelog.Debug().Msgf("upgrade debug")
}
