#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

export RELEASE_NAME="$1"
export MLP_HOST="$2"
export API_PATH="$3"

echo "Running client's examples"
export MLP_API_BASEPATH="http://${MLP_HOST}/v1"

cd ${API_PATH}
for example in client/examples/*; do
    [[ -e $example ]] || continue
    echo $example
    go run $example/main.go
done
