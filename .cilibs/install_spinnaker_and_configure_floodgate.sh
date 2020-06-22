#!/bin/bash -e

GATE_API_BRANCH=$1

echo Install Spinnaker and configure Floodgate
export NEED_SPINNAKER_VERSION=$( echo $GATE_API_BRANCH | egrep -o "[0-9]\.[0-9]+" )
.cilibs/install-and-run-spinnaker.sh
until [ $( curl -w '%{http_code}' -o /dev/null http://spinnaker/api/v1 ) -eq 302 ]
do
    echo "Waiting for Spinnaker"
    sleep 10
done

