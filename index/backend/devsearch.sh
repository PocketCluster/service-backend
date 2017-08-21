#!/usr/bin/env bash

export CONFIG_PATH=${PWD}/DEVASSETS/config-dev.yaml
echo "Read ${CONFIG_PATH} ..."
echo "Run Search..."
go run ${PWD}/exec/buildsearch/main.go
