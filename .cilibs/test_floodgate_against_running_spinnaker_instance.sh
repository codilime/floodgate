#!/bin/bash -e

FLOODGATE_EXTRA_PARAMS=$1

echo Test Floodgate against running Spinnaker instance
/floodgate/bin/floodgate --version
/floodgate/bin/floodgate version
/floodgate/bin/floodgate $FLOODGATE_EXTRA_PARAMS --config ~/floodgate.yaml compare && exit 1 || echo "Found changes"
/floodgate/bin/floodgate $FLOODGATE_EXTRA_PARAMS --config ~/floodgate.yaml sync
/floodgate/bin/floodgate $FLOODGATE_EXTRA_PARAMS --config ~/floodgate.yaml compare

