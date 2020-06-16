#!/bin/bash -e

GATE_API_BRANCH=$1
FLOODGATE_EXTRA_PARAMS=""

echo "Update submodules - sponnet"
git submodule init && git submodule update

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

.travisci/install_toolset.sh

echo "Update /etc/hosts"
sudo bash -c 'echo "127.1.2.3 spinnaker" >> /etc/hosts'

.travisci/wait_for_dpkg.sh

echo Install Spinnaker and configure Floodgate
export NEED_SPINNAKER_VERSION=$( echo $GATE_API_BRANCH | egrep -o "[0-9]\.[0-9]+" )
.travisci/install-and-run-spinnaker.sh
until [ $( curl -w '%{http_code}' -o /dev/null http://spinnaker/api/v1 ) -eq 302 ]
do
    echo "Waiting for Spinnaker"
    sleep 10
done

echo Test Floodgate against running Spinnaker instance
/floodgate/bin/floodgate --version
/floodgate/bin/floodgate version
/floodgate/bin/floodgate $FLOODGATE_EXTRA_PARAMS --config ~/floodgate.yaml compare && exit 1 || echo "Found changes"
/floodgate/bin/floodgate $FLOODGATE_EXTRA_PARAMS --config ~/floodgate.yaml sync
/floodgate/bin/floodgate $FLOODGATE_EXTRA_PARAMS --config ~/floodgate.yaml compare

