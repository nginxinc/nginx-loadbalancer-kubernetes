#!/usr/bin/env bash

set -eo pipefail

if [[ "${CI}" != "true" ]]; then
    echo "This script is meant to be run in the CI."
    exit 1
fi

pttn="^release-[0-9]+\.[0-9]+\.[0-9]+"
if ! [[ "${CI_COMMIT_TAG}" =~ $pttn ]]; then
    echo "CI_COMMIT_TAG needs to be set to valid semver format."
    exit 1
fi

ROOT_DIR=$(git rev-parse --show-toplevel)
source ${ROOT_DIR}/.devops.sh

DOCKERHUB_USERNAME=$(devops.secret.get "kic-dockerhub-creds" | jq -r ".username")
if [[ -z "${DOCKERHUB_USERNAME}" ]]; then
    echo "DOCKERHUB_USERNAME needs to be set."
    exit 1
fi

DOCKERHUB_PASSWORD=$(devops.secret.get "kic-dockerhub-creds" | jq -r ".password")
if [[ -z "${DOCKERHUB_PASSWORD}" ]]; then
    echo "DOCKERHUB_PASSWORD needs to be set."
    exit 1
fi

SRC_REGISTRY="${DEVOPS_DOCKER_URL}"
SRC_PATH="nginx-azure-lb/nginxaas-operator/nginxaas-operator"
SRC_TAG=$(echo "${CI_COMMIT_TAG}" | cut -f 2 -d "-")
SRC_IMG="${SRC_REGISTRY}/${SRC_PATH}:${SRC_TAG}"

DST_REGISTRY="docker.io"
DST_PATH="nginx/nginxaas-operator"
DST_TAG="${CI_COMMIT_TAG}"
DST_IMG="${DST_REGISTRY}/${DST_PATH}:${DST_TAG}"

devops.docker.login
docker pull "${SRC_IMG}"
docker tag "${SRC_IMG}" "${DST_IMG}"

# Login to Dockerhub and push release image to it.
docker login --username "${DOCKERHUB_USERNAME}" --password "${DOCKERHUB_PASSWORD}" "${DST_REGISTRY}"
docker push "${DST_IMG}"
