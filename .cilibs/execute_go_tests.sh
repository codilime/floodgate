#!/bin/bash -e

echo "Execute go tests"
go test -v ./... -coverprofile cover.out
