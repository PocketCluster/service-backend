#!/usr/bin/env bash

export CONFIG_PATH=${PWD}/DEVASSETS/config-dev.yaml
echo "Read ${CONFIG_PATH} ..."
echo "Run upgrade"
go run ${PWD}/exec/upgrader/main.go
