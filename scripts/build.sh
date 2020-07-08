#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# shellcheck source=scripts/config.sh
source "$(dirname "$0")/config.sh"

COMMAND=$1

if git describe --exact-match >/dev/null; then
	VERSION=$(git describe --dirty)
else
	VERSION=$(git describe --tags --abbrev=0 --dirty)+dev.$(git rev-parse HEAD)
fi

echo "Building $GOBIN/$COMMAND"

go install -ldflags="-X main.version=$VERSION" "./cmd/$COMMAND"
