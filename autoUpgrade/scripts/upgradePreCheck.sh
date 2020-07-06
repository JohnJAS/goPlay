#!/bin/bash
write_log() {
    local level=$1
    local msg=$2
    #format 2018-09-17 16:30:01.772388983+08:00
    local timestamp=$(date --rfc-3339='ns')
    #format 2018-09-17T16:30:01.772388983+08:00
    timestamp=${timestamp:0:10}"T"${timestamp:11}
    case $level in
        debug)
            echo "${timestamp} DEBUG $msg  " >> $LOGFILE ;;
        info)
            echo -e "$msg"
            echo "${timestamp} INFO $msg  " >> $LOGFILE ;;
        error)
            echo -e "$msg"
            echo "${timestamp} ERROR $msg  " >> $LOGFILE ;;
        warn)
            echo -e "$msg"
            echo "${timestamp} WARN $msg  " >> $LOGFILE ;;
        fatal)
            echo -e "$msg"
            echo -e "The upgrade pre-check log file is ${LOGFILE}.\n"
            echo -e "${timestamp} FATAL $msg  \n" >> $LOGFILE
            echo "${timestamp} INFO Please refer to the Troubleshooting Guide for help on how to resolve this error.  " >> $LOGFILE
            echo "                  The upgrade pre-check log file is ${LOGFILE}" >> $LOGFILE
            #kill $scriptPID
            exit 1
            ;;
        *)
            echo "${timestamp} INFO $msg  " >> $LOGFILE ;;
    esac
}

print_help() {
    echo "Description:"
    echo "    CDF Upgrae Pre-Check"
    echo "Usage:"
    echo "    upgradePreCheck.sh [-f|--fromVersion <version>] [-t|--targetVersion <version>] [--slient]"
    echo "Option:"
    echo "    -f, --fromVersion        The from version upgrade will check. Format: 202005"
    echo "    -t, --targetVersion      The traget version upgrade will check. Format: 202008"
    echo "    -s, --slient             Pop out error message only."
    echo "    -b, --byok               BYOK mode pre-check"
    echo "    -h, --help               Help message."
}
#main
source /etc/profile.d/itom-cdf.sh 2>/dev/null || source /etc/profile
unset HTTP_PROXY; unset HTTPS_PROXY; unset http_proxy; unset https_proxy;

if [[ "$BYOK" == "true" ]] ; then
    LOGFILE=${CURRENT_DIR}/upgrade-`date "+%Y%m%d%H%M%S"`.log
else
    if [[ ! -d ${K8S_HOME}/log/scripts/upgradePreCheck ]] ; then
        mkdir -p ${K8S_HOME}/log/scripts/upgradePreCheck
    fi
    LOGFILE=${K8S_HOME}/log/scripts/upgradePreCheck/upgradePreCheck-`date "+%Y%m%d%H%M%S"`.log
fi


while [[ $# -ge 1 ]] ; do
    key="${1}"
    #============ process input parameter and validation =========================
    case ${key} in
        -f|--fromVersion)
            FROM_VERSION="${2}"
            shift
        ;;
        -t|--targetVersion)
            TARGET_VERSION="${2}"
            shift
        ;;
        -s|--slient)
            SLIENT_MODE="true"
        ;;
        -b|--byok)
            BYOK="true"
        ;;
        -h|--help)
            print_help
            exit 0
        ;;
        *)
            print_help
            write_log "error" "Unknown argument '${key}'. Run -h|--help to get script help."
        ;;
    esac
    shift
done

[[ $FROM_VERSION == "" ]] && write_log "error" "-f|--fromVersion is mandatory." 
[[ $TARGET_VERSION == "" ]] && write_log "error" "-t|--targetVersion is mandatory." 

CURRENT_DIR=$(cd `dirname $0`;pwd)
JQ=${CURRENT_DIR}/../bin/jq

UPGRADE_CHAIN=($(cat ${CURRENT_DIR}/../autoUpgrade.json | ${JQ} -r '.[].targetVersion' | sort -h | xargs))
if [[ $SLIENT_MODE != "true" ]] ; then
    write_log "info" "UPGRADE_CHAIN: ${UPGRADE_CHAIN[*]}"
fi
