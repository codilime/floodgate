#!/bin/bash -e

GATE_VERSION=release-1.20.x
BUILD_OS=linux
BUILD_ARCH=amd64


while getopts "o:a:g" opt; do
  case ${opt} in
    o) #Build OS
      BUILD_OS=${OPTARG}
      ;;
    a) #Build arch
      BUILD_ARCH=${OPTARG}
      ;;
    g) #Gate version
      GATE_VERSION=${OPTARG}
      ;;
  esac
done

echo "Get dependencies"
go mod download
go get -u golang.org/x/lint/golint

echo "Examine source code with go vet"
go vet -v ./...

echo "Execute go tests"
go test -v ./... -coverprofile cover.out

echo "Compile code"
if [ -z "$TRAVIS_BRANCH" ]
then
    export RELEASE=$(echo $CIRCLE_TAG | sed 's/^v[0-9]\+\.[0-9]\+\.[0-9]\+-\?//')
else
    export RELEASE=$CIRCLE_BRANCH
fi
env GOOS=${BUILD_OS} GOARCH=${BUILD_ARCH} go build -ldflags \
"-X github.com/codilime/floodgate/version.GitCommit=$TRAVIS_COMMIT \
-X github.com/codilime/floodgate/version.BuiltDate=$(date  +%Y-%m-%d_%H:%M:%S) \
-X github.com/codilime/floodgate/version.Release=$RELEASE \
-X github.com/codilime/floodgate/version.GoVersion=$GOLANG_VERSION \
-X github.com/codilime/floodgate/version.GateVersion=$(echo ${GATE_API_BRANCH} | sed 's/release-//') \
"

echo "Calculate code coverage"
REQUIREDCODECOVERAGE=60
go tool cover -func cover.out | tee codecoverage.txt
CURRENTCODECOVERAGE=$(grep 'total:' codecoverage.txt | awk '{print substr($3, 1, length($3)-1)}')
curl \
  --header "Authorization: Token ${SERIESCI_TOKEN}" \
  --header "Content-Type: application/json" \
  --data "{\"value\":\"${CURRENTCODECOVERAGE} %\",\"sha\":\"${CIRCLE_SHA1}\"}" \
  https://seriesci.com/api/codilime/floodgate/coverage/one
if [ ${CURRENTCODECOVERAGE%.*} -lt ${REQUIREDCODECOVERAGE} ]
then
    echo "Not enough code coverage!"
    echo "Current code coverage: ${CURRENTCODECOVERAGE}%"
    echo "Required code coverage: ${REQUIREDCODECOVERAGE}%"
    exit 1
else
    echo "Code coverage is at least ${REQUIREDCODECOVERAGE}% : OK"
fi

echo "Check linting"
for GOSRCFILE in $( find . -type f -name '*.go' -not -path './gateapi/*')
do
  golint -set_exit_status $GOSRCFILE
done

echo "Copy binaries for later use"
mkdir -p /floodgate/bin
chmod 777 /floodgate/bin
ls
cp /go/src/github.com/codilime/floodgate/floodgate /floodgate/bin/floodgate

echo "Generate checksum"
cd /go/src/github.com/codilime/floodgate/
cp floodgate floodgate-$GATE_API_BRANCH.$BUILD_OS.$BUILD_ARCH
sha1sum floodgate-$GATE_API_BRANCH.$BUILD_OS.$BUILD_ARCH > floodgate-$GATE_API_BRANCH.$BUILD_OS.$BUILD_ARCH.sha1sum

