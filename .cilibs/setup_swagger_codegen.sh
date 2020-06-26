#!/bin/bash -e

echo "Setup swagger-codegen"
SWAGGER_VERSION=$(cat gateapi/.swagger-codegen/VERSION)
wget https://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/${SWAGGER_VERSION}/swagger-codegen-cli-${SWAGGER_VERSION}.jar -O swagger-codegen-cli.jar
wget https://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/${SWAGGER_VERSION}/swagger-codegen-cli-${SWAGGER_VERSION}.jar.sha1 -O swagger-codegen-cli.jar.sha1
echo ' swagger-codegen-cli.jar' >> swagger-codegen-cli.jar.sha1
sha1sum -c swagger-codegen-cli.jar.sha1
mv swagger-codegen-cli.jar /floodgate/bin/
