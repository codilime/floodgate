#!/bin/bash -e

echo "Check linting"
for GOSRCFILE in $( find . -type f -name '*.go' -not -path './gateapi/*')
do
  golint -set_exit_status $GOSRCFILE
done

