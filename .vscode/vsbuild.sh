#!/bin/bash
##
## Build mfw discoverd
##
set -e
TARGET=$1
PORT=22
LOCAL_MUSL_BUILD=false
BACKUP=false

while getopts "t:p:m:b:" flag; do
    case "${flag}" in
        t) TARGET=${OPTARG} ;;
        p) PORT=${OPTARG} ;;
        m) LOCAL_MUSL_BUILD=${OPTARG} ;;
        b) BACKUP=${OPTARG}
    esac
done
shift $((OPTIND-1))

echo "Sending package to $TARGET with port: $PORT and local musl build: $LOCAL_MUSL_BUILD and backup: $BACKUP"

if [ "$LOCAL_MUSL_BUILD" = true ]
then
    docker-compose -f build/docker-compose.build.yml up --exit-code-from musl-local --build musl-local
    if [ $? -ne 0 ]
    then 
        echo "Build failed, aborting"
        exit -1
    fi
    docker-compose -f build/docker-compose.build.yml up --exit-code-from musl-lint --build musl-lint
    if [ $? -ne 0 ]
    then 
        echo "Lint failed, aborting"
        exit -1
    fi
else
    docker-compose -f build/docker-compose.build.yml up  --exit-code-from musl --build musl
    if [ $? -ne 0 ]
    then 
        echo "Build failed, aborting"
        exit -1
    fi
fi

ssh-copy-id -p $PORT root@$TARGET

if [ "$BACKUP" = true ]
then
    now=`date +"%N"`
    mkdir "discoverd_backup"
    scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -P $PORT root@$TARGET:/usr/bin/discoverd ./discoverd_backup/discoverd_${now}; 
fi


ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p $PORT root@$TARGET "/etc/init.d/discoverd stop"; 
sleep 5
scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -P $PORT ./cmd/discoverd/discoverd root@$TARGET:/usr/bin/; 
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p $PORT root@$TARGET "/etc/init.d/discoverd start"
