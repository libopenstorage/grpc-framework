#!/bin/bash

WEBROOT=${1:-../docs}

count=0
for f in $(find ${WEBROOT} -name '*.swagger.json') ; do
    path=$(dirname $f | sed -e "s#${WEBROOT}/##" -e "s#/\$##")
    filepath=${path}/$(basename $f)

    if [ $count -eq 0 ] ; then
        URLS="[{ name: \"${filepath}\", url: \"../${filepath}\" }"
    else
        URLS=$URLS",{ name: \"${filepath}\", url: \"../${filepath}\" }"
    fi
    ((count++))
done
URLS=$URLS"]"

if [ $(uname -s) == "Linux" ] ; then
    sed -i "s#url:.*\$#urls: ${URLS},#" ${WEBROOT}/swagger/swagger-initializer.js
elif [ $(uname -s) == "Darwin" ] ; then
    sed -i '' -e "s#url:.*\$#urls: ${URLS},#" ${WEBROOT}/swagger/swagger-initializer.js
fi
