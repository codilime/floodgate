# generate API libraries using swagger in docker

SWAGGER_JSON=${1:-"gate-swagger.json"}
SWAGGER_VERSION=$(cat gateapi/.swagger-codegen/VERSION)

echo '* Using swagger-codegen-cli '${SWAGGER_VERSION}
docker pull swaggerapi/swagger-codegen-cli:${SWAGGER_VERSION}

echo '* Cleaning up gateapi directory'
rm -rf gateapi/api gateapi/docs gateapi/*.go gateapi/*.md

echo '* Generating new gateapi code using '${SWAGGER_JSON}' file...'
docker run --rm -v ${PWD}:/local -u $(id -u):$(id -g) \
    swaggerapi/swagger-codegen-cli:${SWAGGER_VERSION} \
    generate \
    -i /local/${SWAGGER_JSON} \
    -l go \
    -o /local/gateapi

