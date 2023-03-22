#!/bin/bash
nginx -t
rc=$?

# Output based on success or failure
if [ $rc -eq 0 ]; then
    nginx -s reload
    echo "{\"changed\": true}"
    exit 0
else
    echo "{\"changed\": false, \"failed\": true, \"msg\": \"error reload nginx config\"}"
    exit 1
fi
