#!/bin/bash
[ -d ${APP_PATH} ] && rc=0 || rc=1

# Output based on success or failure
if [ $rc -eq 0 ]; then
    echo "{\"changed\": true}"
    exit 0
else
    echo "{\"changed\": false, \"failed\": true, \"msg\": \"release folder doest not exists\"}"
    exit 1
fi
