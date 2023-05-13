#!/bin/bash
set -e

go test -race -v ./...

echo 'Running e2e tests'
go test -v e2e/*.go