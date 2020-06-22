#!/bin/bash -e

GATE_API_BRANCH=$1

echo "Get gate code"
git clone https://github.com/spinnaker/gate.git -b ${GATE_API_BRANCH} /floodgate/gate

