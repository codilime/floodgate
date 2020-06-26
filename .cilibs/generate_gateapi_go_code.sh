#!/bin/bash -e

echo "Generate gateapi go code"
java -jar /floodgate/bin/swagger-codegen-cli.jar generate -l go -i /floodgate/gate/gate-swagger.json -o /floodgate/gateapi


