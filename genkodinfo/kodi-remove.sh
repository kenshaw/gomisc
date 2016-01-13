#!/bin/bash

HOST=localhost:8080
USER=kodi
PASS=kodi

LABEL="$1"

if [ -z "$LABEL" ]; then
  echo "usage: $0 <label>"
  exit 1
fi

REQ=$(cat <<-ENDREQ
{
  "jsonrpc": "2.0",
  "id": "mybash",
  "method": "VideoLibrary.GetTVShows",
  "params": {
    "filter": {
      "field": "title",
      "operator": "contains",
      "value": "${LABEL}"
    }
  }
}
ENDREQ
)

RES=$(
curl -s \
  -H 'content-type: application/json;' \
  --data-binary "$REQ" \
  http://$USER:$PASS@$HOST/jsonrpc
)

if [ ! -z "$DEBUG" ]; then
  echo "$RES"
fi

if [ "$2" = "--clean=yes" ]; then
  TVSHOWID=$(jq -r '.result.tvshows[0].tvshowid' <<< "$RES")

  curl -s \
    -H 'content-type: application/json;' \
    --data-binary '{ "jsonrpc": "2.0", "method": "VideoLibrary.RemoveTVShow", "params": { "tvshowid": '$TVSHOWID' }, "id": "mybash"}' \
    http://$USER:$PASS@$HOST/jsonrpc |jq '.'
else
  jq '.result' <<< "$RES"
fi
