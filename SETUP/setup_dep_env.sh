#!/usr/bin/env bash

echo "(2017/08/25) this is not an impending issue. We'll leave it here"
#source vendor_cleanup.sh

if [[ -z ${GOPATH} || ${GOPATH} != "/opt/gopkg" ]]; then
    echo "GOPATH should point /opt/gopkg"
fi

for d in $(find /opt/gopkg/src/ -maxdepth 3 -type d)
do
    if [ -d "$d/vendor" ]; then
        pushd ${PWD}
        # ls -lat "$d/vendor"
        echo "$d/vendor"
        popd
    fi
done

