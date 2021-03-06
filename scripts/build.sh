#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# Download dependencies
echo "Downloading dependencies..."
go mod download

# Set GOPATH
GOPATH="$(go env GOPATH)"

GECKO_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd ) # Directory above this script
BUILD_DIR=$GECKO_PATH/build # Where binaries go
PLUGIN_DIR="$BUILD_DIR/plugins" # Where plugin binaries (namely coreth) go


"$GECKO_PATH/scripts/build_avalanche.sh"

"$GECKO_PATH/scripts/build_coreth.sh"

if [[ -f "$BUILD_DIR/avalanche" && -f "$PLUGIN_DIR/evm" ]]; then
        echo "Build Successful"
else
        echo "Build failure" 
fi
