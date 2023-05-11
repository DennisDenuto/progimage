#!/bin/bash

go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
oapi-codegen --config models/v1/config.yaml models/v1/spec.yaml > models/v1/v1.gen.go