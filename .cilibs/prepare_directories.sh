#!/bin/bash -e

BUILD_OS=linux
BUILD_ARCH=amd64

echo "Prepare directories"
sudo mkdir /floodgate
sudo chmod 777 /floodgate
mkdir -p /floodgate/bin
mkdir -p /floodgate/libs
mkdir -p /floodgate/resources
cp -r sponnet /floodgate/libs/
cp -r examples /floodgate/resources/
cp floodgate /floodgate/bin/
chmod +x /floodgate/bin/floodgate

