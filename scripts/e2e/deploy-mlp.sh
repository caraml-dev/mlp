#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

CHART_PATH="$1"

helm repo add caraml https://caraml-dev.github.io/helm-charts/

helm install mlp caraml/mlp --namespace=mlp \
    --values=${CHART_PATH}/values-e2e.yaml \
    --set mlp.image.tag=${GITHUB_REF#refs/*/} \
    --dry-run

helm install mlp caraml/mlp --namespace=mlp \
    --values=${CHART_PATH}/values-e2e.yaml \
    --wait --timeout=5m \
    --set mlp.image.tag=${GITHUB_REF#refs/*/}

kubectl get all --namespace=mlp
