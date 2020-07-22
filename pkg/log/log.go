package log

import (
	"log"
	"strings"
	"time"

	cdfCommon "github.com/goPlay/pkg/common"
)

func getLevel(level int) string {
	switch level {
	case cdfCommon.DEBUG:
		return "DEBUG"
	case cdfCommon.INFO:
		return "INFO"
	case cdfCommon.WARN:
		return "WARN"
	case cdfCommon.ERROR:
		return "ERROR"
	case cdfCommon.FATAL:
		return "FATAL"
	default:
		return ""
	}
}

func getLogLevel(loglevel int) int {
	switch loglevel {
	case cdfCommon.DEBUG:
		return cdfCommon.DEBUG | cdfCommon.INFO | cdfCommon.WARN | cdfCommon.ERROR | cdfCommon.FATAL
	case cdfCommon.INFO:
		return cdfCommon.INFO | cdfCommon.WARN | cdfCommon.ERROR | cdfCommon.FATAL
	case cdfCommon.WARN:
		return cdfCommon.WARN | cdfCommon.ERROR | cdfCommon.FATAL
	case cdfCommon.ERROR:
		return cdfCommon.ERROR | cdfCommon.FATAL
	case cdfCommon.FATAL:
		return cdfCommon.FATAL
	default:
		return 0
	}
}

func TransferLogLevel(logLevel string) int {
	switch logLevel {
	case "DEBUG":
		return cdfCommon.DEBUG
	case "INFO":
		return cdfCommon.DEBUG
	case "WARN":
		return cdfCommon.WARN
	case "ERROR":
		return cdfCommon.ERROR
	case "FATAL":
		return cdfCommon.FATAL
	default:
		return 0
	}
}

//level    : definite by program
//loglevel : definite by user
func WriteLog(logger *log.Logger, level int, loglevel int, msg string, filePath ...string) {
	timeStamp := time.Now().UTC().Format(time.RFC3339Nano)
	log.SetFlags(0)
	logger.SetPrefix(timeStamp + " " + getLevel(level) + " ")

	loglevel = getLogLevel(loglevel)

	switch level {
	case cdfCommon.DEBUG:
		if level & loglevel == level {
			logger.Println(msg)
		}
	case cdfCommon.FATAL:
		if level & loglevel == level {
			log.Println(msg)
			log.Println("The log file is " + strings.Join(filePath, ""))
			logger.Println(msg)
			logger.SetPrefix(timeStamp + " " + getLevel(level) + " ")
			logger.Println("Please refer to the Troubleshooting Guide for help on how to resolve this error.  The log file is " + strings.Join(filePath, ""))
		}
	default:
		if level & loglevel == level {
			log.Println(msg)
			logger.Println(msg)
		}
	}
}
