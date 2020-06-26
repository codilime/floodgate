#!/bin/bash -e

echo "Examine source code with go vet"
go vet -v ./...
