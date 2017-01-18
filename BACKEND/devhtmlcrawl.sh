#!/usr/bin/env bash

export CONFIG_PATH=${PWD}/DEVASSETS/config-dev.yaml
echo "Read ${CONFIG_PATH} ..."
echo "Run Scrap..."
go run ${PWD}/exec/htmlscrap/main.go
