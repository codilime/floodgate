#!/bin/bash -e

SEND_COVERITY=$1

echo "Calculate code coverage"
REQUIREDCODECOVERAGE=60
go tool cover -func cover.out | tee codecoverage.txt
CURRENTCODECOVERAGE=$(grep 'total:' codecoverage.txt | awk '{print substr($3, 1, length($3)-1)}')

echo "Send coverity report to SeriesCI"
if [[ $SEND_COVERITY == "send" ]]
then 
  curl \
  --header "Authorization: Token ${SERIESCI_TOKEN}" \
  --header "Content-Type: application/json" \
  --data "{\"value\":\"${CURRENTCODECOVERAGE} %\",\"sha\":\"${CIRCLE_SHA1}\"}" \
  https://seriesci.com/api/codilime/floodgate/coverage/one
else
  echo "Skipping"
fi
if [ ${CURRENTCODECOVERAGE%.*} -lt ${REQUIREDCODECOVERAGE} ]
then
    echo "Not enough code coverage!"
    echo "Current code coverage: ${CURRENTCODECOVERAGE}%"
    echo "Required code coverage: ${REQUIREDCODECOVERAGE}%"
    exit 1
else
    echo "Code coverage is at least ${REQUIREDCODECOVERAGE}% : OK"
fi
