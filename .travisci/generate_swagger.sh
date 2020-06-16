#!/bin/bash -e

GATE_API_BRANCH=$1

echo "Prepare extra directories"
mkdir /floodgate
chmod 777 /floodgate
mkdir /floodgate/bin

echo "Setup swagger-codegen"
SWAGGER_VERSION=$(cat gateapi/.swagger-codegen/VERSION)
wget https://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/${SWAGGER_VERSION}/swagger-codegen-cli-${SWAGGER_VERSION}.jar -O swagger-codegen-cli.jar
wget https://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/${SWAGGER_VERSION}/swagger-codegen-cli-${SWAGGER_VERSION}.jar.sha1 -O swagger-codegen-cli.jar.sha1
echo ' swagger-codegen-cli.jar' >> swagger-codegen-cli.jar.sha1
sha1sum -c swagger-codegen-cli.jar.sha1
mv swagger-codegen-cli.jar /floodgate/bin/

echo "Get gate code"
git clone https://github.com/spinnaker/gate.git -b ${GATE_API_BRANCH} >> /floodgate/gate

echo "Generate swagger.json"
cd /floodgate/gate
./gradlew clean
./gradlew gate-web:test --tests *GenerateSwagger* --max-workers 2
cat gate-web/swagger.json | json_pp > ./gate-swagger.json

echo "Generate gateapi go code"
java -jar /floodgate/bin/swagger-codegen-cli.jar generate -l go -i /floodgate/gate/gate-swagger.json -o /floodgate/gateapi


