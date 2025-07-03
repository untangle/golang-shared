#!/bin/bash

libc=$1

user_uid=`id -u` user_gid=`id -g` docker-compose -f build/docker-compose.build.yml up --build ${libc}-local