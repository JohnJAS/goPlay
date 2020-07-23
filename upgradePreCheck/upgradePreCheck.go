package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"

	cdflog "github.com/JohnJAS/goPlay/pkg/log"
)

//program env
var currentDir string
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
}

func main() {
	os.Args = append(os.Args, "-f")
	os.Args = append(os.Args, "202005")
	os.Args = append(os.Args, "-t")
	os.Args = append(os.Args, "202011")

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

func preCheck(c *cli.Context) (err error) {
	return
}
