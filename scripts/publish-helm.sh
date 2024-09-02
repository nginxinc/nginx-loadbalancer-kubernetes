#!/usr/bin/env bash

set -eo pipefail

ROOT_DIR=$(git rev-parse --show-toplevel)

publish_helm() {
    pkg="nginx-loadbalancer-kubernetes-${VERSION}.tgz"
    helm package --version "${VERSION}" --app-version "${VERSION}" charts/nlk
    helm push "${pkg}" "${repo}"
}

init_ci_vars() {
    if [ -z "$CI_PROJECT_NAME" ]; then
        CI_PROJECT_NAME=$(basename "$ROOT_DIR")
    fi
    if [ -z "$CI_COMMIT_REF_SLUG" ]; then
        CI_COMMIT_REF_SLUG=$(
            git rev-parse --abbrev-ref HEAD | tr "[:upper:]" "[:lower:]" \
                | LANG=en_US.utf8 sed -E -e 's/[^a-zA-Z0-9]/-/g' -e 's/^-+|-+$$//g' \
                | cut -c 1-63
        )
    fi
}

# MAIN
init_ci_vars

# shellcheck source=/dev/null
source "${ROOT_DIR}/.devops.sh"
if [ "$CI" != "true" ]; then
    devops.backend.docker.set "azure.container-registry-dev"
fi
repo="oci://${DEVOPS_DOCKER_URL}/nginx-azure-lb/${CI_PROJECT_NAME}/charts/${CI_COMMIT_REF_SLUG}"
# shellcheck source=/dev/null
# shellcheck disable=SC2153
version=$(source "${ROOT_DIR}/version";echo "$VERSION")

publish_helm
