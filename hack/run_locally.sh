#!/bin/bash

docker run -p 8080:8080 -e port=8080 -e basePath='/tmp/' $(ko build --local ./)