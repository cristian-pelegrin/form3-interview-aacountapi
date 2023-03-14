#!/bin/sh

cd /go/src/app || exit

echo "----------- Running Unit Tests  ---------------"
go test -v ./form3
echo "------------ Finish Unit Tests  ---------------"

echo "-------- Running Integration Tests ------------"
go test -v ./form3/test/integration
echo "--------- Finish Integration Tests ------------"
