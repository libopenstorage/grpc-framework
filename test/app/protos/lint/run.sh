#!/bin/bash

fail() {
    echo "$1"
    exit 1
}

if [ -z "$SCRIPTSDIR" ] ; then
    fail "SCRIPTSDIR is not defined. This file expects to be called from the Makefile"
fi

for script in ${SCRIPTSDIR}/*.sh ; do
    if [ -f "${script}" ] ; then
        short=$(echo ${script} | sed -e "s#^.*lint#lint#")
        /bin/bash ${script} || fail "*** FAILED to run ${short}"
    fi
done
