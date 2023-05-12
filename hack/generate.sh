#!/bin/bash

go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
oapi-codegen --config pkg/models/v1/config.yaml pkg/models/v1/spec.yaml > pkg/models/v1/v1.gen.go