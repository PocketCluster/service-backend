#!/usr/bin/env bash

export CONFIG_PATH=${PWD}/DEVASSETS/config-dev.yaml
echo "Read ${CONFIG_PATH} ..."
echo "Check DEV Config..."
go run ${PWD}/exec/config/main.go

echo "\nCheck Product Config..."
export CONFIG_PATH=${PWD}/config.yaml
${PWD}/exec/config/config
