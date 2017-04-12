#!/bin/sh
if [ $# -eq 0 ]; then
  proxycheck-client -in=${PROXYLIST} -url=${URL} -apiurl=${APIURL}/api/v1/addproxy -treds=${TREDS}
else
  exec "$@"
fi
