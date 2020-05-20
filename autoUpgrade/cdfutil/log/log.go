package log

import (
	"log"
	"strings"
	"time"

	cdfCommon "autoUpgrade/common"
)

func WriteLog(logger *log.Logger, level string, msg string, filePath ...string) {
	timeStamp := time.Now().UTC().Format(time.RFC3339Nano)
	log.SetFlags(0)
	logger.SetPrefix(timeStamp + " " + level + " ")

	switch level {

	case cdfCommon.DEBUG:
		logger.Println(msg)
	case cdfCommon.FATAL:
		log.Println(msg)
		log.Println("The CDF autoUpgrade log file is "+strings.Join(filePath,""))
		logger.Println(msg)
		logger.SetPrefix(timeStamp + " " + cdfCommon.INFO + " ")
		logger.Println("Please refer to the Troubleshooting Guide for help on how to resolve this error.  The CDF autoUpgrade log file is "+strings.Join(filePath,""))
	default:
		log.Println(msg)
		logger.Println(msg)
	}
}
