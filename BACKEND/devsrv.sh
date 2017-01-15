#!/usr/bin/env bash

export CONFIG_PATH=${PWD}/config-dev.yaml
echo "Read ${CONFIG_PATH} ..."
echo "Run server..."
go run ${PWD}/exec/indexsrv/main.go
