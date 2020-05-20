package log

import (
	"io/ioutil"
	"log"
	"time"

	cdfCommon "autoUpgrade/common"
)

func WriteLog(logger *log.Logger, level string, msg string) {
	timeStamp := time.Now().UTC().Format(time.RFC3339Nano)
	logger.SetPrefix(timeStamp + " " + level + " ")

	switch timeStamp {

	case cdfCommon.DEBUG:
		log.SetOutput(ioutil.Discard)
		logger.SetOutput(ioutil.Discard)
		logger.Println(msg)
	case cdfCommon.FATAL:
		log.Println(msg)
		log.Println("The CDF autoUpgrade log file is ")
		logger.Println(msg)
		logger.SetPrefix(timeStamp + " " + cdfCommon.INFO + " ")
		logger.Println("Please refer to the Troubleshooting Guide for help on how to resolve this error.  The CDF autoUpgrade log file is ")
	default:
		log.Println(msg)
		logger.Println(msg)

	}
}
