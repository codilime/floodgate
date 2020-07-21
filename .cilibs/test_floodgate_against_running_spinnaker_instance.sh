#!/bin/bash -e

FLOODGATE_EXTRA_PARAMS=$1

echo Test Floodgate against running Spinnaker instance

echo "Print version using version flag"
/floodgate/bin/floodgate --version
echo "Print version using version command"
/floodgate/bin/floodgate version
echo "Comare changes - clean Spinnaker"
/floodgate/bin/floodgate $FLOODGATE_EXTRA_PARAMS --config ~/floodgate.yaml compare && exit 1 || echo "Found changes"
echo "Apply local resources"
/floodgate/bin/floodgate $FLOODGATE_EXTRA_PARAMS --config ~/floodgate.yaml sync
echo "Compare changes - synced resources"
/floodgate/bin/floodgate $FLOODGATE_EXTRA_PARAMS --config ~/floodgate.yaml compare

