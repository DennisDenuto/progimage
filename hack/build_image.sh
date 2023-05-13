#!/bin/bash
set -xeu

REPO=${1:?"Usage - REPO"}

# https://ko.build/install/
KO_DOCKER_REPO=${REPO} ko build ./