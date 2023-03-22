#!/bin/bash
systemctl is-active --quiet ${SYSTEMD_UNIT} && rc=0 || rc=1
rc=$?

# Output based on success or failure
if [ $rc -eq 0 ]; then
    echo "{\"changed\": true}"
    exit 0
else
    echo "{\"changed\": false, \"failed\": true, \"msg\": \"App is not running\"}"
    exit 1
fi
