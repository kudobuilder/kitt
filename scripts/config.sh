#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

export GOLANGCILINT_VERSION='1.29.0'

export GOBIN=$PWD/bin
export PATH=$GOBIN:$PATH
