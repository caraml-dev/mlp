#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

RELEASE_NAME="$1"

echo "Deleting release: $RELEASE_NAME"
helm delete --purge $RELEASE_NAME
