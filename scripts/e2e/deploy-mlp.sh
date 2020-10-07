#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

CHART_PATH="$1"

helm install mlp ${CHART_PATH} --namespace=mlp \
    --values=${CHART_PATH}/values-e2e.yaml \
    --set mlp.image.tag=${GITHUB_REF#refs/*/} \
    --dry-run

helm install mlp ${CHART_PATH} --namespace=mlp \
    --values=${CHART_PATH}/values-e2e.yaml \
    --set mlp.image.tag=${GITHUB_REF#refs/*/} \
    --wait --timeout=5m

kubectl get all --namespace=mlp
