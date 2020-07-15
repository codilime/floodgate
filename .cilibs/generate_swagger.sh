#!/bin/bash -e

GATE_API_BRANCH=$1

.cilibs/prepare_extra_directories.sh

.cilibs/setup_swagger.sh $GATE_API_BRANCH

