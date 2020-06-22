#!/bin/bash -e

GATE_API_BRANCH=release-1.20.x
BUILD_OS=linux
BUILD_ARCH=amd64


while getopts "o:a:g:" opt; do
  case ${opt} in
    o) #Build OS
      BUILD_OS=${OPTARG}
      ;;
    a) #Build arch
      BUILD_ARCH=${OPTARG}
      ;;
    g) #Gate version
      GATE_API_BRANCH=${OPTARG}
      ;;
  esac
done

echo "Generate checksum"
cd /go/src/github.com/codilime/floodgate/
cp floodgate floodgate-$GATE_API_BRANCH.$BUILD_OS.$BUILD_ARCH
sha1sum floodgate-$GATE_API_BRANCH.$BUILD_OS.$BUILD_ARCH > floodgate-$GATE_API_BRANCH.$BUILD_OS.$BUILD_ARCH.sha1sum

