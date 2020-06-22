#!/bin/bash -e

GATE_API_BRANCH=$1

.cilibs/prepare_extra_directories.sh

.cilibs/setup_swagger_codegen.sh

.cilibs/get_gate_code.sh $GATE_API_BRANCH

.cilibs/generate_swagger_json.sh

.cilibs/generate_gateapi_go_code.sh

