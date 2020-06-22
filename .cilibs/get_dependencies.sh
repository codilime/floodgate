#!/bin/bash -e

echo "Get dependencies"
go mod download
go get -u golang.org/x/lint/golint
