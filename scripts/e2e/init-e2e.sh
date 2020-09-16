#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

export RELEASE_NAME="$1"
export MLP_HOST="$2"

echo "Creating merlin project: ${RELEASE_NAME}"
curl "${MLP_HOST}/v1/projects" -d "{
  \"name\"   : \"${RELEASE_NAME}\",
  \"team\"   : \"dsp\",
  \"stream\" : \"dsp\"
}"
sleep 5
