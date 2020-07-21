#!/bin/bash -e

BUILD_OS=linux
BUILD_ARCH=amd64

sudo bash -c "source .cilibs/prepare_extra_directories.sh"
mkdir -p /floodgate/libs
mkdir -p /floodgate/resources
cp -r sponnet /floodgate/libs/
cp -r examples /floodgate/resources/
cp floodgate /floodgate/bin/
chmod +x /floodgate/bin/floodgate

