# generate API libraries using swagger

docker pull swaggerapi/swagger-codegen-cli:latest

docker run --rm -v ${PWD}:/local -u $(id -u):$(id -g) swaggerapi/swagger-codegen-cli generate \
    -i /local/gate-swagger.json \
    -l go \
    -o /local/gateapi