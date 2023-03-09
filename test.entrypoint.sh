#!/bin/sh

cd /go/src/app
echo "Running Unit Tests..."
go test -v ./...