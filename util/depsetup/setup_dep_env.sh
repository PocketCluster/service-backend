#!/usr/bin/env bash

source ./devenv.sh

CLEAN_PATH=${GOPATH}/src

./bin/python ./condense_dep.py

source vendor_cleanup.sh

# --- cleanup dependencies ---
for d in $(find ${CLEAN_PATH} -maxdepth 3 -type d)
do
    if [[ -d "${d}/vendor" ]]; then
        pushd ${PWD}
        echo "clean ${d}/vendor"
        cd $d/vendor && clean_vendor
        popd
    fi
done
