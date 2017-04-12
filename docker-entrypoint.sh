#!/bin/sh
if [ $# -eq 0 ]; then
  proxycheck-client -in=${PROXYLIST} -url=${URL} -apiurl=${APIURL} -treds=${TREDS}
else
  exec "$@"
fi
