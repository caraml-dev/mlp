#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

CHART_PATH="$1"

helm install mlp ${CHART_PATH} --namespace=mlp \
    --values=${CHART_PATH}/values/values.yaml \
    --values=${CHART_PATH}/values/e2e-values.yaml \
    --set mlp.image.tag=${GITHUB_REF#refs/heads/} \
    --dry-run

helm install mlp ${CHART_PATH} --namespace=mlp \
    --values=${CHART_PATH}/values/values.yaml \
    --values=${CHART_PATH}/values/e2e-values.yaml \
    --set mlp.image.tag=${GITHUB_REF#refs/heads/} \
    --wait

kubectl get all --namespace=mlp
