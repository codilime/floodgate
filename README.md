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

TBD (Following #21)

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

## License

Floodgate is licensed under Apache 2.0 License, following other Spinnaker's components.

