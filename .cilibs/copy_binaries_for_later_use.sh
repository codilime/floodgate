#!/bin/bash -e

echo "Copy binaries for later use"
mkdir -p /floodgate/bin
chmod 777 /floodgate/bin
cp /go/src/github.com/codilime/floodgate/floodgate /floodgate/bin/floodgate

