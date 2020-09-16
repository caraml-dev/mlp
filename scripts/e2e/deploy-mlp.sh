#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

RELEASE_NAME="$1"
IMAGE_TAG="$2"
NAMESPACE="$3"
CHART_PATH="$4"
VALUES_PATH="$5"
MLP_HOST="$6"

waiting_pod_running(){
    namespace=$1
    TIMEOUT=180
    PODNUM=$(kubectl get deployments,statefulset -n ${namespace} -l release=${RELEASE_NAME}| grep -v NAME | grep -v -e '^$'| wc -l)
    until kubectl get pods -n ${namespace}  -l release=${RELEASE_NAME}| grep -E "Running" | [[ $(wc -l) -ge $PODNUM ]]; do
        echo Pod Running: $(kubectl get pods -n ${namespace}  -l release=${RELEASE_NAME}| grep -E "Running" | wc -l)/$PODNUM

        sleep 10
        TIMEOUT=$(( TIMEOUT - 10 ))
        if [[ $TIMEOUT -eq 0 ]];then
            echo "Timeout to waiting for pod start."
            kubectl describe pods -n ${namespace}  -l release=${RELEASE_NAME}
            exit 1
        fi
    done
    sleep 15
}

echo "Deploying mlp api with release name: $RELEASE_NAME"
helm install ${CHART_PATH} --name ${RELEASE_NAME} --namespace ${NAMESPACE} \
    -f ${VALUES_PATH} \
    --set mlp.encryption.key=${ENCRYPTION_KEY} \
    --set mlp.image.tag=${IMAGE_TAG} \
    --set mlp.ingress.host=${MLP_HOST} \
    --set mlp.oauthClientID=${OAUTH_CLIENT_ID} \
    --set mlp.sentryDSN=${SENTRY_DSN} \
    --set gitlab.clientId=${GITLAB_CLIENT_ID} \
    --set gitlab.clientSecret=${GITLAB_CLIENT_SECRET} \
    --wait
echo "Waiting for mlp api deployment"

waiting_pod_running "mlp"
