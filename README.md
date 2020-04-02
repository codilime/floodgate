# Floodgate

[![codilime/floodgate](https://circleci.com/gh/codilime/floodgate/tree/master.svg?style=svg)](https://app.circleci.com/pipelines/github/codilime/floodgate)

## Motivation

This project integrates multiple parts of "as-code" experience in Spinnaker, eg. API, Sponnet, Pipeline Templates for a complete GitOps workflow. 

## Features

- [x] Allows creating configuration of applications, pipeline templates and pipelines as file types:
  - [x] JSONNET
  - [x] JSON
  - [x] YAML
- [x] Updates only the parts of Spinnaker configuration that have actually changed
- [x] Reports diffs in managed objects 
- [ ] Works with all currently supported versions of Spinnaker
- [ ] Is well suited to run in a CI system (single binary!)

## Upcoming features

- [ ] Run as a microservice within Spinnaker installation for seamless integration

## Build process

This repository contains a "known-working" version of `gate-swagger.json` file, which is a definition of Gate's API, used to generate client code in go using `swagger-codegen`.

CircleCI process uses the `swagger-codegen-cli` JAR file to generate the API client at build time.

As a developer, you can use `generate_swaggerapi.sh` script, which will use Docker to generate go client code for Gate API in the `gateapi` directory. For example:

```
./generate_swaggerapi.sh gate-swagger.json
* Using swagger-codegen-cli 2.4.12
2.4.12: Pulling from swaggerapi/swagger-codegen-cli
e7c96db7181b: Already exists
f910a506b6cb: Already exists
b6abafe80f63: Already exists
36ebbdce0651: Pull complete
Digest: sha256:fc24e10784390e27fae893a57da7353d695e641920f3eb58f706ee54af92ebed
Status: Downloaded newer image for swaggerapi/swagger-codegen-cli:2.4.12
docker.io/swaggerapi/swagger-codegen-cli:2.4.12
* Cleaning up gateapi directory
* Generating new gateapi code using gate-swagger.json file...
[main] INFO io.swagger.parser.Swagger20Parser - reading from /local/gate-swagger.json
[...]
```

**Note:** This will remove current contents of `gateapi` directory!

You can also obtain the same using raw java file, like in CI:

```
$ SWAGGER_VERSION=$(cat gateapi/.swagger-codegen/VERSION)
$ wget https://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/${SWAGGER_VERSION}/swagger-codegen-cli-${SWAGGER_VERSION}.jar -O /tmp/swagger-codegen-cli.jar
$ rm -r gateapi
$ java -jar /tmp/swagger-codegen-cli.jar generate -l go -i gate-swagger.json -o gateapi
```

**Note:** This will remove current contents of `gateapi` directory!

## WIP: How to run

This project uses the Go module system.
First, make sure you have some version of Spinnaker running.

Get the repository using `go get`:
```bash
go get github.com/codilime/floodgate
```
In its current development state Floodgate checks if a `config.yaml` is present in the root directory.
Navigate to the directory:
```bash
cd $GOROOT/github.com/codilime/floodgate
```
Create a `config.yaml` file. An example config can be found in `example.yaml`.
Additionally, you can find example resources in the `resources` directory.

Execute `go build` and run the binary:
```bash
go build && ./floodgate
```
You can also use `go install` so that the `floodgate` command is available from any directory.

To run tests use:
```bash
go test ./...
```

## How to use?

WIP (Following #19)

JSON comparison is achieved using an external library. For full output specification please check https://github.com/josephburnett/jd#diff-language

Simple example:

a.json:
`{"hungry":"true", "pizza":{"eat":"true","like":"true"},"pasta":{"eat":"true","like":"false"}}`
b.json:
`{"hungry":"false", "pizza":{"eat":"true","like":"true"},"pasta":{"eat":"false","like":"true"}}`

Difference between a.json and b.json:
```
@ ["hungry"]
- "true"
+ "false"
@ ["pasta","eat"]
- "true"
+ "false"
@ ["pasta","like"]
- "false"
+ "true"
```

## License

Floodgate is licensed under Apache 2.0 License, following other Spinnaker's components.

