#!/usr/bin/env bash

if [[ ! -z ${GOPATH} ]]; then
    unset GOPATH
fi

export GOPATH="${HOME}/INDEX/GOPKG"
export PATH=$GOROOT/bin:$GOPATH/bin:$HOME/.util:$NATIVE_PATH
