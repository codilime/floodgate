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

## How to use?

WIP (Following #19)

## License

Floodgate is licensed under Apache 2.0 License, following other Spinnaker's components.

