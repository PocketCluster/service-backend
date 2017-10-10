#!/usr/bin/env bash

# workspace root
export WORK_ROOT=${PWD}
export GOREPO="/opt/gopkg"
export GODEPS="${WORK_ROOT}/dep"

echo "(2017/08/25) this is not an impending issue. We'll leave it here"
#source vendor_cleanup.sh

if [[ -z ${GOPATH} || ${GOPATH} != "/opt/gopkg" ]]; then
    echo "GOPATH should point /opt/gopkg"
fi

# --- setup github.com/blevesearch/bleve (for log search ekanite) ---
pushd ${WORK_ROOT}
echo "github.com/blevesearch/bleve ..."
BLSEARCH="${GOREPO}/src/github.com/blevesearch"
if [[ ! -d ${BLSEARCH} ]]; then
    mkdir -p ${BLSEARCH}
fi
BLEVE="${BLSEARCH}/bleve"
LINK=$(readlink ${BLEVE})
if [[ ! -d ${BLEVE} ]] || [[ ${LINK} != "../../../dep/bleve" ]]; then
    echo "cleanup old link ${BLEVE} and rebuild ..."
    cd ${BLSEARCH} && (rm ${BLEVE} || true) && ln -s ../../../dep/bleve ./bleve
fi
popd

# --- cleanup dependencies ---
for d in $(find /opt/gopkg/src/ -maxdepth 3 -type d)
do
    if [ -d "$d/vendor" ]; then
        pushd ${PWD}
        # ls -lat "$d/vendor"
        echo "$d/vendor"
        popd
    fi
done

