#!/bin/bash

while [ $# -gt 0 ];do
    case "$1" in
    -i|--infra) UPGRADE_INFRA=true
        case "$2" in
            --docker) UPGRADE_INFRA_DOCKER=true
                case "$3" in
                    --k8s) UPGRADE_INFRA_K8S=true
                        case "$4" in
                            --cdf) UPGRADE_INFRA_CDF=true
                                shift 4 ;;
                            *)
                                shift 3 ;;
                        esac ;;
                    --cdf) UPGRADE_INFRA_CDF=true
                        case "$4" in
                            --k8s) UPGRADE_INFRA_K8S=true
                                shift 4 ;;
                            *)
                                shift 3 ;;
                        esac ;;
                    *) 
                        shift 2 ;;
                esac ;;
            --k8s) UPGRADE_INFRA_K8S=true
            case "$3" in
                --docker) UPGRADE_INFRA_DOCKER=true
                    case "$4" in
                        --cdf) UPGRADE_INFRA_CDF=true
                            shift 4 ;;
                        *)
                            shift 3 ;;
                    esac ;;
                --cdf) UPGRADE_INFRA_CDF=true
                    case "$4" in
                        --docker) UPGRADE_INFRA_DOCKER=true
                            shift 4 ;;
                        *)
                            shift 3 ;;
                    esac ;;
                *) 
                    shift 2 ;;
            esac ;;
            --cdf) UPGRADE_INFRA_CDF=true
            case "$3" in
                --k8s) UPGRADE_INFRA_K8S=true
                    case "$4" in
                        --docker) UPGRADE_INFRA_DOCKER=true
                            shift 4 ;;
                        *)
                            shift 3 ;;
                    esac ;;
                --docker) UPGRADE_INFRA_DOCKER=true
                    case "$4" in
                        --k8s) UPGRADE_INFRA_K8S=true
                            shift 4 ;;
                        *)
                            shift 3 ;;
                    esac ;;
                *) 
                    shift 2 ;;
            esac ;;
            *) UPGRADE_INFRA_DOCKER=true; UPGRADE_INFRA_K8S=true; UPGRADE_INFRA_CDF=true;
                shift ;;
        esac ;;
    -u|--upgrade) UPGRADE_CDF=true
        shift;;
    -e|--evict) DRAIN=true
        shift;;
    --drain) DRAIN=true
        shift;;
    --drain-timeout)
        case "$2" in
          -*) echo "--drain-timeout option requires a value." ; exit 1 ;;
          *)  if [[ -z "$2" ]] ; then echo "--drain-timeout option requires a length of time.(second)" ; exit 1 ; fi ; DRAIN_TIMEOUT=$2 ; shift 2 ;;
        esac ;;
    -c|--clean) CLEAN=true
        shift;;
    -y|--yes) FORCE_YES=true
        shift;;
    -t|--temp)
        case "$2" in
          -*) echo "-t|--temp option requires a value." ; exit 1 ;;
          *)  if [[ -z "$2" ]] ; then echo "-t|--temp option needs to provide folder path. " ; exit 1 ; fi ; TEMP_FOLDER=$2 ; shift 2 ;;
        esac ;;
    -dev|--developerMode) DEVELOPOR_MODE=true
        shift;;
    -sio|--skipImageOperation) SKIP_IMAGE_OPERATION=true
        shift;;
    -h|--help)
        exit 1;;
    *) 
        echo -e "The input parameter $1 is not a supported parameter or not used in a correct way. Please refer to the following usage.\n"
        exit 1;;
    esac
done

if [[ $UPGRADE_INFRA == "true" ]] ; then
    echo "Pretend to execute upgrade.sh -i in 2020.02 subfolder..."

elif [[ $UPGRADE_CDF == "true" ]] ; then
    echo "Pretend to execute upgrade.sh -u in 2020.02 subfolder..."

elif [[ $CLEAN == "true" ]] ; then
    echo "Pretend to execute upgrade.sh -c in 2020.02 subfolder..."
fi