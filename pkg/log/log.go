package log

import (
	"log"
	"strings"
	"time"
)

//Log level
const (
	DEBUG = 1 << iota
	INFO
	WARN
	ERROR
	FATAL
)

func LogLevelItoa(level int) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return ""
	}
}

func LogLevelAtoi(logLevel string) int {
	switch logLevel {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return DEBUG
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return 0
	}
}

func transferLogLevel(loglevel int) int {
	switch loglevel {
	case DEBUG:
		return DEBUG | INFO | WARN | ERROR | FATAL
	case INFO:
		return INFO | WARN | ERROR | FATAL
	case WARN:
		return WARN | ERROR | FATAL
	case ERROR:
		return ERROR | FATAL
	case FATAL:
		return FATAL
	default:
		return 0
	}
}

//level    : definite by program
//loglevel : definite by user
func WriteLog(logger *log.Logger, level int, loglevel int, msg string, filePath ...string) {
	timeStamp := time.Now().UTC().Format(time.RFC3339Nano)
	log.SetFlags(0)
	logger.SetPrefix(timeStamp + " " + LogLevelItoa(level) + " ")

	loglevel = transferLogLevel(loglevel)

	switch level {
	case DEBUG:
		if level & loglevel == level {
			logger.Println(msg)
		}
	case FATAL:
		if level & loglevel == level {
			log.Println(msg)
			log.Println("The log file is " + strings.Join(filePath, ""))
			logger.Println(msg)
			logger.SetPrefix(timeStamp + " " + LogLevelItoa(level) + " ")
			logger.Println("Please refer to the Troubleshooting Guide for help on how to resolve this error.  The log file is " + strings.Join(filePath, ""))
		}
	default:
		if level & loglevel == level {
			log.Println(msg)
			logger.Println(msg)
		}
	}
}
