#!/bin/bash

HOST=localhost:8080
USER=kodi
PASS=kodi

RES=$(
curl -s \
  -H 'content-type: application/json;' \
  --data-binary '{ "jsonrpc": "2.0", "method": "VideoLibrary.Clean", "id": "mybash"}' \
  http://$USER:$PASS@$HOST/jsonrpc
)

if [ ! -z "$DEBUG" ]; then
  echo "got: $RES"
fi

OK=$(jq -r '.result' <<< "$RES")

if [ "$OK" = "OK" ]; then
  echo "OK"
else
  echo "error: $RES"
fi
