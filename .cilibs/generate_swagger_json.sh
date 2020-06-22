#!/bin/bash -e

echo "Generate swagger.json"
cd /floodgate/gate
./gradlew clean
./gradlew gate-web:test --tests *GenerateSwagger* --max-workers 2
cat gate-web/swagger.json | json_pp > ./gate-swagger.json
