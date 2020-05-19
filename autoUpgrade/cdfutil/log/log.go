package log

import (
	"log"
	"time"

	cdfCommon "autoUpgrade/common"
)

func write_log() {
	//local level=$1
	//local msg=$2
	//#format 2018-09-17 16:30:01.772388983+08:00
	//local timestamp=$(date --rfc-3339='ns')
	//#format 2018-09-17T16:30:01.772388983+08:00
	//timestamp=${timestamp:0:10}"T"${timestamp:11}
	//case $level in
	//debug)
	//echo "${timestamp} DEBUG $msg  " >> $LOGFILE
	//;;
	//info)
	//echo -e "$msg"
	//echo "${timestamp} INFO $msg  " >> $LOGFILE
	//;;
	//warn)
	//echo -e "$msg"
	//echo "${timestamp} WARN $msg  " >> $LOGFILE
	//;;
	//error)
	//echo -e "$msg"
	//echo "${timestamp} ERROR $msg  " >> $LOGFILE
	//exit 1
	//;;
	//fatal)
	//echo -e "$msg"
	//echo -e "The CDF autoUpgrade log file is ${LOGFILE}.\n"
	//echo -e "${timestamp} FATAL $msg  \n" >> $LOGFILE
	//echo "${timestamp} INFO Please refer to the Troubleshooting Guide for help on how to resolve this error.  " >> $LOGFILE
	//echo "                                         The CDF autoUpgrade log file is ${LOGFILE}" >> $LOGFILE
	//#kill $scriptPID
	//exit 1
	//;;
	//*)
	//echo "${timestamp} INFO $msg  " >> $LOGFILE ;;
	//esac
}

func WriteLog(logger *log.Logger, level string, msg string) {
	timeStamp := time.Now().UTC().Format(cdfCommon.RFC3339Nano)

	logger.SetPrefix(timeStamp + " " + level + " ")

	logger.Println(msg)
}
