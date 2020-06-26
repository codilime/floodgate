#!/bin/bash -e

while getopts "o:a:g:e:" opt; do
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
    e) #Floodgate extra params
      FLOODGATE_EXTRA_PARAMS=${OPTARG}
      ;;
  esac
done


.cilibs/prepare_directories.sh 

.cilibs/install_toolset.sh

.cilibs/update_hosts.sh

.cilibs/wait_for_dpkg.sh

.cilibs/install_spinnaker_and_configure_floodgate.sh $GATE_API_BRANCH

.cilibs/test_floodgate_against_running_spinnaker_instance.sh $FLOODGATE_EXTRA_PARAMS
