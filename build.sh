#!/bin/bash

libc=$1

docker-compose -f build/docker-compose.build.yml up --build ${libc}-local