package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	cdfCommon "github.com/JohnJAS/goPlay/pkg/common"
	cdfLog "github.com/JohnJAS/goPlay/pkg/log"
	cdfOS "github.com/JohnJAS/goPlay/pkg/os"
	"github.com/urfave/cli/v2"
)

//program env
var currentDir string
var logPath string
var logger *log.Logger

var k8sHome string

//decouple with autoUpgrade json, currently I think it's better
var upgradeChain = []string{
	"000000",
	"201811",
	"201902", "201905", "201908", "201911",
	"202002", "202005", "202008", "202011",
	"202102", "202105", "202108", "202111",
	"202202", "202205", "202208", "202211",
	"999999",
}

//cli params
var fromVersion, targetVersion string
var silentMode, byokMode, debugMode bool

func init() {
	var err error
	//get current dir
	currentDir, err = os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	//test on windows
	os.Setenv(cdfCommon.K8SHome, "C:\\tmp")
	k8sHome = os.Getenv(cdfCommon.K8SHome)
	if k8sHome == "" {
		log.Fatal("failed to get the value of K8S_HOME")
	}

	logPath = filepath.Join(k8sHome, "log", "scripts", "upgradePreCheck-"+time.Now().UTC().Format(cdfCommon.TIMESTAMP)+".log")

}

func main() {
	os.Args = append(os.Args, "-f")
	os.Args = append(os.Args, "202005")
	os.Args = append(os.Args, "-t")
	os.Args = append(os.Args, "202011")

	logger = initLogger(logPath)
	cdfLog.WriteLog(logger, cdfLog.DEBUG, cdfLog.DEBUG, "Current directory : "+currentDir)

	app := &cli.App{
		Name:            "upgradePreCheck",
		Usage:           "Precheck before CDF upgrade",
		Description:     "This program is executed by upgrade automatically and it will check the environment for different cases, including version-less upgrade, classic upgrade, byok upgrade and so on.",
		UsageText:       "upgradePrecheck -f <from_version> -t <target_version>",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "f",
				Aliases:     []string{"fromVersion"},
				Required:    true,
				Destination: &fromVersion,
				Usage:       "Current CDF version.(mandatory)",
			},
			&cli.StringFlag{
				Name:        "t",
				Aliases:     []string{"targetVersion"},
				Required:    true,
				Destination: &targetVersion,
				Usage:       "The target CDF version needed to be upgraded.(mandatory)",
			},
			&cli.BoolFlag{
				Name:        "silent",
				Value:       false,
				Destination: &silentMode,
				Usage:       "Pre-check in silent Mode. Only pop out error message",
			},
			&cli.BoolFlag{
				Name:        "byok",
				Value:       false,
				Destination: &byokMode,
				Usage:       "Pre-check for BYOK upgrade",
			},
			&cli.BoolFlag{
				Name:        "debug",
				Value:       false,
				Destination: &debugMode,
				Usage:       "Debug mode",
			},
		},
		Action: preCheck,
	}
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}

func initLogger(path string) (logger *log.Logger) {
	var file *os.File
	var err error
	file, err = cdfOS.CreateFile(path)
	if err != nil {
		log.Fatal(err)
	}

	//initialize logger
	return log.New(file, "", 0)
}

func preCheck(c *cli.Context) (err error) {

	cdfLog.WriteLog(logger, cdfLog.DEBUG, cdfLog.DEBUG, fmt.Sprintf("fromVersion   : %v", fromVersion))
	cdfLog.WriteLog(logger, cdfLog.DEBUG, cdfLog.DEBUG, fmt.Sprintf("targetVersion : %v", targetVersion))
	cdfLog.WriteLog(logger, cdfLog.DEBUG, cdfLog.DEBUG, fmt.Sprintf("silentMode    : %v", silentMode))
	cdfLog.WriteLog(logger, cdfLog.DEBUG, cdfLog.DEBUG, fmt.Sprintf("byokMode      : %v", byokMode))
	cdfLog.WriteLog(logger, cdfLog.DEBUG, cdfLog.DEBUG, fmt.Sprintf("debugMode     : %v", debugMode))

	return
}
