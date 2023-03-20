#!/bin/bash
# Run from the top of the source tree
# Requires the test app to be built


./test/app/bin/server &
pid=$!
sleep 3
./test/app/bin/client
ret=$?
kill -9 $!
exit $?
